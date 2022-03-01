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
	//用于读写etcd的键值对
	kv := clientv3.NewKV(cli)
	if putResp, err := kv.Put(context.TODO(), "/cron/jobs/job2", "bye", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}
