package worker

import (
	"context"
	"crongo/crontabD/common"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"time"
)

type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIP string
}

var G_register *Register

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, isIPNet := addr.(*net.IPNet)
		if isIPNet && ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 := ipNet.IP.String()
				return ipv4, nil
			}
		}
	}
	return "", errors.New("没有找到ip网卡")
}

func (r *Register) keepOnline() {

	var cancel context.CancelFunc
	for {
		regKey := common.JOB_WORKERS_DIR + r.localIP
		cancel = nil

		leaseGrantResp, err := r.lease.Grant(context.TODO(), 10)
		if err != nil {
			time.Sleep(1 * time.Second)
			if cancel != nil {
				cancel()
			}

		}

		keepAliveChan, err := r.lease.KeepAlive(context.TODO(), leaseGrantResp.ID)
		if err != nil {
			time.Sleep(1 * time.Second)
			if cancel != nil {
				cancel()
			}
		}

		cancelCtx, cancel := context.WithCancel(context.TODO())

		_, err = r.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID))
		if err != nil {
			time.Sleep(1 * time.Second)
			if cancel != nil {
				cancel()
			}
		}

		for {
			select {
			case keepAliveResp := <-keepAliveChan:
				if keepAliveResp == nil {
					time.Sleep(1 * time.Second)
					if cancel != nil {
						cancel()
					}
				}

			}
		}

	}

}

func InitRegister() error {
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

	localIP, err := getLocalIP()
	if err != nil {
		return err
	}
	G_register = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIP,
	}

	go G_register.keepOnline()
	return nil
}
