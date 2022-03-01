package common

import (
	"context"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type JobExecuteInfo struct {
	Job      *Job
	PlanTime time.Time
	RealTime time.Time

	CancelCtx context.Context
	Cancel    context.CancelFunc
}

type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type JobEvent struct {
	EventType int
	Job       *Job
}

type JobSchedulerPlan struct {
	Job      *Job
	Expr     cron.Schedule
	NextTime time.Time
}

type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
}

type JobLog struct {
	JobName      string `bson:"jobName"`
	Command      string `bson:"command"`
	Err          string `bson:"err"`
	OutPut       string `bson:"output"`
	PlanTime     int64  `bson:"planTime"`
	ScheduleTime int64  `bson:"scheduleTime"`
	StartTime    int64  `bson:"startTime"`
	EndTime      int64  `bson:"endTime"`
}

type LogBatch struct {
	Logs []interface{}
}

type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
}

func BuildResponse(errno int, msg string, data interface{}) ([]byte, error) {
	response := &Response{
		Errno: errno,
		Msg:   msg,
		Data:  data,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resp, nil
}

func UnpackJob(value []byte) (*Job, error) {
	job := &Job{}
	err := json.Unmarshal(value, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// /cron/jobs/job10 ==> job10
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func ExtractKillerName(KillerKey string) string {
	return strings.TrimPrefix(KillerKey, JOB_KILLER_DIR)
}

func BuildJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildJobSchedulerPlan(job *Job) (*JobSchedulerPlan, error) {

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	scheduler, err := parser.Parse(job.CronExpr)
	if err != nil {

		return nil, err
	}

	return &JobSchedulerPlan{
		Job:      job,
		Expr:     scheduler,
		NextTime: scheduler.Next(time.Now()),
	}, nil
}

func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulerPlan) *JobExecuteInfo {
	cancelCtx, cancel := context.WithCancel(context.TODO())
	return &JobExecuteInfo{
		Job:       jobSchedulePlan.Job,
		PlanTime:  jobSchedulePlan.NextTime,
		RealTime:  time.Now(),
		CancelCtx: cancelCtx,
		Cancel:    cancel,
	}
}

func ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKERS_DIR)
}
