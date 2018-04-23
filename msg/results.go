package msg

import (
	"encoding/json"
	"sync"
)

//response body content
//like {"code":200,"response":"success"}
type fcmError struct {
	Code int `json:"code,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

var errorsMap map[int]fcmError
var onceResponse sync.Once

//init once
func initErrors()  map[int]fcmError {
	onceResponse.Do(func() {
		errorsMap = make(map[int]fcmError)
		errorsMap[400] = fcmError{400, "錯誤的消息體"}
		errorsMap[401] = fcmError{401, "授權未通過"}
		errorsMap[200] = fcmError{0,"success"}
	})
	return errorsMap
}

//get error []byte to write
func GetError(key int) []byte {
	initErrors()

	if _,ok := errorsMap[key]; ok {
		b,_ := json.Marshal(errorsMap[key])
		return b
	}

	return nil
}

func GetSuccess(msgId string) []byte {

	if msgId == "" {
		return GetSuccessAsync()
	}

	return GetSuccessSync(msgId)
}
func GetSuccessSync(msgId string) []byte {
	f := errorsMap[200]
	type response struct{
		MsgId string `json:"msg_id"`
	}
	f.Response = response{MsgId:msgId}
	b,_ := json.Marshal(f)
	return b
}

func GetSuccessAsync() []byte {
	b,_ := json.Marshal(errorsMap[200])
	return b
}
