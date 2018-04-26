package main

import (
	"github.com/lwl1989/TTTask/msg"
	"fmt"
	"github.com/lwl1989/TTTask/cmd"
	"github.com/takama/daemon"
	"log"
)

var config *msg.ApiConfig

func main() {
	c := cmd.GetCommand()

	h := c.Get("help")
	path := c.Get("config")

	if h != "" {
		cmd.ShowHelp()
		return
	}

	if path == "" {
		cmd.ShowHelp()
		return
	}


	config = msg.GetConfig(path)

	srv, err := daemon.New("ttpush", "ttpush fcm proxy service")
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	service := &cmd.Service{ Daemon:srv,Config:config }
	status, err := service.Manage()
	if err != nil {
		log.Fatalln(status, "\nError: ", err)
	}
	fmt.Println(status)
}
