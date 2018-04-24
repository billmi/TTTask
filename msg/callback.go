package msg

import (
	"net/http"
	"encoding/json"
	"bytes"
	"log"
	"os"
	"time"
)

type callback struct {
	MsgId string `json:"msg_id"`
	Status uint8 `json:"status"`
}

type CallBack struct {
	callback *callback
	fcmResponse *fcmResponse
	config *ApiConfig
}

func (callback *callback) GetCallBackResponse() *fcmResponse {
	code := 0

	if callback.Status != 200 {
		code = 500
	}

	return &fcmResponse{
		Code: code,
		Response: callback,
	}
}

func GetCallBack(msgId string, status int) *CallBack {
	callback:=&callback{
		MsgId:  msgId,
		Status: uint8(status),
	}
	return &CallBack{
		callback:callback,
		fcmResponse: callback.GetCallBackResponse(),
	}
}

func (CallBack *CallBack) SetConfig(config *ApiConfig) {
	CallBack.config = config
}

func (CallBack *CallBack) Do() {
	b,_ := json.Marshal(CallBack.fcmResponse)
	response,err := http.Post(CallBack.config.CallBack, FromTypeJson, bytes.NewReader(b))

	if err != nil {
		//todo: save log
	}else{
		t := time.Now()
		fileName := t.Format("")
		file, err := os.Create(CallBack.config.LogFile+fileName+".log")
		if err != nil {
			log.Fatalln("fail to create test.log file!")
		}
		logger := log.New(file, "", log.LstdFlags|log.Llongfile)
		logger.Println(response)
	}


}