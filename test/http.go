package test

import (
	"github.com/lwl1989/TTTask/msg"
	"net/http"
)

func main()  {
	config := msg.GetConfig("/www/go_path/src/github.com/lwl1989/TTTask/conf/config.json.back")
	handler := msg.GetHandle()
	handler.SetConfig(config)

	server := &http.Server{
		Addr: ":8080",
		Handler:handler,
	}

	server.ListenAndServe()
}