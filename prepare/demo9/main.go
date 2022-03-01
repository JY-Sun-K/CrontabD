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

	putOp := clientv3.OpPut("/cron/jobs/job8", "123")

	opResp, err := kv.Do(context.TODO(), putOp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入Revision", opResp.Put().Header.Revision)

	getOp := clientv3.OpGet("/cron/jobs/job8")

	opResp, err = kv.Do(context.TODO(), getOp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("数据Revision", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value", string(opResp.Get().Kvs[0].Value))

}
