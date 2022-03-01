package worker

import (
	"context"
	"crongo/crontabD/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type LogSink struct {
	Client         *mongo.Client
	LogCollection  *mongo.Collection
	LogChan        chan *common.JobLog
	AutoCommitChan chan *common.LogBatch
}

var (
	G_logSink *LogSink
)

func (l *LogSink) saveLogs(batch *common.LogBatch) {
	l.LogCollection.InsertMany(context.TODO(), batch.Logs)
}

func (l *LogSink) writeLoop() {
	var (
		logBatch    *common.LogBatch
		commitTimer *time.Timer
	)

	for {
		select {
		case log := <-l.LogChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				commitTimer = time.AfterFunc(time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						//发出超时通知，不要直接提交batch ，防止并发调度logBatch
						//l.AutoCommitChan<- logBatch
						return func() {
							l.AutoCommitChan <- batch
						}

					}(logBatch),
				)
			}
			logBatch.Logs = append(logBatch.Logs, log)

			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				l.saveLogs(logBatch)

				logBatch = nil

				commitTimer.Stop()
			}
		case timeoutBatch := <-l.AutoCommitChan:
			if timeoutBatch != logBatch {
				continue
			}

			l.saveLogs(timeoutBatch)
			logBatch = nil

		}
	}
}

func InitLogSink() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(G_config.MongodbUri))
	if err != nil {
		log.Println(err)
		return err
	}

	G_logSink = &LogSink{
		Client:         client,
		LogCollection:  client.Database("cron").Collection("log"),
		LogChan:        make(chan *common.JobLog, 1000),
		AutoCommitChan: make(chan *common.LogBatch, 1000),
	}

	go G_logSink.writeLoop()

	return nil
}

func (l *LogSink) Append(jobLog *common.JobLog) {
	select {
	case l.LogChan <- jobLog:
	default:

	}

}
