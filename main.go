/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: gef.go
 * @Version: 1.0.0
 * @Date: 2022/6/26 12:21
 */

package gef

import (
	"embed"
	"encoding/gob"
	"errors"
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Event"
	"github.com/gef/GoEasy/Registry"
	"github.com/gef/GoEasy/Routers"
	"github.com/gef/GoEasy/Templates"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/pool"
	"github.com/gef/GoEasy/Utils/serv"
	"github.com/gef/config"
	"github.com/gef/routers"
	"github.com/gef/rpcPlugins"
	"github.com/gef/static"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"os"
)

func init() {
	//! 初始化配置信息，自动读取配置文件
	err := config.Init()
	if err != nil {
		panic(err)
	}
	//! 查找log配置文件，设置日志配置
	if _, err = os.Stat(config.WorkPath + "/configs/log.json"); err == nil {
		logger.SetLogger(config.WorkPath + "/configs/log.json")
	} else {
		if _, err = os.Stat(config.AppPath + "/configs/log.json"); err == nil {
			logger.SetLogger(config.AppPath + "/configs/log.json")
		} else {
			panic(errors.New("获取日志配置信息失败！"))
		}
	}
	//! 初始化数据库
	db.Init()
	//! 初始化Redis
	pool.RedisInit()
	//! 初始化GoCache
	pool.GocacheInit()
	//! rpc插件需要注册一个类型
	gob.Register(map[string]interface{}{})
}

// New 创建新的Gef应用
func New() Gef {
	gef := Gef{}
	return gef
}

// Gef Gef应用结构
type Gef struct {
	Servers      []Server //服务器列表
	selfServers  []Server //自定义服务列表
	StaticFile   embed.FS //静态资源文件
	TemplateFile embed.FS //模板文件
}

// Server web服务结构
type Server struct {
	Name       string      //名称
	Port       string      //端口
	Router     interface{} //路由
	RouterType int         //路由类型 0httprouter 1 mux
}

// SetAdminPages 设置后台页面
func (g *Gef) SetAdminPages(pages map[string]EasyApp.AppPage) {
	for k, v := range pages {
		Registry.AdminPages[k] = v
	}
}

// SetFrontRouters 设置前台路由
func (g *Gef) SetFrontRouters(FrontRouters interface{}) {
	routers.FrontRouters = FrontRouters
}

// SetFrontRouterType 设置前台路由
func (g *Gef) SetFrontRouterType(routerType int) {
	routers.FrontRouterType = routerType
}

// SetEvent 补充监听事件
func (g *Gef) SetEvent(EventAdd map[string][]Event.Listener) {
	Event.BindEvents(EventAdd)
}

// SetAdminStatic 设置静态文件
func (g *Gef) SetAdminStatic(f embed.FS) {
	static.FilesAdd = f
}

// SetAdminTemplate 设置静态文件
func (g *Gef) SetAdminTemplate(f embed.FS) {
	Templates.FilesAdd = f
}

// SetServer 设置服务器
func (g *Gef) SetServer(serv Server) {
	g.selfServers = append(g.selfServers, serv)
}

// Run 启动Gef应用
func (g *Gef) Run() {
	//! 设置web服务组
	if config.AdminPort != "" {
		g.Servers = append(g.Servers, Server{
			Name:       "gef-admin",
			Port:       config.AdminPort,
			Router:     Routers.AdminRouters(),
			RouterType: 0,
		})
	}
	if config.FrontPort != "" && routers.FrontRouters != nil {
		g.Servers = append(g.Servers, Server{
			Name:       "gef-front",
			Port:       config.FrontPort,
			Router:     routers.FrontRouters,
			RouterType: routers.FrontRouterType,
		})
	}
	for _, server := range g.selfServers {
		g.Servers = append(g.Servers, server)
	}
	
	//!退出web服务清除RPC链接
	defer rpcPlugins.RpcPluginClientsKill()
	
	//!启动web服务组
	HttpServers := serv.Server{}
	for _, serv := range g.Servers {
		if serv.RouterType == 0 {
			HttpServers.Set(serv.Name, serv.Port, serv.Router.(*httprouter.Router))
		} else if serv.RouterType == 1 {
			HttpServers.SetMux(serv.Name, serv.Port, serv.Router.(*mux.Router))
		}
	}
	HttpServers.Run()
}
