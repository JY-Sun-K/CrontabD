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

	lease := clientv3.NewLease(cli)

	leaseGrantResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	leaseid := leaseGrantResp.ID
	//ctx,_:=context.WithTimeout(context.TODO(),5*time.Second)

	//keepRespChan,err:=lease.KeepAlive(ctx,leaseid)
	keepRespChan, err := lease.KeepAlive(context.TODO(), leaseid)
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

	putResp, err := kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseid))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功", putResp.Header.Revision)

	for {
		getResp, err := kv.Get(context.TODO(), "/cron/lock/job1")
		if err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("kv 过期了")
			break
		}
		fmt.Println("还没过期", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}

}
