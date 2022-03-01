package master

import (
	"context"
	"crongo/crontabD/common"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var (
	G_jobMgr *JobMgr
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
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

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	log.Println("ETCD启动成功。。。")
	return nil

}

func (j *JobMgr) SaveJob(job *common.Job) (*common.Job, error) {
	jobKey := common.JOB_SAVE_DIR + job.Name

	jobValue, err := json.Marshal(job)
	if err != nil {
		return nil, err

	}

	putResp, err := j.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	oldJobObj := common.Job{}
	if putResp.PrevKv != nil {
		err := json.Unmarshal(putResp.PrevKv.Value, &oldJobObj)
		if err != nil {
			return nil, nil
		}
		oldJob := &oldJobObj
		return oldJob, nil
	}

	return nil, nil

}

func (j *JobMgr) DeleteJob(name string) (*common.Job, error) {
	jobKey := common.JOB_SAVE_DIR + name

	delResp, err := j.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}
	oldJobObj := common.Job{}
	if len(delResp.PrevKvs) != 0 {
		err := json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj)
		if err != nil {
			return nil, nil
		}
		oldJob := &oldJobObj
		return oldJob, nil
	}

	return nil, nil

}

func (j *JobMgr) ListJobs() ([]*common.Job, error) {
	dirKey := common.JOB_SAVE_DIR

	getResp, err := j.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	jobList := make([]*common.Job, 0)
	for _, kvPair := range getResp.Kvs {
		job := &common.Job{}
		err := json.Unmarshal(kvPair.Value, job)
		if err != nil {
			continue
		}
		jobList = append(jobList, job)
	}
	return jobList, nil

}

func (j *JobMgr) KillJob(name string) error {
	killerKey := common.JOB_KILLER_DIR + name

	leaseGrantResp, err := j.lease.Grant(context.TODO(), 1)
	if err != nil {
		return err
	}

	leaseId := leaseGrantResp.ID

	_, err = j.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}
	return nil

}
