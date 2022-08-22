/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 后台账户登录、退出、修改密码、修改信息
 * @File: Account
 * @Version: 1.0.0
 * @Date: 2021/10/14 10:42 下午
 */

package adminHandle

import (
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/pool"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/config"
	"github.com/yqstech/gef/util"
	"github.com/yqstech/gef/util/captcha"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
)

type Account struct {
	Base
}

// NodeInit 初始化
func (ac *Account) NodeInit(pageBuilder *builder.PageBuilder) {
	//清除默认handle
	ac.NodePageActions = map[string]httprouter.Handle{
		"resetpwd":   ac.ResetPwd,
		"login":      ac.Login,
		"logout":     ac.LogOut,
		"userinfo":   ac.UserInfo,
		"verifyCode": ac.verifyCode,
	}
}

// Login 登录页面
func (ac Account) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if r.Method == "POST" {
		//获取账户和密码
		account := util.PostValue(r, "account")
		pwd := util.PostValue(r, "pwd")
		if account == "" {
			ac.ApiResult(w, 101, "账号不得为空！", "")
			return
		}
		if pwd == "" {
			ac.ApiResult(w, 101, "密码不能为空！", "")
			return
		}
		//校验是否需要验证码
		failTimes, ok := pool.Gocache.Get("admin_login_fail_times")
		if ok {
			//假如失败次数达到三次以上
			if failTimes.(int) >= 3 {
				//校验码必填
				code := util.PostValue(r, "code")
				captchaId := util.PostValue(r, "captchaId")
				if code == "" {
					ac.ApiResult(w, 101, "请填入验证码！", map[string]interface{}{
						"verify": "reload",
					})
					return
				}
				if captchaId == "" {
					ac.ApiResult(w, 101, "验证码ID获取失败！", map[string]interface{}{
						"verify": "reload",
					})
					return
				}
				//获取真实验证码
				if !captcha.Verify(captchaId, code) {
					ac.ApiResult(w, 102, "验证码校验失败！", map[string]interface{}{
						"verify": "reload",
					})
					return
				}
			}
		}
		//按账户和密码查询账户
		where := map[string]interface{}{"is_delete": 0}
		where["account"] = account
		where["password"] = util.GetPassword(pwd)
		where["status"] = 1
		where["is_delete"] = 0
		data, err := ac.DbModel("tb_admin").Where(where).First()
		if err != nil {
			logger.Error(err.Error())
			ac.ApiResult(w, 500, "系统异常！", "")
			return
		}
		//查询为空
		if data == nil {
			failTimes, ok := pool.Gocache.Get("admin_login_fail_times")
			if ok {
				failTimes = failTimes.(int) + 1
			} else {
				failTimes = 1
			}
			pool.Gocache.Set("admin_login_fail_times", failTimes, time.Second*3600)
			if failTimes.(int) >= 3 {
				ac.ApiResult(w, 102, "账号或密码不正确！", map[string]interface{}{
					"verify": "verify",
				})
			} else {
				ac.ApiResult(w, 102, "账号或密码不正确！", nil)
			}
			return
		}
		//查询成功创建token
		token := util.MD5(util.TimeNow() + "xxxx")
		//记录token
		_, err = ac.DbModel("tb_admin_token").Insert(map[string]interface{}{
			"account_id":  data["id"],
			"token":       token,
			"create_time": util.TimeNow(),
		})
		if err != nil {
			logger.Error(err.Error())
			ac.ApiResult(w, 500, "系统异常！", "")
			return
		}
		ck := &http.Cookie{
			Name:   "admin_token",
			Value:  token,
			Path:   config.AdminPath + "/",
			MaxAge: 86600 * 7,
		}
		w.Header().Set("set-cookie", ck.String())
		ac.ApiResult(w, 200, "success", map[string]string{
			"token": token,
			"url":   config.AdminPath + "/index/index",
		})
		return
	}
	tpl := builder.Displayer{
		TplName: "login.html",
	}
	//默认不需要校验码
	verify := ""
	failTimes, ok := pool.Gocache.Get("admin_login_fail_times")
	if ok {
		if failTimes.(int) >= 3 {
			verify = "verify"
		}
	}
	tpl.SetDate("verify", verify)
	tpl.SetDate("title", "后台登录")
	tpl.SetDate("submit_url", config.AdminPath+"/account/login")
	pageBuilder := builder.PageBuilder{}
	pageBuilder.DataReset()
	ac.ActShow(w, tpl, &pageBuilder)
}

