package main

import (
	"net/http"
	"github.com/lwl1989/TTTask/msg"
	"fmt"
	"github.com/lwl1989/TTTask/cmd"
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
	fmt.Println("load config success")
	handler := msg.GetHandle()
	handler.SetConfig(config)
	http.ListenAndServe(":"+config.GetServerPort(), handler)

}
