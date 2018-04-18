package main

import (
	"net/http"
	"github.com/lwl1989/TTTask/msg"
	"fmt"
)

var config *msg.ApiConfig

func main()  {

	path := ""
	fmt.Println("Please input your config path(default /usr/etc/fcm.json): ")
	fmt.Scanln(&path)

	if path == "" {
		panic("error! no input field!")
	}
	config = msg.GetConfig(path)

	handler := msg.GetHandle()
	handler.SetConfig(config)

	http.ListenAndServe(":"+config.GetServerPort(), msg.GetHandle())
}