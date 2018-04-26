package cmd


import (
	"log"
	"os"
	"net/http"
	"github.com/takama/daemon"
	"github.com/lwl1989/TTTask/msg"
	"os/signal"
	"syscall"
	"fmt"
)



var stdlog *log.Logger

// Service has embedded daemon
type Service struct {
	Daemon daemon.Daemon
	Config *msg.ApiConfig
	Server *http.Server
	start chan bool
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	service.start = make(chan bool)
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	usage := "Usage: command install | remove | start | stop | status"
	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Daemon.Install()
		case "remove":
			return service.Daemon.Remove()
		case "start":
			return service.Daemon.Start()
		case "stop":
			return service.Daemon.Stop()
		case "status":
			return service.Daemon.Status()
		default:
			return usage, nil
		}
	}


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	go service.DoStart()

	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			if service.Server != nil {
				service.Server.Close()
			}

			if killSignal == os.Interrupt {
				return "Server exit", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

func (service *Service) DoStart() {
	service.Server = getServer(service.Config)
	err := service.Server.ListenAndServe()

	fmt.Println(err)
}

func getServer(config *msg.ApiConfig)  *http.Server {
	handler := msg.GetHandle()
	handler.SetConfig(config)

	server := &http.Server{
		Addr: ":"+config.GetServerPort(),
		Handler:handler,
	}
	return server
}

