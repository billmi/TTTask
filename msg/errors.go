package msg

import (
	"encoding/json"
	"sync"
)

type fcmError struct {
	code int `json:"code:omitempty"`
	response interface{} `json:"response:omitempty"`
}

var errors map[int]fcmError
var onceResponse sync.Once

func initErrors()  map[int]fcmError {
	onceResponse.Do(func() {
		errors = make(map[int]fcmError)
		errors[400] = fcmError{400, "錯誤的消息體"}
		errors[401] = fcmError{401, "授權未通過"}
	})
	return errors
}

func GetError(key int) []byte {
	initErrors()

	if _,ok := errors[key]; ok {
		b,_ := json.Marshal(errors[key])
		return b
	}

	return nil
}
