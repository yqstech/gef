/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AdminRouters
 * @Version: 1.0.0
 * @Date: 2022/3/8 12:44 下午
 */

package Routers

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Handles/commHandle"
	"github.com/gef/GoEasy/Middleware"
	"github.com/gef/GoEasy/Registry"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gef/config"
	"github.com/gef/static"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// AdminRouters 后台路由
func AdminRouters() *httprouter.Router {
	//引入公共路由
	router := httprouter.New()
	//静态资源 //不封装静态资源则使用 router.ServeFiles("/static/*filepath", http.Dir("static/"))
	//静态资源，累加静态资源库
	fs := util.FileSystems{}
	for _, f := range static.FileSystems {
		fs = append(fs, f)
	}
	fs = append(fs, http.FS(static.Files))
	router.ServeFiles("/static/*filepath", fs)
	//存储上传
	router.ServeFiles("/uploads/*filepath", http.Dir("data/uploads/"))
	
	//工具包，打印post提交信息
	router.POST("/utils/requestPrint", commHandle.Utils{}.RequestPrint)
	
	//当前服务控制
	router.GET("/server/manage/pid", commHandle.Server{}.Pid)         //PID
	router.GET("/server/manage/restart", commHandle.Server{}.Restart) //重启
	
	//后台模块启动器
	AdminBoot := Middleware.AdminBoot{AppPages: map[string]EasyApp.AppPage{}}
	//绑定Page列表
	AdminBoot.BindPages(Registry.AdminPages)
	//后台全部请求，引导到后台启动器的入口
	router.POST(config.AdminPath+"/:pageName/:actionName", AdminBoot.Gateway)
	router.GET(config.AdminPath+"/:pageName/:actionName", AdminBoot.Gateway)
	return router
}
