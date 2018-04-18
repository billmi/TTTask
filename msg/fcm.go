package msg

import (
	"github.com/appleboy/go-fcm"
	"reflect"
)


type fcmMsg struct {
	types   string
	message *fcm.Message
	conf    *ApiConfig
}

func GetNewFcm() *fcmMsg  {
	return &fcmMsg{}
}

func (fcmMsg *fcmMsg) SetTopic(topics interface{}) {
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

func (fcmMsg *fcmMsg) SetConfig(config *ApiConfig) {
	fcmMsg.conf = config
}

func (fcmMsg *fcmMsg) SetTo(token string) {
	fcmMsg.message.To = token
}

func (fcmMsg *fcmMsg) SetMessage(msg *fcm.Message) {
	fcmMsg.message = msg
}

func (fcmMsg *fcmMsg) Send() (*fcm.Response,error)  {
	client, err := fcm.NewClient(fcmMsg.conf.GetKey())
	if err != nil {
		return nil,err
	}

	response, err := client.Send(fcmMsg.message)

	return response,err
}

func getMultiTopicsArray(arr []string) (str string) {
	for _,v := range arr {
		str += v + " && "
	}
	return str
}
