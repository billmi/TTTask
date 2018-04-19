package msg

import (
	"github.com/appleboy/go-fcm"
	"net/http"
	"encoding/json"
)

const (
	MessageBody   = "message"
	MessageNotice = "notice"
	MessageStand  = "stand"
)


// Get Message from type default MessageBody
func GetMessageType(r *http.Request) (types string)  {
	types = r.FormValue("type")
	if types == "" {
		types = MessageBody
	}

	return types
}

// Get Message Content
func GetMessage(r *http.Request, types string) *fcm.Message {

	if types == MessageNotice {
		return getNotice(r)
	}

	if types == MessageStand {
		return getStand(r)
	}

	return getMessage(r)
}

// Get Message Content with type eq MessageBody
func getMessage(r *http.Request) *fcm.Message {
	data := []byte(r.FormValue("message"))
	var message map[string]interface{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		panic("message body error" + string(data[:]))
	}
	return &fcm.Message{
		Data: message,
	}
}

// Get Message Content with type eq MessageNotice
func getNotice(r *http.Request) *fcm.Message {
	return &fcm.Message{
		Notification: &fcm.Notification{
			Title:       r.FormValue("title"),
			Body:        r.FormValue("body"),
			Icon:        r.FormValue("icon"),
			ClickAction: r.FormValue("click_action"),
		},
	}
}

// Get Message Content with type eq MessageStand
func getStand(r *http.Request) *fcm.Message  {
	var msg  fcm.Message
	//remove type first
	r.Form.Del("type")
	b,_ := interfaceToByteArray(r.Form)
	json.Unmarshal(b.Bytes(), msg)
	return &msg
}

