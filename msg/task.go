package msg

import (
	"sync"
	"github.com/lwl1989/timing"
)

var taskOnce sync.Once
var job *timing.OnceCron


func GetTask() *timing.OnceCron  {
	taskOnce.Do(func() {
		job = timing.NewCron()
	})
	return job
}

func AddToTask(fcmMsg *FcmMsg)  {

	GetTask().AddTask(&timing.Task{
		Job:fcmMsg,
		RunTime:fcmMsg.sendTime,
	})
	go GetTask().Start()
}