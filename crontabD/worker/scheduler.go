package worker

import (
	"crongo/crontabD/common"
	"errors"
	"log"
	"time"
)

type Scheduler struct {
	JobEventChan      chan *common.JobEvent
	JobPlanTable      map[string]*common.JobSchedulerPlan
	JobExecutingTable map[string]*common.JobExecuteInfo
	JobResultChan     chan *common.JobExecuteResult
}

var (
	G_scheduler *Scheduler
)

func (s *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		jobSchedulePlan, err := common.BuildJobSchedulerPlan(jobEvent.Job)
		if err != nil {
			log.Println(err)
			return
		}
		s.JobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
		log.Println("ADD jobSchedulePlan", jobEvent.Job.Name, "&", jobSchedulePlan)

	case common.JOB_EVENT_DELETE:
		if _, ok := s.JobPlanTable[jobEvent.Job.Name]; ok {
			delete(s.JobPlanTable, jobEvent.Job.Name)
			log.Println("ADD jobSchedulePlan", jobEvent.Job.Name)
		}

	case common.JOB_EVENT_KILL:
		if jobExecuteInfo, ok := s.JobExecutingTable[jobEvent.Job.Name]; ok {
			jobExecuteInfo.Cancel()
		}

	}
}

func (s *Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan) {

	if _, ok := s.JobExecutingTable[jobPlan.Job.Name]; ok {
		log.Println("正在运行，跳过执行:", jobPlan.Job.Name)
		return
	}
	jobExecuteInfo := common.BuildJobExecuteInfo(jobPlan)

	s.JobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	log.Println("执行任务:", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	G_executor.ExecuteJob(jobExecuteInfo)
}

func (s *Scheduler) TrySchedule() time.Duration {
	var nearTime *time.Time
	if len(s.JobPlanTable) == 0 {
		return 1 * time.Second
	}

	//range all jobs
	now := time.Now()

	for _, jobPlan := range s.JobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//todo: 尝试执行任务
			log.Println("执行的任务是:", jobPlan.Job.Name)
			s.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}
		//统计最近一个任务的到期时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}

	}

	return (*nearTime).Sub(now)

}

func (s *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	delete(s.JobExecutingTable, result.ExecuteInfo.Job.Name)

	if result.Err != errors.New("锁已经被占用") {
		jobLog := &common.JobLog{
			JobName: result.ExecuteInfo.Job.Name,
			Command: result.ExecuteInfo.Job.Command,
			//Err:          result.Err.Error(),
			OutPut:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}

		if result.Err == nil {
			jobLog.Err = ""
		} else {
			jobLog.Err = result.Err.Error()
		}

		G_logSink.Append(jobLog)
	}
	//log.Println("任务执行完成",result.ExecuteInfo.Job.Name,string(result.Output),result.Err)
	log.Println("任务执行完成", result.ExecuteInfo.Job.Name, result.Err)

}

func (s *Scheduler) schedulerLoop() {

	scheduleAfter := s.TrySchedule()

	scheduleTimer := time.NewTimer(scheduleAfter)

	for {
		select {
		case jobEvent := <-s.JobEventChan:
			s.handleJobEvent(jobEvent)
		case <-scheduleTimer.C:
		case jobResult := <-s.JobResultChan:
			s.handleJobResult(jobResult)

		}

		scheduleAfter = s.TrySchedule()

		scheduleTimer.Reset(scheduleAfter)
	}

}

func (s *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	s.JobEventChan <- jobEvent
}

func InitScheduler() {
	G_scheduler = &Scheduler{
		JobEventChan:      make(chan *common.JobEvent, 1000),
		JobPlanTable:      make(map[string]*common.JobSchedulerPlan),
		JobExecutingTable: make(map[string]*common.JobExecuteInfo),
		JobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}

	go G_scheduler.schedulerLoop()

}

func (s *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	s.JobResultChan <- jobResult
}
