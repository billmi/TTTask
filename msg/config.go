package msg

import (
	"github.com/jingweno/conf"
)

type ApiConfig struct {
	ApiKey     string
	ServerPort string
	MaxTtl     uint
	Icon	   string
	Uri        string
}

func GetConfig(path string) *ApiConfig {

	c, err := conf.NewLoader().Argv().Env().File(path).Load()

	if err != nil {
		panic("load config err")
	}

	return &ApiConfig{
		ApiKey:     c.Get("api_key").(string),
		ServerPort: c.Get("server_port").(string),
		MaxTtl:     uint(c.Get("max_ttl").(float64)),
		Icon:       c.Get("icon").(string),
		Uri:        c.Get("uri").(string),
	}
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
