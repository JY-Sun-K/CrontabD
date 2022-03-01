package master

import (
	"context"
	"crongo/crontabD/common"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var G_workerMgr *WorkerMgr

func (w *WorkerMgr) ListWorkers() ([]string, error) {
	workerArr := make([]string, 0)
	getResp, err := w.kv.Get(context.TODO(), common.JOB_WORKERS_DIR, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range getResp.Kvs {
		workerIP := common.ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, workerIP)
	}
	return workerArr, nil
}

func InitWorkerMgr() error {
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

	G_workerMgr = &WorkerMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return nil
}
