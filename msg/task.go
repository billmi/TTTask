package msg

import (
	"sync"
	"github.com/robfig/cron"
	"time"
	"strconv"
	"github.com/appleboy/go-fcm"
	"log"
)

var taskOnce sync.Once
var job *cron.Cron
var taskEntity []*task
type task struct {
	fcmMsg *FcmMsg
	response *fcm.Response
}

func GetTask() *cron.Cron  {
	taskOnce.Do(func() {
		job = cron.New()
		job.AddFunc("@weekly", func() {
			log.Panicln("每周運行任務，保證協程不被退出")
		})
		taskEntity = make([]*task,0)
	})
	return job
}

func AddToTask(fcmMsg *FcmMsg)  {
	//分 时 日 月 年
	sendTime := time.Unix(fcmMsg.sendTime, 0)
	spec := strconv.Itoa(sendTime.Second())+" "+strconv.Itoa(sendTime.Minute())+" "+strconv.Itoa(sendTime.Hour())+" "+strconv.Itoa(sendTime.Day())+" "+strconv.Itoa(int(sendTime.Month()))
	log.Println("Add sendTask with time "+ spec)
	GetTask().AddFunc(spec, func() {
		response,err := fcmMsg.Send()
		log.Println("Send!")
		taskEntity = append(taskEntity, &task{fcmMsg,response})
		log.Println(taskEntity)
		// and add task list restful
		status := 200
		if err != nil {
			status = 500
		}
		call := GetCallBack(string(fcmMsg.messageId),  status)
		call.SetConfig(fcmMsg.conf)
		call.Do()
	})
	GetTask().Start()
}