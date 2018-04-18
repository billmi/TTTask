package msg

import (
	"net/http"
	"encoding/json"
)

const (
	FromTypeNormal = "application/x-www-form-urlencoded"
	FromTypeMulti  = "multipart/form-data"
	FromTypeEmpty  = ""
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
	}else{
		err = r.ParseForm()
	}
	return err
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