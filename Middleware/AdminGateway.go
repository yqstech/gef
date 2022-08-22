/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 后台应用启动器，封装应用的前置校验，并启动应用
 * @File: AdminGateway
 * @Version: 1.0.0
 * @Date: 2021/10/14 5:33 下午
 */

package Middleware

import (
	"github.com/yqstech/gef/Handles/adminHandle"
	"github.com/yqstech/gef/Handles/commHandle"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/config"
	"github.com/yqstech/gef/util"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
)

// AdminGateway 后台网关,校验token、校验权限、easyModel自动注册   ====> Run() 启动nodePage或空白页
type AdminGateway struct {
	NodePages map[string]builder.NodePager
}

// Gateway 网关入口
func (that *AdminGateway) Gateway(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//功能封装
	that.checkToken(
		that.checkAuth(
			that.log(
				that.EasyModelAutoRegister(
					//启动注册的nodePage或空白页
					that.Run,
				),
			),
		),
	)(w, r, ps)
}

// EasyModelAutoRegister 动态注册页面；判断页面名称是em_开头，且页面集未定义，则自动注册到页面集中
func (that *AdminGateway) EasyModelAutoRegister(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		pageName := ps.ByName("pageName")
		//校验是否设置了对应的页面
		if _, ok := that.NodePages[pageName]; !ok {
			//em_开头的未定义页面，都设置到EasyModel页面
			if pageName[0:3] == "em_" {
				that.NodePages[pageName] = &adminHandle.EasyModelHandle{ModelKey: pageName[3:]}
			}
		}
		next(w, r, ps)
	}
}

//校验登录身份token
func (that *AdminGateway) checkToken(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//defer func() {
		//	if r := recover(); r != nil {
		//		fmt.Fprint(w, "程序异常：", r)
		//	}
		//}()

		//获取当前链接
		url := ps.ByName("pageName") + "/" + ps.ByName("actionName")
		//token是否免检
		checkTokenExclude := []interface{}{"account/login", "account/verifyCode"}
		if util.IsInArray(url, checkTokenExclude) {
			next(w, r, ps)
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
			that.toLogin(w, r)
			return
		}
		//根据token获取用户信息
		userinfo := Models.Admin{}.GetAccountInfoByToken(token)
		if userinfo == nil {
			that.toLogin(w, r)
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

		next(w, r, ps)
	}
}

//
//  checkAuth
//  @Description: 校验权限
//  @receiver admin
//  @param next
//  @return BootHandle
//
func (that *AdminGateway) checkAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		accountId := ps.ByName("account_id")
		mainAccountId := ps.ByName("main_account_id")
		//未登录的无需检查
		if mainAccountId == "" && accountId == "" {
			next(w, r, ps)
			return
		}
		//校验权限
		if mainAccountId == accountId {
			//主账户登录不校验权限
			next(w, r, ps)
			return
		}
		//子账户登录查询权限
		url := "/" + ps.ByName("pageName") + "/" + ps.ByName("actionName")
		if (Models.Admin{}).CheckAuth(url, util.String2Int(accountId)) {
			next(w, r, ps)
		} else {
			builder.NodePage{}.ErrResult(w, r, 120, "您无权进行此操作！", "")
		}
	}
}

//记录操作日志
func (that *AdminGateway) log(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.Method == "POST" {
			//操作地址
			url := "/" + ps.ByName("pageName") + "/" + ps.ByName("actionName")
			//查询信息看看是否需要记录日志
			conn := db.New().Table("tb_admin_rules")
			ruleInfo, err := conn.Where("is_delete", 0).Where("status", 1).Where("open_log", 1).Where("route", url).First()
			if err != nil {
				logger.Error(err.Error())
				builder.NodePage{}.ErrResult(w, r, 500, "系统出错了！", "")
				return
			}
			if ruleInfo != nil {
				//整理日志信息
				logInfo := map[string]interface{}{
					"rule":         url,
					"rule_name":    ruleInfo["name"],
					"url":          r.RequestURI,
					"account_id":   ps.ByName("account_id"),
					"account_name": ps.ByName("account_name"),
					"account":      ps.ByName("account"),
					"data":         "",
					"create_time":  util.TimeNow(),
				}
				id := util.PostValue(r, "id")
				if id != "" {
					logInfo["notice"] = "id=" + id
				}
				insertId, err := db.New().Table("tb_admin_log").InsertGetId(logInfo)
				if err != nil {
					logger.Error(err.Error())
					builder.NodePage{}.ErrResult(w, r, 500, "系统出错了！", "")
					return
				}
				logId := httprouter.Param{Key: "_log_id", Value: util.Int642String(insertId)}
				ps = append(ps, logId)
			}
		}
		next(w, r, ps)
	}
}

// 跳转到登录页面或者返回需要登录信息
func (that *AdminGateway) toLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		commHandle.Base{}.ApiResult(w, 100, "no token", nil)
	} else {
		w.Header().Set("Location", config.AdminPath+"/account/login")
		w.WriteHeader(302)
	}
}

// BindPages 绑定页面列表
//func (that *AdminGateway) BindPages(s map[string]EasyApp.AppPage) {
//	for k, v := range s {
//		that.NodePages[k] = v
//	}
//}

func (that *AdminGateway) Run(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//新建页面构建器
	pageBuilder := builder.PageBuilder{
		RequestID:          ps.ByName("requestID"),
		HttpResponseWriter: w,
		HttpRequest:        r,
		HttpParams:         ps,
	}
	//数据初始化
	pageBuilder.DataReset()

	//在页面集中查找页面
	if nodePage, ok := that.NodePages[ps.ByName("pageName")]; ok {
		//拷贝一个节点页面对象
		nodePageCopy := builder.NodePageCopy(nodePage)
		//页面运行
		nodePageCopy.Run(&pageBuilder, nodePageCopy)
	} else {
		//检索不到，创建一个空白页
		nodePage = &builder.EmptyPage{}
		//运行空白页
		nodePage.Run(&pageBuilder, nodePage)
	}
}
