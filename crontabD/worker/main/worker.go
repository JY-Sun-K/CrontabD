package main

import (
	"crongo/crontabD/worker"
	"flag"
	"log"
	"runtime"
	"time"
)

var (
	confFile string
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {

	//worker -config ./worker.json
	//worker -h
	flag.StringVar(&confFile, "config", "./worker.json", "-config ./worker.json")
	flag.Parse()
}

func main() {
	var err error
	initArgs()
	//初始化线程
	initEnv()
	//errChan := make(chan error,1)

	err = worker.IntiConfig(confFile)
	if err != nil {
		log.Println("初始化错误:", err)
		return
	}
	err = worker.InitRegister()
	if err != nil {
		log.Println("服务注册:", err)
		return
	}

	err = worker.InitLogSink()
	if err != nil {
		log.Println(err)
		return
	}

	err = worker.InitExecutor()
	if err != nil {
		log.Println(err)
		return
	}

	worker.InitScheduler()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	err = worker.InitJobMgr()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("服务已运行")

	for {
		time.Sleep(1 * time.Second)
	}

}
