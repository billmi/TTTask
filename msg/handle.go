package msg

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"encoding/gob"
	"golang.org/x/tools/cmd/vet/testdata/divergent"
)

const (
	FromTypeNormal 	= "application/x-www-form-urlencoded"
	FromTypeMulti  	= "multipart/form-data"
	FromTypeEmpty  	= ""
	FromTypeJson	= "application/json"
)




type handler struct {
	conf    *ApiConfig
}

func GetHandle() *handler {
	return &handler{}
}

func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := initRequest(r)
	if err != nil {
		responseError(w,400)
		return
	}
	fcmMsg := &fcmMsg{}
	fcmMsg.types = GetMessageType(r)

	fcmMsg.message = GetMessage(r, fcmMsg.types)
	fcmMsg.conf = handler.conf
	fcmMsg.message.TimeToLive = &fcmMsg.conf.MaxTtl

	fcmMsg.SetUsers(r)
	response,err := fcmMsg.Send()

	if err != nil {
		responseError(w,400)
		return
	}

	if response.Failure > 0 {
		responseError(w,400)
		return
	}

	if response.Error != nil {
		responseErrorMessage(w,500, err)
	}

	w.Header().Set("Content-Type","application/json")
	w.Write([]byte("success"))
	w.WriteHeader(200)
}

func (handler *handler) SetConfig(config *ApiConfig)  {
	handler.conf = config
}

func initRequest(r *http.Request)  (err error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == FromTypeEmpty || contentType == FromTypeNormal {
		err = r.ParseForm()
	}else if contentType == FromTypeMulti {
		err = r.ParseMultipartForm(2 << 20)
	}else if contentType == FromTypeJson{
		parseJsonRequest(r)
	}else {
		err = r.ParseForm()
	}

	return err
}

func parseJsonRequest(r *http.Request) (err error) {
	result, err := ioutil.ReadAll(r.Body)
	if err == nil {
		data := make(map[string]interface{})
		err = json.Unmarshal(result, data)
		if len(data) > 0 {
			for key,value := range data {
				buff,err := interfaceToByteArray(value)
				if err != nil {
					return err
				}
				r.Form.Set(key,buff.String())

			}
		}
	}
	return err
}

func interfaceToByteArray(value interface{}) (buff bytes.Buffer, err error)  {
	enc := gob.NewEncoder(&buff)
	err = enc.Encode(value)

	return buff,err
}

func responseError(w http.ResponseWriter, code int)  {
	w.Write(GetError(code))
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
}

func responseErrorMessage(w http.ResponseWriter, code int, message interface{})  {
	b,_ := json.Marshal(fcmError{code, message})
	w.Write(b)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
}
