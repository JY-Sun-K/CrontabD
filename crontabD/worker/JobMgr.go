package worker

import (
	"context"
	"crongo/crontabD/common"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var (
	G_jobMgr *JobMgr
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

//监听任务变化
func (j *JobMgr) watchJobs() error {
	getResp, err := j.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kvpair := range getResp.Kvs {
		job, _ := common.UnpackJob(kvpair.Value)
		jobEvent := common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
		fmt.Println(jobEvent)
		//TODO:同步
		G_scheduler.PushJobEvent(jobEvent)
	}

	go func() {
		watchStartRevision := getResp.Header.Revision + 1
		watchChan := j.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {

				switch watchEvent.Type {
				case mvccpb.PUT:
					job, _ := common.UnpackJob(watchEvent.Kv.Value)

					//构建一个更新的event
					jobEvent := common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
					fmt.Println(jobEvent)
					//todo:推给scheduler
					G_scheduler.PushJobEvent(jobEvent)

				case mvccpb.DELETE:
					//Delete /cron/jobs/job10
					jobName := common.ExtractJobName(string(watchEvent.Kv.Key))
					job := &common.Job{
						Name: jobName,
					}
					jobEvent := common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
					fmt.Println(jobEvent)
					G_scheduler.PushJobEvent(jobEvent)
					//todo:推一个删除事件给scheduler
				}
			}
		}
	}()

	return nil

}

func (j *JobMgr) watchKiller() {
	go func() {

		watchChan := j.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {

				switch watchEvent.Type {
				case mvccpb.PUT:
					jobName := common.ExtractKillerName(string(watchEvent.Kv.Key))
					job := &common.Job{
						Name: jobName,
					}
					jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
					G_scheduler.PushJobEvent(jobEvent)

				case mvccpb.DELETE:

				}
			}
		}
	}()
}

func InitJobMgr() error {
	config := clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	client, err := clientv3.New(config)
	if err != nil {
		return err
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	watcher := clientv3.NewWatcher(client)

	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	err = G_jobMgr.watchJobs()
	if err != nil {
		return err
	}

	G_jobMgr.watchKiller()

	log.Println("ETCD启动成功。。。")
	return nil

}

func (j *JobMgr) CreateJobLock(jobName string) *JobLock {
	jobLock := InitJobLOck(jobName, j.kv, j.lease)
	return jobLock
}
