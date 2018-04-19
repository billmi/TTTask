package test

import (
	"github.com/robfig/cron"
	"fmt"
	"time"
	"strconv"
)

func main()  {
	//分 时 日 月 年
	now := time.Now().Add(5*time.Second)
	str := strconv.Itoa(now.Second())+" "+strconv.Itoa(now.Minute())+" "+strconv.Itoa(now.Hour())+" "+strconv.Itoa(now.Day())+" "+strconv.Itoa(int(now.Month()))
	fmt.Println(str)
	c := cron.New()
	err := c.AddFunc(str, func() { fmt.Println("time") })
	c.AddFunc("1 * * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	fmt.Println(err)
	go c.Start()
	path := ""
	fmt.Scanln(&path)

}