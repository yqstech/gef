package serv

import (
	"flag"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/GoEasy/Utils/util"
	"github.com/yqstech/gef/config"
	"net/http"
	"os"
)

// Server 封装gracehttp服务器对象
type Server struct {
	ServeList    []Serve
	MuxServeList []MuxServe
}

// Serve 单个服务结构
type Serve struct {
	Name   string
	Port   int
	Router *httprouter.Router
}

type MuxServe struct {
	Name   string
	Port   int
	Router *mux.Router
}

// Set 服务器组设置一个web服务
func (that *Server) Set(name string, port interface{}, router *httprouter.Router) {
	that.ServeList = append(that.ServeList, Serve{
		Name:   name,
		Port:   util.String2Int(util.Interface2String(port)),
		Router: router,
	})
}
func (that *Server) SetMux(name string, port interface{}, router *mux.Router) {
	that.MuxServeList = append(that.MuxServeList, MuxServe{
		Name:   name,
		Port:   util.String2Int(util.Interface2String(port)),
		Router: router,
	})
}

// Run 启动服务器组
func (that *Server) Run() {
	if that.ServeList == nil {
		logger.Error("请先设置web服务器!")
		return
	}
	//格式化服务器组
	servers := []*http.Server{}
	for _, server := range that.ServeList {
		TServ := flag.String(server.Name, ":"+util.Int2String(server.Port), server.Name)
		servers = append(servers,
			&http.Server{
				Addr:    *TServ,
				Handler: server.Router,
			},
		)
		logger.Alert(server.Name+" running by port ", server.Port)
	}
	for _, server := range that.MuxServeList {
		TServ := flag.String(server.Name, ":"+util.Int2String(server.Port), server.Name)
		servers = append(servers,
			&http.Server{
				Addr:    *TServ,
				Handler: server.Router,
			},
		)
		logger.Alert(server.Name+"(mux) running by port ", server.Port)
	}

	//启动服务器组
	err := gracehttp.Serve(servers...)
	if err != nil {
		logger.Error("GoEasy* running error:", err.Error())
		return
	}
}

// Args 获取命令行数据
func Args() string {
	//设置默认端口
	port := config.FrontPort
	//命令行参数配置
	args := os.Args
	if len(args) > 1 {
		port = args[1]
	}
	return port
}
