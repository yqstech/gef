/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 首页与公共方法
 * @File: Index
 * @Version: 1.0.0
 * @Date: 2021/10/15 2:43 下午
 */

package adminHandle

import (
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/config"
	"github.com/yqstech/gef/util"
	"net/http"

	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
)

type Index struct {
	Base
}

// NodeInit NodeInit 设置Handel路由
func (that *Index) NodeInit(pageBuilder *builder.PageBuilder) {
	that.NodePageActions["get_menus"] = that.GetMenus
	that.NodePageActions["main"] = that.Main
}

// Index 后台主页框架
func (that *Index) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tpl := builder.Displayer{
		TplName: "index.html",
	}
	tpl.SetDate("title", Models.AppConfigs{}.Value("app_name")+" - 后台管理中心")
	tpl.SetDate("login_path", config.AdminPath+"/account/login")
	tpl.SetDate("menu_path", config.AdminPath+"/index/get_menus")
	tpl.SetDate("logo", "/static/images/logo.png")
	that.ActShow(w, tpl, that.PageBuilder)
}

// Main 后台欢迎页
func (that *Index) Main(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tpl := builder.Displayer{
		TplName: "main.html",
	}
	tpl.SetDate("app", Models.AppConfigs{}.Values(2))
	//当前账户信息
	tpl.SetDate("account", map[string]interface{}{
		"account_id":   ps.ByName("account_id"),
		"account_name": ps.ByName("account_name"),
		"account":      ps.ByName("account"),
	})
	that.ActShow(w, tpl, that.PageBuilder)
}

