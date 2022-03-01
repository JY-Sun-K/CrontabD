package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"time"
)

type cronJob struct {
	JobName string
}

func main() {
	c1 := &cronJob{JobName: "c1"}
	c2 := &cronJob{JobName: "c2"}
	//specParser := cron.NewParser(cron.Second )
	//sched, err := specParser.Parse("*/3")
	//if err != nil {
	//	log.Println(err)
	//}
	c := cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron process: ", log.LstdFlags))),
		cron.WithParser(cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
	)
	c.AddFunc("*/2 * * * * ", func() {

		fmt.Println(c1.JobName)
		fmt.Println("lol1")
	})
	c.AddFunc("@every 2s", func() {
		fmt.Println(c2.JobName)
		fmt.Println("lol2")
	})

	fmt.Println(c.Entries())
	c.Start()
	select {
	case <-time.NewTimer(8 * time.Second).C:

	}
}
