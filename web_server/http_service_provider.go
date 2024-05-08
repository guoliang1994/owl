package web_server

import (
	"github.com/gin-gonic/gin"
	"net"
	"owl"
	"owl/log"
)

type HttpOptions struct {
	*WebServerOptions
	Port int `json:"port"`
}

func NewHttpOptionFromConfigFile(cfgManager *owl.ConfManager, cfgFile string) (opt *HttpOptions) {

	err := cfgManager.GetConfig(cfgFile, &opt)
	if err != nil {
		return nil
	}
	err = cfgManager.GetConfig(cfgFile+".http", &opt)
	if err != nil {
		return nil
	}
	return opt
}

func NewDefaultHttpOption() *HttpOptions {
	opt := &HttpOptions{
		WebServerOptions: &WebServerOptions{
			Domain:       "",
			MaxCons:      1024,
			ReadTimeout:  100,
			WriteTimeout: 100,
			IdleTimeout:  100,
			Mode:         "release",
		},
		Port: 80,
	}
	return opt
}

type HttpService struct {
	opt *HttpOptions
	*WebServer
}

func NewHttpService(stage *owl.Stage, e *gin.Engine, opt *HttpOptions) *HttpService {
	if opt == nil {
		opt = NewDefaultHttpOption()
	}
	httpServer := &HttpService{
		opt:       opt,
		WebServer: NewWebServer(stage, e, opt.WebServerOptions),
	}
	httpServer.BlockRun()
	return httpServer
}

func (i *HttpService) BlockRun() {
	server, listener := i.getServerAndListener(i.opt.Port)
	i.opt.Port = listener.Addr().(*net.TCPAddr).Port
	log.PrintLnBlue("http server start on port:", i.opt.Port)
	go func() {

		err := server.Serve(listener)
		if err != nil {
			log.PrintRed("http server start on port:", i.opt.Port)
		}
	}()
}
func (i *HttpService) GetOptions() *HttpOptions {
	return i.opt
}
