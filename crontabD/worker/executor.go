package worker

import (
	"crongo/crontabD/common"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os/exec"
	"time"
)

var (
	G_executor *Executor
)

type Executor struct {
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func (e *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		//cmd:=exec.CommandContext(context.TODO(),"/bin/bash","-c",info.Job.Command)
		//output,err:=cmd.CombinedOutput()

		jobLock := G_jobMgr.CreateJobLock(info.Job.Name)

		//time.Sleep(time.Duration(rand.Intn(1000))*time.Millisecond)

		err := jobLock.TryLock()
		defer jobLock.UnLock()
		if err != nil {
			result := &common.JobExecuteResult{
				ExecuteInfo: info,
				Output:      make([]byte, 0),
				Err:         err,
				StartTime:   time.Now(),
				EndTime:     time.Now(),
			}
			G_scheduler.PushJobResult(result)
		} else {
			startTime := time.Now()
			cmd := exec.CommandContext(info.CancelCtx, "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", info.Job.Command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(err)

			}
			endTime := time.Now()
			cmdRe := ConvertByte2String(output, "GB18030")
			//log.Println(string(cmdRe))

			result := &common.JobExecuteResult{
				ExecuteInfo: info,
				Output:      cmdRe,
				Err:         err,
				StartTime:   startTime,
				EndTime:     endTime,
			}
			G_scheduler.PushJobResult(result)

		}

	}()
}

func InitExecutor() error {
	return nil
}

func ConvertByte2String(byt []byte, charset Charset) []byte {
	var str []byte
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byt)
		str = decodeBytes
	case UTF8:
		fallthrough
	default:
		str = byt
	}
	return str
}
