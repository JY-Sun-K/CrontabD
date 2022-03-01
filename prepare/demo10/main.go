package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.101.138:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Println(err)
	}
	defer cli.Close()

	// lease实现锁自动过期:
	// op操作
	// txn事务: if else then

	// 1, 上锁 (创建租约, 自动续租, 拿着租约去抢占一个key)

	lease := clientv3.NewLease(cli)

	leaseGrantResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	leaseid := leaseGrantResp.ID
	//ctx,_:=context.WithTimeout(context.TODO(),5*time.Second)

	//keepRespChan,err:=lease.KeepAlive(ctx,leaseid)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseid)

	keepRespChan, err := lease.KeepAlive(ctx, leaseid)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经停止")
					goto END
				} else {
					fmt.Println("收到自动续租", keepResp.ID)
				}
			}
		}
	END:
	}()

	kv := clientv3.NewKV(cli)

	txn := kv.Txn(context.TODO())

	// 定义事务

	// 如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseid))).
		Else(clientv3.OpGet("/cron/lock/job9")) // 否则抢锁失败

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		fmt.Println(err)
		return // 没有问题
	}

	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 2, 处理业务

	fmt.Println("处理任务")
	time.Sleep(5 * time.Second)

	// 3, 释放锁(取消自动续租, 释放租约)
	// defer 会把租约释放掉, 关联的KV就被删除了

}
