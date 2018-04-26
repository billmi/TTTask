package msg

import (
	"github.com/jingweno/conf"
)

type DefaultNotification struct {
	Title string
	Body  string
	Icon  string
	Uri   string
}

//load from config.json
//ApiKey https://firebase.google.com/docs/server/setup#prerequisites
//ServerPort http listen server
//MaxTtl https://firebase.google.com/docs/cloud-messaging/concept-options#ttl
//Notification
//normal notification content if no post value
type ApiConfig struct {
	ApiKey       string
	ServerPort   string
	MaxTtl       uint
	CallBack     string
	LogFile      string
	Proxy		 string
	Notification *DefaultNotification
	Maps 		 map[string]string
}

//load config from path
//@return *ApiConfig
func GetConfig(path string) *ApiConfig {

	c, err := conf.NewLoader().Argv().Env().File(path).Load()

	if err != nil {
		panic("load config err "+path)
	}

	maps := c.Get("notification").(map[string]interface{})
	notification := &DefaultNotification{}
	for key, value := range maps {
		if key == "uri" {
			notification.Uri = value.(string)
			continue
		}

		if key == "body" {
			notification.Body = value.(string)
			continue
		}

		if key == "icon" {
			notification.Icon = value.(string)
			continue
		}

		if key == "title" {
			notification.Title = value.(string)
			continue
		}
	}
	return &ApiConfig{
		ApiKey:       c.String("api_key"),
		ServerPort:   c.String("server_port"),
		MaxTtl:       uint(c.Get("max_ttl").(float64)),
		Notification: notification,
		CallBack:     c.String("notify_callback"),
		LogFile:      c.String("log_file"),
		Proxy:        c.String("proxy"),
	}
}

//need todo
func (conf *ApiConfig) Get(key string) string {
	return ""
}

func (conf *ApiConfig) GetKey() string {
	return conf.ApiKey
}

func (conf *ApiConfig) GetServerPort() string {
	return conf.ServerPort
}

func (conf *ApiConfig) GetMaxTtl() uint {
	return conf.MaxTtl
}
