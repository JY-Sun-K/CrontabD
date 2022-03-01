package master

import (
	"crongo/crontabD/common"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

var G_apiServer *ApiServer

type ApiServer struct {
	httpServer *http.Server
}

//POST job={"name":"job1,"command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	postJob := r.PostForm.Get("job")
	job := common.Job{}
	err = json.Unmarshal([]byte(postJob), &job)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	oldJob, err := G_jobMgr.SaveJob(&job)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	bytes, err := common.BuildResponse(0, "success", oldJob)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	w.Write(bytes)

}

func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	name := r.PostForm.Get("name")

	oldJob, err := G_jobMgr.DeleteJob(name)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	bytes, err := common.BuildResponse(0, "success", oldJob)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	w.Write(bytes)

}

func handleJobList(w http.ResponseWriter, r *http.Request) {
	jobList, err := G_jobMgr.ListJobs()
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	bytes, err := common.BuildResponse(0, "success", jobList)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	w.Write(bytes)
}

func handleJobKill(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	name := r.PostForm.Get("name")

	err = G_jobMgr.KillJob(name)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	bytes, err := common.BuildResponse(0, "success", nil)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	w.Write(bytes)
}

func handleJobLog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	// /job/log?name=job10&skip=0&limit=10
	name := r.Form.Get("name")
	skipParam := r.Form.Get("skip")
	limitParam := r.Form.Get("limit")

	skip, err := strconv.Atoi(skipParam)
	if err != nil {
		skip = 0
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 20
	}

	logArr, err := G_logMgr.ListLogs(name, skip, limit)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}

	bytes, err := common.BuildResponse(0, "success", logArr)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	w.Write(bytes)

}

func handleWorkerList(w http.ResponseWriter, r *http.Request) {
	workerArr, err := G_workerMgr.ListWorkers()
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	bytes, err := common.BuildResponse(0, "success", workerArr)
	if err != nil {
		log.Println(err)
		bytes, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(bytes)
		return
	}
	w.Write(bytes)
}

func InitApiServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/worker/list", handleWorkerList)

	staticDir := http.Dir("./webroot")
	staticHandler := http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort))
	if err != nil {

		return err
	}

	httpServer := &http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
	}

	G_apiServer = &ApiServer{httpServer: httpServer}

	go httpServer.Serve(listener)

	log.Println("apiHttp启动成功。。。")
	return nil

}
