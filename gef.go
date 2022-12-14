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
	"errors"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Event"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/gdb"
	"github.com/yqstech/gef/Utils/pool"
	"github.com/yqstech/gef/Utils/serv"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/builder/adminTemplates"
	"github.com/yqstech/gef/config"
	"github.com/yqstech/gef/registry"
	"github.com/yqstech/gef/routers"
	"github.com/yqstech/gef/static"
	"github.com/yqstech/gef/util"
	"net/http"
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
	//! 初始化Redis
	pool.RedisInit()
	//! 初始化GoCache
	pool.GocacheInit()

	//! 初始化数据库
	db.Init()
	gdb.Init()
	//! 自动维护数据库
	dbm := DbManager{}
	dbm.AutoTable(tables)
	dbm.AutoAdminRules(adminRules)
	dbm.AutoInsideData(insideData)
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

// SetSoftInfo 设置软件信息
func (g *Gef) SetSoftInfo(key, version string) {
	config.SoftKey = key
	config.SoftVersion = version
}

// SetLicenceResult 设置软件版权信息
func (g *Gef) SetLicenceResult(result map[string]interface{}) {
	config.LicenceResult = result
}

// SetAdminPages 设置后台页面，支持追加
func (g *Gef) SetAdminPages(pages map[string]builder.NodePager) {
	for k, v := range pages {
		registry.AdminPages[k] = v
	}
}

// SetFrontRouters 设置前台路由，追加使用【AddFrontRouters】方法
func (g *Gef) SetFrontRouters(FrontRouters interface{}) {
	routers.FrontRouters = FrontRouters
}

// SetFrontRouterType 设置前台路由类型 0 httprouter 1 mux
func (g *Gef) SetFrontRouterType(routerType int) {
	routers.FrontRouterType = routerType
}

type FrontRouter struct {
	Method             string //GET POST ServeFiles
	Url                string //访问路径
	HandleOrFileSystem interface{}
}

// AddFrontRouters 追加前台路由
func (g *Gef) AddFrontRouters(FrontRouters []FrontRouter) {
	//判断前台路由类型
	if routers.FrontRouterType == 0 {
		if routers.FrontRouters == nil {
			//假如前台路由为空，则设置
			routers.FrontRouters = httprouter.New()
		}
		for _, fr := range FrontRouters {
			if fr.Method == "ServeFiles" {
				routers.FrontRouters.(*httprouter.Router).ServeFiles(fr.Url, fr.HandleOrFileSystem.(http.FileSystem))
			}
			if fr.Method == "GET" || fr.Method == "POST" {
				routers.FrontRouters.(*httprouter.Router).Handle(fr.Method, fr.Url, fr.HandleOrFileSystem.(httprouter.Handle))
			}
		}
	} else {
		if routers.FrontRouters == nil {
			//假如前台路由为空，则设置
			routers.FrontRouters = mux.NewRouter()
		}
		for _, fr := range FrontRouters {
			if fr.Method == "ServeFiles" {
				routers.FrontRouters.(*mux.Router).PathPrefix(fr.Url).Handler(http.StripPrefix(fr.Url, http.FileServer(fr.HandleOrFileSystem.(http.FileSystem))))
			} else {
				routers.FrontRouters.(*mux.Router).HandleFunc(fr.Url, util.ToMux(fr.HandleOrFileSystem.(func(http.ResponseWriter, *http.Request, httprouter.Params))))
			}
		}
	}
}

// SetEvent 设置 监听事件 ，支持追加
func (g *Gef) SetEvent(EventAdd map[string][]Event.Listener) {
	Event.BindEvents(EventAdd)
}

// SetAdminStatic 设置静态文件，支持追加
func (g *Gef) SetAdminStatic(f embed.FS) {
	static.FileSystems = append(static.FileSystems, http.FS(f))
}

// SetAdminTemplate 设置后台静态文件，支持追加
func (g *Gef) SetAdminTemplate(f embed.FS) {
	adminTemplates.AdminTemplatesAdd = append(adminTemplates.AdminTemplatesAdd, f)
}

// SetServer 设置服务器，支持追加
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
			Router:     routers.AdminRouters(),
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
