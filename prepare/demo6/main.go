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

	// 用于读写etcd的键值对

	// 删除KV
	delResp, err := kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
		return
	}

	// 被删除之前的value是什么
	if len(delResp.PrevKvs) != 0 {
		for _, kvpair := range delResp.PrevKvs {
			fmt.Println("删除了:", string(kvpair.Key), string(kvpair.Value))
		}
	}
}