// ResetPwd 修改密码页面
func (ac Account) ResetPwd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if r.Method == "POST" {
		//当前账号ID
		accountId := ps.ByName("account_id")
		//获取数据
		password := util.PostValue(r, "password")
		newpassword := util.PostValue(r, "newpassword")
		if password == "" {
			ac.ApiResult(w, 101, "密码不得为空！", "")
			return
		}
		if newpassword == "" {
			ac.ApiResult(w, 101, "请输入新密码！", "")
			return
		}
		//查询当前账号信息
		accountInfo, err := ac.DbModel("tb_admin").Where("id", accountId).First()
		if err != nil {
			logger.Error(err.Error())
			ac.ApiResult(w, 500, "系统异常", "")
			return
		}
		//判断密码是否一致
		oldPwd := util.PostValue(r, "password")
		oldPwd = util.GetPassword(oldPwd)
		if accountInfo["password"].(string) != oldPwd {
			ac.ApiResult(w, 101, "旧密码不正确！", "")
			return
		}
		//修改密码
		_, err = ac.DbModel("tb_admin").Where("id", accountId).Update(map[string]string{
			"password": util.GetPassword(newpassword),
		})
		if err != nil {
			logger.Error(err.Error())
			ac.ApiResult(w, 500, "系统异常", "")
			return
		}
		ac.ApiResult(w, 200, "success", "success")
		return
	}
	tpl := builder.Displayer{
		TplName: "resetpwd.html",
	}
	tpl.SetDate("title", "修改密码")
	tpl.SetDate("postUrl", config.AdminPath+"/account/resetpwd")
	tpl.SetDate("successUrl", config.AdminPath+"/account/logout")
	pageBuilder := builder.PageBuilder{}
	pageBuilder.DataReset()
	ac.ActShow(w, tpl, &pageBuilder)
}

func (ac Account) UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//当前账号ID
	accountId := ps.ByName("account_id")

	//查询当前账号信息
	accountInfo, err := ac.DbModel("tb_admin").Where("id", accountId).First()
	if err != nil {
		logger.Error(err.Error())
		ac.ApiResult(w, 500, "系统异常", "")
		return
	}
	if accountInfo == nil {
		ac.ErrResult(w, r, 404, "账户信息不存在！", "")
		return
	}
	if r.Method == "POST" {
		//获取数据
		name := util.PostValue(r, "name")
		account := util.PostValue(r, "account")
		if name == "" {
			ac.ApiResult(w, 101, "名称不能为空", "")
			return
		}
		if account == "" {
			ac.ApiResult(w, 101, "账户不能为空", "")
			return
		}
		//判断是否已经使用
		otherAccount, err := ac.DbModel("tb_admin").
			Where("id", "!=", accountId).
			Where("account", account).First()
		if err != nil {
			logger.Error(err.Error())
			ac.ApiResult(w, 500, "系统异常", "")
			return
		}
		if otherAccount != nil {
			ac.ApiResult(w, 102, "账户名不可用！", "")
			return
		}
		//修改账户名
		_, err = ac.DbModel("tb_admin").Where("id", accountId).Update(map[string]string{
			"account": account,
			"name":    name,
		})
		if err != nil {
			logger.Error(err.Error())
			ac.ApiResult(w, 500, "系统异常", "")
			return
		}
		ac.ApiResult(w, 200, "success", "success")
		return
	}
	tpl := builder.Displayer{
		TplName: "account_info.html",
	}
	tpl.SetDate("account_name", accountInfo["name"])
	tpl.SetDate("account", accountInfo["account"])
	tpl.SetDate("title", "修改资料")
	tpl.SetDate("postUrl", config.AdminPath+"/account/userinfo")
	tpl.SetDate("successUrl", config.AdminPath+"/index/index")
	pageBuilder := builder.PageBuilder{}
	pageBuilder.DataReset()
	ac.ActShow(w, tpl, &pageBuilder)
}

func (ac Account) LogOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//清空账户Cookie
	accountId := util.String2Int(ps.ByName("account_id"))
	db.New().Table("tb_admin_token").Where("account_id", accountId).Delete()
	ck := &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   config.AdminPath + "/",
		MaxAge: -1,
	}
	w.Header().Set("set-cookie", ck.String())
	tpl := builder.Displayer{
		TplName: "logout.html",
	}
	tpl.SetDate("successUrl", config.AdminPath+"/account/login")

	pageBuilder := builder.PageBuilder{}
	pageBuilder.DataReset()
	ac.ActShow(w, tpl, &pageBuilder)
}

func (ac Account) verifyCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//数字验证码配置
	captchaId, base64, err := captcha.GetCaptchaBase64("auto", 40, 130, 6, captcha.ColorWight)
	if err != nil {
		logger.Error(err.Error())
		ac.ApiResult(w, 500, "获取验证码失败！", nil)
		return
	}
	ac.ApiResult(w, 200, "success", map[string]interface{}{
		"captchaId":     captchaId,
		"captchaBase64": base64,
	})
}
