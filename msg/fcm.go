package msg

import (
	"github.com/appleboy/go-fcm"
	"reflect"
	"errors"
)

// FcmMsg struct
type FcmMsg struct {
	messageId int        // relation sql save id  or other id
	types   string       //  message type  message|notice|stand
	message *fcm.Message  // message json obj
	conf    *ApiConfig   //  serverConfig
	sendTime int64      //  if sendTime > now()   add to cronTab
	sendType int
}

func GetNewFcm() *FcmMsg  {
	return &FcmMsg{}
}

//set fcm receive topics
func (fcmMsg *FcmMsg) SetTopic(topics interface{}) {
	kind := reflect.ValueOf(topics).Type().Kind()
	switch kind {
		case reflect.String:
			fcmMsg.message.To = topics.(string)
			break
		case reflect.Array:
			fcmMsg.message.Condition = getMultiTopicsArray(topics.([]string))
			break
		//case reflect.Map:
		//	fcmMsg.message.Condition = getMultiTopicsMap(topics.(map[string]interface{}))
		default:
			fcmMsg.message.To = topics.(string)
			break
	}
}
//set fcm send Time
func (fcmMsg *FcmMsg) SetSendTime(sendTime int64)  {
	fcmMsg.sendTime = sendTime
}
//set fcm relation id
func (fcmMsg *FcmMsg) SetMessageId(id int) {
	fcmMsg.messageId = id
}
//set server config for get default value(icon maxttl click_action)
func (fcmMsg *FcmMsg) SetConfig(config *ApiConfig) {
	fcmMsg.conf = config
}
//set fcm receive tokens
func (fcmMsg *FcmMsg) SetTo(token string) {
	fcmMsg.message.To = token
}
//set fcm save google server ttl if user not online
func (fcmMsg *FcmMsg) SetTtl(ttl *uint) {
	fcmMsg.message.TimeToLive = ttl
}
//set fcm receive with topics condition
func (fcmMsg *FcmMsg) SetCondition(condition string) {
	fcmMsg.message.Condition = condition
}
//set fcm message body
func (fcmMsg *FcmMsg) SetMessage(msg *fcm.Message) {
	fcmMsg.message = msg
}
//if sendTime > now() Add to cron
func (fcmMsg *FcmMsg) Task()  {
	AddToTask(fcmMsg)
}
//try to send google
func (fcmMsg *FcmMsg) Send() (*fcm.Response,error)  {

	client, err := fcm.NewClient(fcmMsg.conf.GetKey())
	if err != nil {
		return nil,err
	}

	response, err := client.Send(fcmMsg.message)

	return response,err
}
//packaging set fcm receiver
func (fcmMsg *FcmMsg) SetUsers(handler *handler) (err error)  {

	if fcmMsg.types == MessageStand {
		return nil
	}

	condition := handler.GetFromKey("condition")
	if condition != "" {
		fcmMsg.sendType = TypeTopic
		fcmMsg.SetCondition(condition)
		return nil
	}

	token := handler.GetFromKey("token")
	if token != "" {
		fcmMsg.sendType = TypeToken
		fcmMsg.SetTo(token)
		return nil
	}

	topic := handler.GetFromKey("topic")
	if topic != "" {
		fcmMsg.sendType = TypeTopic
		fcmMsg.SetTopic(topic)
		return nil
	}

	return errors.New("not receive send user")
}
//topics merge to string
func getMultiTopicsArray(arr []string) (str string) {
	for _,v := range arr {
		str += v + " && "
	}
	return str
}
