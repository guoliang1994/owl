package web_server

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net"
	"owl"
	"owl/log"
)

type Options struct {
	*WebServerOptions
	HttpPort int `json:"http-port"`
}

func NewOption(stage *owl.Stage) (opt *Options) {
	err := jsoniter.UnmarshalFromString(stage.ConfManager.GetConfig("app").ToString(), &opt)
	if err != nil {
		return nil
	}
	return opt
}

type HttpService struct {
	opt *Options
	*WebServer
}

func NewHttpService(stage *owl.Stage, e *gin.Engine, options *Options) *HttpService {
	return &HttpService{
		opt:       options,
		WebServer: NewWebServer(stage, e, options.WebServerOptions),
	}
}

func (i *HttpService) Boot() {
	server, listener := i.getServerAndListener(i.opt.HttpPort)
	i.opt.HttpPort = listener.Addr().(*net.TCPAddr).Port
	log.PrintLnBlue("http server start on port:", i.opt.HttpPort)
	go func() {

		err := server.Serve(listener)
		if err != nil {
			log.PrintRed("http server start on port:", i.opt.HttpPort)
		}
	}()
}
