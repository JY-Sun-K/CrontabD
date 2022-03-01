package main

import (
	"crongo/crontabD/master"
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

	//master -config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "-config ./master.json")
	flag.Parse()
}

func main() {
	var err error
	initArgs()
	//初始化线程
	initEnv()
	//errChan := make(chan error,1)
	err = master.IntiConfig(confFile)
	if err != nil {
		log.Println(err)
		return
	}
	err = master.InitWorkerMgr()
	if err != nil {
		log.Println(err)
		return
	}

	err = master.InitJobMgr()
	if err != nil {
		log.Println(err)
		return
	}
	err = master.IntiLogMgr()
	if err != nil {
		log.Println(err)
		return
	}

	//启动http api 服务
	err = master.InitApiServer()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("服务已运行")

	for {
		time.Sleep(1 * time.Second)
	}

}
