package msg

import (
	"encoding/json"
	"sync"
)

//response body content
//like {"code":200,"response":"success"}
type fcmResponse struct {
	Code     int         `json:"code,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

var errorsMap map[int]fcmResponse
var onceResponse sync.Once

//init once
func initResponse() map[int]fcmResponse {
	onceResponse.Do(func() {
		errorsMap = make(map[int]fcmResponse)
		errorsMap[400] = fcmResponse{400, "錯誤的消息體"}
		errorsMap[401] = fcmResponse{401, "授權未通過"}
		errorsMap[200] = fcmResponse{200, "success"}
		errorsMap[406] = fcmResponse{406, "Method Not Allowed"}
	})
	return errorsMap
}

//get error []byte to write
func GetError(key int) []byte {
	initResponse()

	if _, ok := errorsMap[key]; ok {
		b, _ := json.Marshal(errorsMap[key])
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
	type response struct {
		MsgId string `json:"msg_id"`
	}
	f.Code = 200
	f.Response = response{MsgId: msgId}
	b, _ := json.Marshal(f)
	return b
}

func GetSuccessAsync() []byte {
	f := errorsMap[200]
	f.Code = 200
	b, _ := json.Marshal(f)
	return b
}
