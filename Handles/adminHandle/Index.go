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
	"github.com/yqstech/gef/EasyApp"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/util"
	"github.com/yqstech/gef/config"
	"net/http"
	
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
)

type Index struct {
	Base
}

// PageInit 设置Handel路由
func (index Index) PageInit(pageData *EasyApp.PageData) {
	pageData.ActionAdd("get_menus", index.GetMenus)
	pageData.ActionAdd("main", index.Main)
}

// Index 后台主页框架
func (index Index) Index(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tpl := EasyApp.Template{
		TplName: "index.html",
	}
	tpl.SetDate("title", Models.AppConfigs{}.Value("app_name")+" - 后台管理中心")
	tpl.SetDate("login_path", config.AdminPath+"/account/login")
	tpl.SetDate("menu_path", config.AdminPath+"/index/get_menus")
	tpl.SetDate("logo", "/static/images/logo.png")
	index.ActShow(w, tpl, pageData)
}

// Main 后台欢迎页
func (index Index) Main(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tpl := EasyApp.Template{
		TplName: "main.html",
	}
	tpl.SetDate("title", "首页")
	index.ActShow(w, tpl, pageData)
}

// GetMenus 获取菜单接口，支持顶部菜单，左侧菜单，右侧菜单
func (index Index) GetMenus(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//主账户ID
	main_account_id := ps.ByName("main_account_id")
	//当前账户ID
	account_id := ps.ByName("account_id")
	account_name := ps.ByName("account_name")
	account := ps.ByName("account")
	//当前账户所属分组角色
	group_id := ps.ByName("group_id")
	
	//定义权限表
	var rules interface{}
	var err error
	
	userMenus := []map[string]interface{}{
		{
			"name":   "退出",
			"icon":   "layui-icon-release",
			"url":    config.AdminPath + "/account/logout",
			"target": "",
		},
		{
			"name":   "修改密码",
			"icon":   "layui-icon-password",
			"url":    config.AdminPath + "/account/resetpwd",
			"target": "main_area",
		},
		{
			"name":   account_name,
			"icon":   "layui-icon-username",
			"url":    config.AdminPath + "/account/userinfo",
			"target": "main_area",
		},
	}
	
	if account_id == main_account_id {
		//当为主账户时，获取所有菜单
		rules, err = db.New().Table("tb_admin_rules").
			Where("is_delete", "=", 0).
			Where("type", "=", 1).
			Where("status", "=", 1).
			OrderBy("index_num asc,id asc").
			Get()
		if err != nil {
			logger.Error(err.Error())
			index.ApiResult(w, 110, "error", "获取菜单失败！")
			return
		}
	} else {
		//否则，获取账户角色的菜单
		if group_id != "" && group_id != "0" {
			logger.Info(group_id, "子账户角色ID")
			groupInfo, err := db.New().Table("tb_admin_group").
				Where("id", int64(util.String2Int(group_id))).
				First()
			if err != nil {
				logger.Error(err.Error())
				index.ApiResult(w, 110, "error", "获取角色失败！")
				return
			}
			if groupInfo == nil {
				index.ApiResult(w, 110, "error", "角色无效！")
				return
			}
			rule_ids := []int64{}
			util.JsonDecode(groupInfo["rules"].(string), &rule_ids)
			ruleIDs := []int64{}
			for _, v := range rule_ids {
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
				index.ApiResult(w, 110, "error", "获取菜单失败！")
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
				index.ApiResult(w, 110, "error", "获取菜单失败！")
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
	topMenus := []map[string]interface{}{}
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
	ruleIndex := []int64{}
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
				"name":   rule["name"],
				"icon":   rule["icon"],
				"url":    rule["route"],
				"hidden": false,
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
					"name":   rule["name"].(string),
					"url":    rule["route"].(string),
					"active": "",
				})
			}
		}
	}
	//转换成菜单结构
	menuArr := []interface{}{}
	for _, Index := range ruleIndex {
		menuArr = append(menuArr, ruleMap[Index])
	}
	
	index.ApiResult(w, 200, "success", map[string]interface{}{
		"menus":        menuArr,
		"topMenus":     topMenus,
		"userMenus":    userMenus,
		"account":      account,
		"account_name": account_name,
	})
}
