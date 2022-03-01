package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

func main() {
	//  c := cron.New()
	//   c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	// c.AddFunc("30 3-6,20-23 * * *", func() { fmt.Println(".. in the range 3-6am, 8-11pm") })
	// c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })
	//   c.AddFunc("@hourly", func() { fmt.Println("Every hour, starting an hour from now") })
	//  c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty, starting an hour thirty from now") })
	//  c.AddFunc("@every 1s", func() {fmt.Println("Every 1 second, starting an hour thirty from now")})
	// c.Start()
	//select {}

	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron process: ", log.LstdFlags))))
	_, err := c.AddFunc("*/6 * * * *", func() {
		log.Println("owo")
	})
	if err != nil {
		log.Println(err)
	}
	c.Start()
	select {}
}
