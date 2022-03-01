package worker

import (
	"context"
	"crongo/crontabD/common"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type JobLock struct {
	kv      clientv3.KV
	Lease   clientv3.Lease
	JobName string
	LeaseId clientv3.LeaseID
	cancel  context.CancelFunc
	isLock  bool
}

func InitJobLOck(jobName string, kv clientv3.KV, lease clientv3.Lease) *JobLock {
	return &JobLock{
		kv:      kv,
		Lease:   lease,
		JobName: jobName,
	}
}

func (j *JobLock) TryLock() error {

	//1.创建租约
	leaseGrantResp, err := j.Lease.Grant(context.TODO(), 5)
	if err != nil {
		return err
	}

	//2.自动续租
	cancelCtx, cancel := context.WithCancel(context.TODO())

	leaseId := leaseGrantResp.ID

	keepRespChan, err := j.Lease.KeepAlive(cancelCtx, leaseId)
	if err != nil {
		cancel()
		j.Lease.Revoke(context.TODO(), leaseId)
		return err
	}

	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:

				if keepResp == nil {
					//fmt.Println("租约已经停止")
					goto END
				} else {
					//fmt.Println("收到自动续租",keepResp.ID)
				}

			}
		}
	END:
	}()

	//3.创建事务txn
	txn := j.kv.Txn(context.TODO())
	lockKey := common.JOB_LOCK_DIR + j.JobName

	//4.事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	txnResp, err := txn.Commit()
	if err != nil {
		cancel()
		j.Lease.Revoke(context.TODO(), leaseId)
		return err
	}
	//5.成功返回，失败释放租约
	if !txnResp.Succeeded {
		cancel()
		j.Lease.Revoke(context.TODO(), leaseId)
		return errors.New("锁已经被占用")
	}

	j.LeaseId = leaseId
	j.cancel = cancel
	j.isLock = true

	return nil
}

func (j *JobLock) UnLock() {
	if j.isLock {
		j.cancel()
		j.Lease.Revoke(context.TODO(), j.LeaseId)
	}

}
