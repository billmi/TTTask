package msg

import (
	"github.com/appleboy/go-fcm"
	"encoding/json"
	"bytes"
)

const (
	MessageBody   = "message"
	MessageNotice = "notice"
	MessageStand  = "stand"
)

const TypeTopic  = 1
const TypeToken  = 2

// Get Message from type default MessageBody
func GetMessageType(handler *handler) (types string)  {
	types = handler.GetFromKey("type")
	if types == "" {
		types = MessageBody
	}

	return types
}

// Get Message Content
func GetMessage(handler *handler, types string) *fcm.Message {

	if types == MessageNotice {
		return getNotice(handler)
	}

	if types == MessageStand {
		return getStand(handler)
	}

	return getMessage(handler)
}

// Get Message Content with type eq MessageBody
func getMessage(handler *handler) *fcm.Message {
	data := []byte(handler.GetFromKey("message"))

	var message map[string]interface{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		panic("message body error" + string(data[:]) + "! error:"+ err.Error())
	}
	return &fcm.Message{
		Data: message,
	}
}

// Get Message Content with type eq MessageNotice
func getNotice(handler *handler) *fcm.Message {
	return &fcm.Message{
		Notification: &fcm.Notification{
			Title:       handler.GetFromKeyDefault("title",handler.conf.Notification.Title),
			Body:        handler.GetFromKeyDefault("body",handler.conf.Notification.Body),
			Icon:        handler.GetFromKeyDefault("icon",handler.conf.Notification.Icon),
			ClickAction: handler.GetFromKeyDefault("click_action",handler.conf.Notification.Uri),
		},
	}
}

// Get Message Content with type eq MessageStand
func getStand(handler *handler) *fcm.Message  {
	var msg  fcm.Message
	//remove type first
	var b *bytes.Buffer
	if handler.types == FromTypeJson {
		b, _ = interfaceToByteArray(handler.raw)
	}else {
		b, _ = interfaceToByteArray(handler.from)
	}
	json.Unmarshal(b.Bytes(), msg)
	return &msg
}

