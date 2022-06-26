/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 后台应用启动器，封装应用的前置校验，并启动应用
 * @File: AdminBoot
 * @Version: 1.0.0
 * @Date: 2021/10/14 5:33 下午
 */

package Middleware

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Handles/adminHandle"
	"github.com/gef/GoEasy/Handles/commHandle"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gef/config"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
)

// AdminBoot
// 后台应用启动器,继承自EasyApp应用启动器
// 封装应用的前置校验，并启动应用
type AdminBoot struct {
	EasyApp.AppBoot
	AppPages map[string]EasyApp.AppPage
}

// BindPages 绑定页面列表
func (admin *AdminBoot) BindPages(s map[string]EasyApp.AppPage) {
	for k, v := range s {
		admin.AppPages[k] = v
	}
}

//
// Gateway
// 后台入口，一个普通的Handel
// 中间件封装BootHandle
func (admin *AdminBoot) Gateway(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//嵌套中间件，执行BootHandle方法
	admin.SafeHandle(
		admin.checkToken(
			admin.checkAuth(
				admin.log(
					admin.EasyModel(
						admin.Run,
					),
				),
			),
		),
	)(admin.AppPages, w, r, ps)
}

// EasyModel em_开头的未定义页面都转发到easyModel页面
func (admin AdminBoot) EasyModel(next EasyApp.BootHandle) EasyApp.BootHandle {
	return func(appPages map[string]EasyApp.AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		pageName := ps.ByName("pageName")
		//校验是否设置了对应的页面
		if _, ok := admin.AppPages[pageName]; !ok {
			//em_开头的未定义页面，都设置到EasyModel页面
			if pageName[0:3] == "em_" {
				admin.AppPages[pageName] = adminHandle.EasyModel{ModelKey: pageName[3:]}
			}
		}
		next(appPages, w, r, ps)
	}
}

//校验登录身份token
func (admin AdminBoot) checkToken(next EasyApp.BootHandle) EasyApp.BootHandle {
	return func(appPages map[string]EasyApp.AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//校验后台域名
		//adminDomain := models.GroupConfigs{}.ConfigValue("admin_domain")
		//if adminDomain != "" {
		//	//判断域名
		//
		//}

		//获取当前链接
		url := ps.ByName("pageName") + "/" + ps.ByName("actionName")
		//token是否免检
		checkTokenExclude := []interface{}{"account/login", "account/verifyCode"}
		if util.IsInArray(url, checkTokenExclude) {
			next(appPages, w, r, ps)
			return
		}
		//获取token
		token := util.GetValue(r, "admin_token")
		if token == "" {
			token = util.PostValue(r, "admin_token")
		}
		if token == "" {
			tk, err := r.Cookie("admin_token")
			if err == nil {
				token = tk.Value
			}
		}
		//token为空
		if token == "" {
			logger.Info("token为空")
			admin.toLogin(w, r)
			return
		}
		//根据token获取用户信息
		userinfo := Models.Admin{}.GetAccountInfoByToken(token)
		if userinfo == nil {
			logger.Info("token无效")
			admin.toLogin(w, r)
			return
		}
		//用户信息传递到处理程序
		//账户ID
		accountId := httprouter.Param{Key: "account_id", Value: util.Int642String(userinfo["account_id"].(int64))}
		//主账户ID 后台主账户为1
		mainAccountId := httprouter.Param{Key: "main_account_id", Value: "1"}
		//所属角色
		groupId := httprouter.Param{Key: "group_id", Value: util.Interface2String(userinfo["group_id"])}
		//名称
		accountName := httprouter.Param{Key: "account_name", Value: userinfo["name"].(string)}
		//账户
		account := httprouter.Param{Key: "account", Value: userinfo["account"].(string)}
		ps = append(ps, accountId)
		ps = append(ps, accountName)
		ps = append(ps, account)
		ps = append(ps, mainAccountId)
		ps = append(ps, groupId)

		//新增请求ID编码
		requestID := "admin_" + util.Int642String(userinfo["account_id"].(int64)) + "_" + util.Int642String(time.Now().UnixNano())
		ps = append(ps, httprouter.Param{Key: "requestID", Value: requestID})

		//上传资源分组（1后台组）
		uploadGroupID := httprouter.Param{Key: "uploadGroupID", Value: "1"}
		ps = append(ps, uploadGroupID)
		//上传资源用户
		uploadUserID := httprouter.Param{Key: "uploadUserID", Value: util.Int642String(userinfo["account_id"].(int64))}
		ps = append(ps, uploadUserID)

		next(appPages, w, r, ps)
	}
}

