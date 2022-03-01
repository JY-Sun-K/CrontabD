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
	kv := clientv3.NewKV(cli)
	getResp, err := kv.Get(context.TODO(), "/cron/jobs/job1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getResp.Kvs)
}
