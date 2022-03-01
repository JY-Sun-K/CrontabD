package master

import (
	"context"
	"crongo/crontabD/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var G_logMgr *LogMgr

type LogMgr struct {
	Client        *mongo.Client
	LogCollection *mongo.Collection
}

func IntiLogMgr() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(G_config.MongodbUri))
	if err != nil {
		log.Println(err)
		return err
	}
	G_logMgr = &LogMgr{
		Client:        client,
		LogCollection: client.Database("cron").Collection("log"),
	}

	return nil
}

func (l *LogMgr) ListLogs(name string, skip int, limit int) ([]*common.JobLog, error) {
	logArr := make([]*common.JobLog, 0)

	filter := &common.JobLogFilter{JobName: name}

	logSort := &common.SortLogByStartTime{SortOrder: -1}

	cursor, err := l.LogCollection.Find(context.TODO(), filter, options.Find().SetSort(logSort), options.Find().SetSkip(int64(skip)), options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog := &common.JobLog{}
		err := cursor.Decode(jobLog)
		if err != nil {
			continue
		}
		logArr = append(logArr, jobLog)

	}

	return logArr, nil

}