// GetMenus 获取菜单接口，支持顶部菜单，左侧菜单，右侧菜单
func (that Index) GetMenus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//主账户ID
	mainAccountId := ps.ByName("main_account_id")
	//当前账户ID
	accountId := ps.ByName("account_id")
	accountName := ps.ByName("account_name")
	account := ps.ByName("account")
	//当前账户所属分组角色
	groupId := ps.ByName("group_id")

	//定义权限表
	var rules interface{}
	var err error

	userMenus := []map[string]interface{}{
		{
			"name":   "退出",
			"icon":   "ri-logout-circle-r-line",
			"url":    config.AdminPath + "/account/logout",
			"target": "",
		},
		{
			"name":   "修改密码",
			"icon":   "ri-shield-user-line",
			"url":    config.AdminPath + "/account/resetpwd",
			"target": "main_area",
		},
		{
			"name":   accountName,
			"icon":   "ri-user-3-line",
			"url":    config.AdminPath + "/account/userinfo",
			"target": "main_area",
		},
	}

	if accountId == mainAccountId {
		//当为主账户时，获取所有菜单
		rules, err = db.New().Table("tb_admin_rules").
			Where("is_delete", "=", 0).
			Where("type", "=", 1).
			Where("status", "=", 1).
			OrderBy("index_num asc,id asc").
			Get()
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 110, "error", "获取菜单失败！")
			return
		}
	} else {
		//否则，获取账户角色的菜单
		if groupId != "" && groupId != "0" {
			logger.Info(groupId, "子账户角色ID")
			groupInfo, err := db.New().Table("tb_admin_group").
				Where("id", int64(util.String2Int(groupId))).
				First()
			if err != nil {
				logger.Error(err.Error())
				that.ApiResult(w, 110, "error", "获取角色失败！")
				return
			}
			if groupInfo == nil {
				that.ApiResult(w, 110, "error", "角色无效！")
				return
			}
			var ruleIds []int64
			util.JsonDecode(groupInfo["rules"].(string), &ruleIds)
			var ruleIDs []int64
			for _, v := range ruleIds {
				ruleIDs = append(ruleIDs, v)
			}
			//获取有效菜单
			conn := db.New().Table("tb_admin_rules")
			rules, err = conn.
				Where("is_delete", "=", 0).
				Where("type", "=", 1).
				Where("status", "=", 1).
				Where(func() {
					conn.Where("is_compel", 1).OrWhere("id", "in", ruleIDs)
				}).
				OrderBy("index_num asc,id asc").
				Get()
			if err != nil {
				logger.Error(err.Error())
				that.ApiResult(w, 110, "error", "获取菜单失败！")
				return
			}
			logger.Info(rules)

		} else {
			//未分配角色，获取默认菜单
			conn := db.New().Table("tb_admin_rules")
			rules, err = conn.
				Where("is_delete", "=", 0).
				Where("type", "=", 1).
				Where("status", "=", 1).
				Where("is_compel", 1).
				OrderBy("index_num asc,id asc").
				Get()
			if err != nil {
				logger.Error(err.Error())
				that.ApiResult(w, 110, "error", "获取菜单失败！")
				return
			}
		}
	}

	//!cookie 获取menuGroupID
	menuGroupID := int64(0)
	menuGroup, err := r.Cookie("menuGroupID")
	if err == nil {
		menuGroupID = int64(util.String2Int(menuGroup.Value))
	}

	//获取顶部一级菜单，并确定选中的菜单ID =》topMenuActiveID
	//topMenus
	var topMenus []map[string]interface{}
	topMenuActiveID := int64(0)
	for _, rule := range rules.([]gorose.Data) {
		if rule["pid"].(int64) == 0 {
			active := false
			if menuGroupID == rule["id"].(int64) {
				active = true
				topMenuActiveID = rule["id"].(int64)
			}
			//#开头的路径不对外显示
			if rule["route"].(string) != "" && rule["route"].(string)[0:1] == "#" {
				rule["route"] = ""
			}
			topMenus = append(topMenus, map[string]interface{}{
				"id":     rule["id"],
				"name":   rule["name"].(string),
				"icon":   rule["icon"],
				"url":    rule["route"],
				"active": active,
			})
		}
	}
	if topMenuActiveID == int64(0) {
		for k, menu := range topMenus {
			topMenus[k]["active"] = true
			topMenuActiveID = menu["id"].(int64)
			break
		}
	}

	//菜单结构
	ruleMap := map[int64]map[string]interface{}{}
	//记录一级菜单顺序，map是无序的
	var ruleIndex []int64
	for _, rule := range rules.([]gorose.Data) {
		if rule["pid"].(int64) == topMenuActiveID {
			ruleIndex = append(ruleIndex, rule["id"].(int64))

			//#开头的路径不对外显示
			if rule["route"].(string) != "" && rule["route"].(string)[0:1] == "#" {
				rule["route"] = ""
			} else if rule["route"].(string) != "" {
				rule["route"] = config.AdminPath + rule["route"].(string)
			}
			ruleMap[rule["id"].(int64)] = map[string]interface{}{
				"id":     util.Int642String(rule["id"].(int64)),
				"name":   rule["name"],
				"icon":   rule["icon"],
				"url":    rule["route"],
				"hidden": len(ruleIndex) > 4,
				"active": false,
				"list":   []map[string]string{},
			}
		}
	}
	for _, rule := range rules.([]gorose.Data) {
		if rule["pid"].(int64) > 0 {
			if _, ok := ruleMap[rule["pid"].(int64)]; ok {
				//#开头的路径不对外显示
				if rule["route"].(string) != "" && rule["route"].(string)[0:1] == "#" {
					rule["route"] = ""
				}
				if rule["route"].(string) != "" {
					rule["route"] = config.AdminPath + rule["route"].(string)
				}
				ruleMap[rule["pid"].(int64)]["list"] = append(ruleMap[rule["pid"].(int64)]["list"].([]map[string]string), map[string]string{
					"id":     util.Int642String(rule["id"].(int64)),
					"name":   rule["name"].(string),
					"url":    rule["route"].(string),
					"active": "",
				})
			}
		}
	}
	//转换成菜单结构
	var menuArr []interface{}
	for _, Index := range ruleIndex {
		menuArr = append(menuArr, ruleMap[Index])
	}

	that.ApiResult(w, 200, "success", map[string]interface{}{
		"menus":        menuArr,
		"topMenus":     topMenus,
		"userMenus":    userMenus,
		"account":      account,
		"account_name": accountName,
	})
}
