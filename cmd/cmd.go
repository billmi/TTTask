package cmd

import (
	"flag"
)

type command struct {
	list map[string]string
	alias map[string]string
}


func GetCommand() *command {
	c := &command{
		make(map[string]string),
		make(map[string]string),
	}
	c.init()
	return c
}


func (command *command) init()  {
	path := ""
	help := ""
	daemon := false
	flag.StringVar(&help, "h", "", "")
	flag.StringVar(&path, "c", "", "")
	flag.BoolVar(&daemon, "d", false, "")
	flag.Parse()

	if path == "" {
		path=getCurrentDirectory()+"/conf/config.json"
	}

	if daemon {
		command.SetAlias("d","daemon","1")
	}else{
		command.SetAlias("d","daemon","0")
	}

	command.SetAlias("h","help",help)
	command.SetAlias("c","config",path)
}

func (command *command) Get(key string) string {
	if _, ok := command.list[key]; ok {
		return command.list[key]
	}else{
		if _, ok := command.alias[key]; ok {
			key = command.alias[key]
		}
		return command.Get(key)
	}

	return ""
}

func (command *command) Set(key string, value string)  {
	command.list[key] = value
}

func (command *command) SetAlias(key string, alias string, value string)  {
	command.Set(key, value)
	command.alias[alias] = key
}

