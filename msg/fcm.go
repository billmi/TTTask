package msg

import (
	"github.com/appleboy/go-fcm"
	"reflect"
)


type FcmMsg struct {
	messageId int
	types   string
	message *fcm.Message
	conf    *ApiConfig
	sendTime int64
}

func GetNewFcm() *FcmMsg  {
	return &FcmMsg{}
}

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

func (fcmMsg *FcmMsg) SetSendTime(sendTime int64)  {
	fcmMsg.sendTime = sendTime
}
func (fcmMsg *FcmMsg) SetMessageId(id int) {
	fcmMsg.messageId = id
}

func (fcmMsg *FcmMsg) SetConfig(config *ApiConfig) {
	fcmMsg.conf = config
}

func (fcmMsg *FcmMsg) SetTo(token string) {
	fcmMsg.message.To = token
}

func (fcmMsg *FcmMsg) SetTtl(ttl *uint) {
	fcmMsg.message.TimeToLive = ttl
}

func (fcmMsg *FcmMsg) SetCondition(condition string) {
	fcmMsg.message.Condition = condition
}

func (fcmMsg *FcmMsg) SetMessage(msg *fcm.Message) {
	fcmMsg.message = msg
}

func (fcmMsg *FcmMsg) Task()  {
	AddToTask(fcmMsg)
}

func (fcmMsg *FcmMsg) Send() (*fcm.Response,error)  {

	client, err := fcm.NewClient(fcmMsg.conf.GetKey())
	if err != nil {
		return nil,err
	}

	response, err := client.Send(fcmMsg.message)

	return response,err
}

func (fcmMsg *FcmMsg) SetUsers(handler *handler) (err error)  {

	if fcmMsg.types == MessageStand {
		return nil
	}

	condition := handler.GetFromKey("condition")
	if condition != "" {
		fcmMsg.SetCondition(condition)
	}

	token := handler.GetFromKey("token")
	if token != "" {
		fcmMsg.SetTo(token)
	}

	topic := handler.GetFromKey("topic")
	fcmMsg.SetTopic(topic)
	//if topic != "" {
	//	var interfaceTopic interface{}
	//	fcmMsg.SetTopic(json.Unmarshal(bytes.NewBufferString(topic).Bytes(), interfaceTopic))
	//}

	return err
}

func getMultiTopicsArray(arr []string) (str string) {
	for _,v := range arr {
		str += v + " && "
	}
	return str
}
