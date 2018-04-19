package test

import (
	"fmt"
	"encoding/json"
)

type fcmError struct {
	Code int `json:"code:omitempty"`
	Response interface{} `json:"response:omitempty"`
}

func main()  {
	fcm := fcmError{400,"dfasfddsgfdg"}
	c,_ := json.Marshal(fcm)
	c1,_ := json.Marshal(&fcm)
	fmt.Println(string(c[:]))
	fmt.Println(string(c1[:]))

	//json. NewDecoder ( strings. NewReader ( fcm ) )
}