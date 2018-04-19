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
)

const (
	FromTypeNormal = "application/x-www-form-urlencoded"
	FromTypeMulti  = "multipart/form-data"
	FromTypeEmpty  = ""
	FromTypeJson   = "application/json"
)

type handler struct {
	conf  *ApiConfig
	types string
	from  url.Values
	raw   map[string]interface{}
}

func GetHandle() *handler {
	return &handler{}
}

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

	if fcmMsg.sendTime > 0 {
		fcmMsg.Task()
	} else{
		//發送消息
		response, err := fcmMsg.Send()

		if err != nil {
			responseError(w, 400)
			return
		}

		if response.Failure > 0 {
			responseError(w, 400)
			return
		}

		if response.Error != nil {
			responseErrorMessage(w, 500, err.Error())
			return
		}
	}

	responseSuccess(w)
}

func (handler *handler) GetFromKey(key string) string {
	if handler.types == FromTypeJson {
		if _, ok := handler.raw[key]; ok {
			t := reflect.ValueOf(handler.raw[key]).Kind()
			if t != reflect.String {
				s, _ := json.Marshal(handler.raw[key])
				return string(s[:])
			}
			return handler.raw[key].(string)
		}
	} else {
		return handler.from.Get(key)
	}
	return ""
}

func (handler *handler) GetFromKeyDefault(key string, defaults string) string {
	var v string
	if handler.types == FromTypeJson {
		if _, ok := handler.raw[key]; ok {
			t := reflect.ValueOf(handler.raw[key]).Kind()
			if t != reflect.String {
				s, _ := json.Marshal(handler.raw[key])
				v = string(s[:])
			}else {
				v = handler.raw[key].(string)
			}
		}
	} else {
		v = handler.from.Get(key)

	}
	if v == "" {
		return defaults
	}
	return v
}

func (handler *handler) SetConfig(config *ApiConfig) {
	handler.conf = config
}

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

func parseJsonRequest(r *http.Request) (map[string]interface{}, error) {
	result, err := ioutil.ReadAll(r.Body)
	if err == nil {
		return parseJsonByte(result)
	}
	return nil, err
}

func parseJsonByte(b []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

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

func responseSuccess(w http.ResponseWriter)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(GetError(200))
}
//func getFromValue(r *http.Request, key string, defaultValue interface{}) interface{} {
//	value := r.Form.Get(key)
//	if value == "" {
//		return defaultValue
//	}
//	return value
//}
