package msg

import (
	"github.com/jingweno/conf"
	"fmt"
	"log"
)


type ApiConfig struct {
	ApiKey string
	ServerPort string
	MaxTtl uint
}

func GetConfig(path string) *ApiConfig {

	c, err := conf.NewLoader().Argv().Env().File(path).Load()

	if err != nil {
		fmt.Println(err)
		log.Fatal("load config err")
		panic("load config err")
	}

	return &ApiConfig{
		ApiKey:c.Get("api_key").(string),
		ServerPort:c.Get("server_port").(string),
		MaxTtl:uint(c.Get("max_ttl").(float64)),
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

