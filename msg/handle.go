package msg

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"encoding/gob"
	"strconv"
	"net/url"
	"reflect"
	"log"
)

const (
	FromTypeNormal = "application/x-www-form-urlencoded"
	FromTypeMulti  = "multipart/form-data"
	FromTypeEmpty  = ""
	FromTypeJson   = "application/json"
)

//handler data
type handler struct {
	conf  *ApiConfig
	types string         //FromType
	from  url.Values
	raw   map[string]interface{}
}

func GetHandle() *handler {
	return &handler{}
}
//hand http server
//if want to any logic
//add path and func
func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType, content, err := initRequest(r)
	handler.types = contentType
	if contentType == FromTypeJson {
		handler.raw = content
	} else {
		handler.from = r.Form
	}

	if err != nil {
		responseError(w, 400)
		return
	}

	fcmMsg := GetNewFcm()

	fcmMsg.types = GetMessageType(handler)

	//設置數據庫綁定消息的ID

	id, err := strconv.Atoi(handler.GetFromKeyDefault("id", "0"))
	if err != nil {
		responseErrorMessage(w, 500, err.Error())
		return
	}
	fcmMsg.SetMessageId(id)

	//設置發送時間
	sendTime, err := strconv.ParseInt(handler.GetFromKeyDefault("send_time", "0"), 10, 64)
	if err != nil {
		responseErrorMessage(w, 500, err.Error())
		return
	}
	r.Form.Del("send_time")
	fcmMsg.SetSendTime(sendTime)

	//設置消息體和配置
	fcmMsg.SetMessage(GetMessage(handler, fcmMsg.types))
	fcmMsg.SetConfig(handler.conf)
	fcmMsg.SetTtl(&handler.conf.MaxTtl)
	fcmMsg.SetUsers(handler)

	msgId := ""
	if fcmMsg.sendTime > 0 {
		fcmMsg.Task()
		responseSuccess(w,"")
	} else{
		//發送消息
		//response : https://firebase.google.com/docs/cloud-messaging/http-server-ref?hl=zh-cn#error-codes
		//multicast_id	必选，数字	用于识别多播消息的唯一 ID（数字）。
		//success	必选，数字	处理时未出错的消息数。
		//failure	必选，数字	无法处理的消息数。
		//canonical_ids	必选，数字	包含规范注册令牌的结果数。
		//规范注册 ID 是客户端应用最后一次请求注册时的注册令牌。这是服务器在向设备发送消息时应当使用的 ID。
		//results
		//      message_id：字符串，用于指定每条成功处理的消息的唯一 ID。
		//		registration_id：可选字符串，用于指定处理和接收消息的客户端应用的规范注册令牌。发送者应使用此值作为未来请求的注册令牌。否则，消息可能被拒绝。
		//		error：字符串，用于指定处理收件人的消息时发生的错误。表 9 中列出了可能的值。
		response, err := fcmMsg.Send()
		log.Println(response,err)
		if err != nil {
			responseErrorMessage(w, 500, err.Error())
			return
		}

		if response.Error != nil {
			responseErrorMessage(w, 500, err.Error())
			return
		}

		log.Println(response)

		if fcmMsg.sendType == TypeToken {
			if response.Failure > 0 {
				responseError(w, 400)
				return
			}
		}

		//todo: topics fail response
		//if fcmMsg.sendType == TypeTopic {
		//
		//}
		responseSuccess(w,msgId)
	}


}
//GetFromKey in json or url.Values(alias r.From)
func (handler *handler) GetFromKey(key string) string {
	if handler.types == FromTypeJson {
		return handler.getRawData(key)
	} else {
		v := handler.from.Get(key)
		if v == "" {
			v = handler.getRawData(key)
		}
		return v
	}
}
//GetFromKey in json or url.Values(alias r.From)
//if empty return default
func (handler *handler) GetFromKeyDefault(key string, defaults string) string {
	var v string
	if handler.types == FromTypeJson {
		v = handler.getRawData(key)
	} else {
		v = handler.from.Get(key)
		if v == "" {
			v = handler.getRawData(key)
		}
	}
	if v == "" {
		return defaults
	}
	return v
}

func (handler *handler) getRawData(key string) string  {
	v := ""
	if _, ok := handler.raw[key]; ok {
		t := reflect.ValueOf(handler.raw[key]).Kind()
		if t != reflect.String {
			s, _ := json.Marshal(handler.raw[key])
			v = string(s[:])
		}else {
			v = handler.raw[key].(string)
		}
	}
	return v
}

//set config
func (handler *handler) SetConfig(config *ApiConfig) {
	handler.conf = config
}
//getContentType and parseFrom or RawData
func initRequest(r *http.Request) (contentType string, content map[string]interface{}, err error) {
	contentType = r.Header.Get("Content-Type")

	if contentType == FromTypeEmpty || contentType == FromTypeNormal {
		err = r.ParseForm()
	} else if contentType == FromTypeMulti {
		err = r.ParseMultipartForm(2 << 20)
	} else if contentType == FromTypeJson {
		content, err := parseJsonRequest(r)
		return contentType, content, err
	} else {
		err = r.ParseForm()
	}

	if len(r.Form) == 1 {
		contentType = FromTypeJson
		for key, _ := range r.Form {
			content, err := parseJsonByte(bytes.NewBufferString(key).Bytes())
			return contentType, content, err
		}
	}

	return contentType, nil, err
}
//Parse Json from *http.Request.Body
func parseJsonRequest(r *http.Request) (map[string]interface{}, error) {
	result, err := ioutil.ReadAll(r.Body)
	if err == nil {
		return parseJsonByte(result)
	}
	return nil, err
}
//Parse Json from []byte
func parseJsonByte(b []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
//interface to []byte
func interfaceToByteArray(value interface{}) (buff *bytes.Buffer, err error) {
	bu := make([]byte, 0)
	buff = bytes.NewBuffer(bu)
	enc := gob.NewEncoder(buff)
	err = enc.Encode(value)

	return buff, err
}

func responseError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(GetError(code))
}

func responseErrorMessage(w http.ResponseWriter, code int, message interface{}) {
	b, _ := json.Marshal(fcmError{code, message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func responseSuccess(w http.ResponseWriter, msgId string)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(GetSuccess(msgId))
}
//func getFromValue(r *http.Request, key string, defaultValue interface{}) interface{} {
//	value := r.Form.Get(key)
//	if value == "" {
//		return defaultValue
//	}
//	return value
//}
