package msg

import (
	"github.com/appleboy/go-fcm"
	"net/http"
	"encoding/json"
)

const (
	MessageBody   = "message"
	MessageNotice = "notice"
)

func GetMessageType(r *http.Request) (types string)  {
	types = r.FormValue("type")
	if types == "" {
		types = MessageBody
	}

	return types
}
func GetMessage(r *http.Request, types string) *fcm.Message {
	if types == MessageBody {
		return getMessage(r)
	}

	return getNotice(r)
}

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


