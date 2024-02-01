package web_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"owl"
	"time"
)

type WebServerOptions struct {
	Domain       string `json:"domain"`
	MaxCons      int    `json:"max-cons"`
	ReadTimeout  int    `json:"read-timeout"`
	WriteTimeout int    `json:"write-timeout"`
	IdleTimeout  int    `json:"idle-timeout"`
	Mode         string `json:"mode"`
}

type WebServer struct {
	e           *gin.Engine
	middlewares []gin.HandlerFunc
	stage       *owl.Stage
	opt         *WebServerOptions
}

func NewWebServer(stage *owl.Stage, e *gin.Engine, options *WebServerOptions) *WebServer {
	return &WebServer{
		e:     e,
		stage: stage,
		opt:   options,
	}
}

func (i *WebServer) getServerAndListener(port int) (*http.Server, net.Listener) {

	address := fmt.Sprintf(":%d", port)
	maxCons := i.opt.MaxCons
	rt := time.Duration(i.opt.ReadTimeout)
	wt := time.Duration(i.opt.WriteTimeout)
	server := &http.Server{
		Addr:           address,          // 服务器监听的地址和端口
		Handler:        i.e,              // Gin 引擎作为处理器
		ReadTimeout:    rt * time.Minute, // 读取请求的超时时间
		WriteTimeout:   wt * time.Minute, // 写入响应的超时时间
		MaxHeaderBytes: 1 << 20,          // 允许的最大请求头大小
	}
	listener, err := net.Listen("tcp", address)

	if err != nil {
		panic(address + " 端口已占用，换个端口或解除端口占用")
	}

	listener = netutil.LimitListener(listener, maxCons)

	//resourcesPath := i.stage.RuntimePath(owl.ResourcesPath)
	//i.e.Static(owl.ResourcesPath, resourcesPath)

	return server, listener
}

func (i *WebServer) Use(middlewares ...gin.HandlerFunc) {
	for _, middleware := range middlewares {
		i.e.Use(middleware)
	}
}
