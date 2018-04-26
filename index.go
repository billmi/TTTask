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
	fmt.Println("like this:")
	fmt.Println(`{"server_port":"8080","max_ttl":2419200,"api_key":"AAAA_1dLSps:APA91......ZHrCUioe-vx6wFvDXfnoh9h"}`)
	fmt.Scanln(&path)
	、、path = "/www/go_path/src/github.com/lwl1989/TTTask/conf/config.json.back"
	if path == "" {
		panic("error! no input field!")
	}
	config = msg.GetConfig(path)
	//os.Setenv("HTTP_PROXY", "http://127.0.0.1:1080")
	//os.Setenv("HTTPS_PROXY", "https://127.0.0.1:1080")
	fmt.Println("load config success")
	handler := msg.GetHandle()
	handler.SetConfig(config)
	http.ListenAndServe(":"+config.GetServerPort(), handler)
}