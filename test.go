package main

import (
	"log"
	"github.com/appleboy/go-fcm"
	"reflect"
	"fmt"
	msg2 "github.com/lwl1989/TTTask/msg"
)

func New(sample interface{}) interface{} {
	t :=reflect.ValueOf(sample).Type()
	fmt.Println(t.Kind() == reflect.String)
	v :=reflect.New(t).Interface()
	return v
}
func getMessage() *fcm.Message  {
	message := map[string]interface{}{
		"title":"fsa",
	}

	return &fcm.Message{
		Data: message,
	}
}
var config1 *msg2.ApiConfig

func main() {
	config1 = msg2.GetConfig("/www/go_path/src/github.com/lwl1989/TTTask/conf/config.json")
	f := getMessage()
	ff := msg2.GetNewFcm()
	ff.SetConfig(config1)
	ff.SetMessage(f)

	ff.SetTopic("/topics/aaa")
	a,b := ff.Send()
	fmt.Println(b)
	log.Printf("%v",a.MessageID)
	return
	// Create the message to be sent.
	msg := &fcm.Message{
		To: "sample_device_token",
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}

	// Create a FCM client to send the message.
	client, err := fcm.NewClient("AAAA_1dLSps:APA91bF11IBdXuZ4OKir5Z5di7_yvLCqShwCKRolrW6Et-orZcTWL5__maDTlGN7mezq4ROXjjkp2xhb1sCNPAhNa1dCU0bAYR65GwVLlgQ5KpX3tZpyaZHrCUioe-vx6wFvDXfnoh9h")
	if err != nil {
		log.Fatalln(err)
	}

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	log.Println(response)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", response)
}