//
//  checkAuth
//  @Description: 校验权限
//  @receiver admin
//  @param next
//  @return BootHandle
//
func (admin AdminBoot) checkAuth(next EasyApp.BootHandle) EasyApp.BootHandle {
	return func(appPages map[string]EasyApp.AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		accountId := ps.ByName("account_id")
		mainAccountId := ps.ByName("main_account_id")
		//未登录的无需检查
		if mainAccountId == "" && accountId == "" {
			next(appPages, w, r, ps)
			return
		}
		//校验权限
		if mainAccountId == accountId {
			//主账户登录不校验权限
			next(appPages, w, r, ps)
			return
		}
		//子账户登录查询权限
		url := "/" + ps.ByName("pageName") + "/" + ps.ByName("actionName")
		if (Models.Admin{}).CheckAuth(url, util.String2Int(accountId)) {
			next(appPages, w, r, ps)
		} else {
			EasyApp.Page{}.ErrResult(w, r, 120, "您无权进行此操作！", "")
		}
	}
}

//记录操作日志
func (admin AdminBoot) log(next EasyApp.BootHandle) EasyApp.BootHandle {
	return func(appPages map[string]EasyApp.AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.Method == "POST" {
			//操作地址
			url := "/" + ps.ByName("pageName") + "/" + ps.ByName("actionName")
			//查询信息看看是否需要记录日志
			conn := db.New().Table("tb_admin_rules")
			ruleInfo, err := conn.Where("is_delete", 0).Where("status", 1).Where("open_log", 1).Where("route", url).First()
			if err != nil {
				logger.Error(err.Error())
				EasyApp.Page{}.ErrResult(w, r, 500, "系统出错了！", "")
				return
			}
			if ruleInfo != nil {

				//记录Request信息
				// rdump, err := httputil.DumpRequest(r, true)
				// if err != nil {
				// 	logger.Error(err.Error())
				// 	nodeFlow.EasyApp{}.ErrResult(w, r, 500, "系统出错了！", "")
				// 	return
				// }
				//整理日志信息
				logInfo := map[string]interface{}{
					"rule":         url,
					"rule_name":    ruleInfo["name"],
					"url":          r.RequestURI,
					"account_id":   ps.ByName("account_id"),
					"account_name": ps.ByName("account_name"),
					"account":      ps.ByName("account"),
					"data":         ps.ByName("_body"),
					"create_time":  util.TimeNow(),
				}
				id := util.PostValue(r, "id")
				if id != "" {
					logInfo["notice"] = "id=" + id
				}
				insertId, err := db.New().Table("tb_admin_log").InsertGetId(logInfo)
				if err != nil {
					logger.Error(err.Error())
					EasyApp.Page{}.ErrResult(w, r, 500, "系统出错了！", "")
					return
				}
				logId := httprouter.Param{Key: "_log_id", Value: util.Int642String(insertId)}
				ps = append(ps, logId)
			}
		}
		next(appPages, w, r, ps)
	}
}

// 跳转到登录页面或者返回需要登录信息
func (admin AdminBoot) toLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		commHandle.Base{}.ApiResult(w, 100, "no token", nil)
	} else {
		w.Header().Set("Location", config.AdminPath+"/account/login")
		w.WriteHeader(302)
	}
}
