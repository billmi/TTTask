package test

import (
	"time"
	"math/rand"
	"fmt"
)

type Schedule interface {
	// Return the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next() *Entity
}

type fff struct {
}

func (fff *fff) Send() {
	fmt.Println(rand.Int63())
}

type Entity struct {
	message    *fff
	timer      int64
	execTime   *time.Timer
	execTick   *time.Ticker
	exec       bool
	NextEntity *Entity
}

type Task interface {
	AddTimer(unixTime int64, message *fff)
	AddTicker(unixTime int64, message *fff)
	Start()
	Clear()
	GetNextEntity() *Entity
}

type JobEntity interface {
	HadExec() (exec bool)
	Send()
}

type Jobs struct {
	Entities []*Entity
}

func (entity *Entity) Send() {
	entity.message.Send()
	entity.exec = true
}

func (entity *Entity) HadExec() (exec bool) {
	return entity.exec
}

func (jobs *Jobs) GetNextEntity() *Entity {

	if len(jobs.Entities) > 0 {
		if len(jobs.Entities) == 1 {
			return jobs.Entities[0]
		}

		i64 := jobs.Entities[0].timer
		temp := jobs.Entities[0]

		for _, e := range jobs.Entities {
			if e.timer <= i64 {
				i64 = e.timer
				temp = e
			}
		}

		return temp
	}
	return nil
}

func (jobs *Jobs) AddTimer(unixTime int64, message *fff) {
	d := time.Unix(unixTime, 0).Sub(time.Now())
	fmt.Println(time.Unix(unixTime, 0))
	fmt.Println(time.Now())
	fmt.Println(d)
	en :=  &Entity{
		message:  message,
		timer:    unixTime,
		execTime: time.NewTimer(d),
	}

	jobs.Entities = append(jobs.Entities,en)

}

func (jobs *Jobs) AddTicker(unixTime int64, message *fff) {
	d := time.Unix(unixTime, rand.Int63()).Sub(time.Now())

	jobs.Entities = append(jobs.Entities, &Entity{
		message:  message,
		timer:    unixTime,
		execTime: time.NewTimer(d),
	})
	//fmt.Println(len(jobs.Entities))
}

func GetJobs() (jobs *Jobs) {
	jobs = &Jobs{}
	jobs.Entities = make([]*Entity, 0)
	return jobs
}

func (jobs *Jobs) Start() {

	entity := jobs.GetNextEntity()
	fmt.Println(time.Unix(entity.timer,0))
	//entity.Send()
	//fmt.Println(entity.timer)
	for {

		//fmt.Println(entity)
		select {
		case <-entity.execTime.C:
			entity.Send()
			break
		}

	}
}

func (jobs *Jobs) Clear() {
	jobs.Entities = make([]*Entity, 0)
}

func main() {
	jobs := GetJobs()


	jobs.AddTicker(time.Now().Add(time.Second * 30).Unix(), &fff{})
	jobs.AddTimer(time.Now().Add(time.Second * 20).Unix(), &fff{})
	jobs.AddTimer(time.Now().Add(time.Second * 10).Unix(), &fff{})
	jobs.Start()
	path := ""
	fmt.Scanln(&path)


}
