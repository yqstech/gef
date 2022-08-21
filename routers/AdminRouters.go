/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AdminRouters
 * @Version: 1.0.0
 * @Date: 2022/3/8 12:44 下午
 */

package routers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/yqstech/gef/Handles/commHandle"
	"github.com/yqstech/gef/Middleware"
	"github.com/yqstech/gef/Registry"
	"github.com/yqstech/gef/Utils/util"
	"github.com/yqstech/gef/config"
	"github.com/yqstech/gef/static"
	"net/http"
)

// AdminRouters 后台路由
func AdminRouters() *httprouter.Router {
	router := httprouter.New()
	
	//!静态资源 (自动包含项目追加)
	//不封装静态资源则使用 router.ServeFiles("/static/*filepath", http.Dir("static/"))
	//静态资源，累加静态资源库
	fs := util.FileSystems{}
	for _, f := range static.FileSystems {
		fs = append(fs, f)
	}
	fs = append(fs, http.FS(static.Files))
	router.ServeFiles("/static/*filepath", fs)
	
	//!静态资源2 —— 上传的文件（锁定固定文件夹 data/uploads ）
	router.ServeFiles("/uploads/*filepath", http.Dir("data/uploads/"))
	
	//!普通handle方法
	router.POST("/utils/requestPrint", commHandle.Utils{}.RequestPrint)
	router.GET("/server/manage/pid", commHandle.Server{}.Pid)         //PID
	router.GET("/server/manage/restart", commHandle.Server{}.Restart) //重启
	
	//!后台页面
	//创建后台网关,传入注册的
	AdminGateway := Middleware.AdminGateway{NodePages: Registry.AdminPages}
	router.POST(config.AdminPath+"/:pageName/:actionName", AdminGateway.Gateway)
	router.GET(config.AdminPath+"/:pageName/:actionName", AdminGateway.Gateway)
	
	return router
}
