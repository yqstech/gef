/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: adminRules
 * @Version: 1.0.0
 * @Date: 2022/8/16 16:33
 */

package database

import (
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
)

// AutoAdminRules 自动维护后台字段
func AutoAdminRules(rules []map[string]interface{}) {
	setAdminRules(0, rules)
}

func setAdminRules(pid int64, rules []map[string]interface{}) {
	for index, rule := range rules {
		//route参数必填
		if rule["route"].(string) == "" {
			continue
		}
		ruleId := int64(0)
		//查找已经插入的数据
		ruleInfo, err := db.New().Table("tb_admin_rules").
			Where("is_delete", 0).
			Where("pid", pid). //同一个pid下不能有重复的route
			Where("route", rule["route"].(string)).
			First()
		if err != nil {
			panic("权限更新失败！" + err.Error())
			return
		}
		//否则新增这个数据
		if ruleInfo == nil {
			insertId, err := db.New().Table("tb_admin_rules").InsertGetId(map[string]interface{}{
				"pid":         pid,
				"name":        rule["name"],
				"type":        rule["type"],
				"is_compel":   rule["is_compel"],
				"icon":        rule["icon"],
				"route":       rule["route"],
				"index_num":   index + 1,
				"create_time": util.TimeNow(),
				"update_time": util.TimeNow(),
			})
			if err != nil {
				panic("权限更新失败！" + err.Error())
				return
			}
			ruleId = insertId
		} else {
			ruleId = ruleInfo["id"].(int64)
			//值不全，不更新
			if _, ok := rule["name"]; !ok {
				goto EndUpdate
			}
			if _, ok := rule["type"]; !ok {
				goto EndUpdate
			}
			if _, ok := rule["is_compel"]; !ok {
				goto EndUpdate
			}
			if _, ok := rule["icon"]; !ok {
				goto EndUpdate
			}
			//更新一下
			db.New().Table("tb_admin_rules").Where("id", ruleId).Update(map[string]interface{}{
				"name":        rule["name"],
				"type":        rule["type"],
				"is_compel":   rule["is_compel"],
				"icon":        rule["icon"],
				"update_time": util.TimeNow(),
			})
			EndUpdate:
		}
		if children, ok := rule["children"]; ok && len(children.([]map[string]interface{})) > 0 {
			setAdminRules(ruleId, children.([]map[string]interface{}))
		}
	}
}

var adminRules = []map[string]interface{}{
	{
		"name":      "后台管理",
		"type":      1,
		"is_compel": 0,
		"icon":      "icon-slideshow-3-line",
		"route":     "home",
		"children": []map[string]interface{}{
			{
				"name":      "仪表盘",
				"type":      1,
				"is_compel": 1,
				"icon":      "icon-dashboard-3-line",
				"route":     "/index/index",
				"children": []map[string]interface{}{
					{
						"name":      "获取菜单",
						"type":      2,
						"is_compel": 1,
						"icon":      "",
						"route":     "/index/get_menus",
						"children":  []map[string]interface{}{},
					}, {
						"name":      "欢迎页",
						"type":      2,
						"is_compel": 1,
						"icon":      "",
						"route":     "/index/main",
						"children":  []map[string]interface{}{},
					},
				},
			},
			{
				"name":      "后台账户",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-shield-user-line",
				"route":     "admin",
				"children": []map[string]interface{}{
					{
						"name":      "角色管理",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/admin_group/index",
						"children":  []map[string]interface{}{},
					}, {
						"name":      "账户管理",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/admin/index",
						"children":  []map[string]interface{}{},
					},
				},
			}, {
				"name":      "操作日志",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-time-line",
				"route":     "log",
				"children": []map[string]interface{}{
					{
						"name":      "操作日志",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/admin_log/index",
						"children":  []map[string]interface{}{},
					},
				},
			}, {
				"name":      "附件管理",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-gallery-upload-line",
				"route":     "attachment",
				"children": []map[string]interface{}{
					{
						"name":      "图片管理",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/attachment_image/index",
						"children":  []map[string]interface{}{},
					}, {
						"name":      "文件管理",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/attachment_file/index",
						"children":  []map[string]interface{}{},
					},
				},
			},
			{
				"name":      "修改资料",
				"type":      2,
				"is_compel": 1,
				"icon":      "",
				"route":     "/account/userinfo",
				"children":  []map[string]interface{}{},
			}, {
				"name":      "修改密码",
				"type":      2,
				"is_compel": 1,
				"icon":      "",
				"route":     "/account/resetpwd",
				"children":  []map[string]interface{}{},
			}, {
				"name":      "退出",
				"type":      2,
				"is_compel": 1,
				"icon":      "",
				"route":     "/account/logout",
				"children":  []map[string]interface{}{},
			},
		},
	},
	{
		"name":      "应用管理",
		"type":      1,
		"is_compel": 0,
		"icon":      "ri-apps-line",
		"route":     "app",
		"children":  []map[string]interface{}{},
	},
	{
		"name":      "设置",
		"type":      1,
		"is_compel": 0,
		"icon":      "icon-settings-6-line",
		"route":     "config",
		"children": []map[string]interface{}{
			{
				"name":      "通用设置",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-settings-4-line",
				"route":     "common",
				"children": []map[string]interface{}{
					{
						"name":      "系统设置",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/app_configs_g1/index",
						"children":  []map[string]interface{}{},
					},
				},
			},
		},
	},
	{
		"name":      "开发",
		"type":      1,
		"is_compel": 0,
		"icon":      "icon-code-box-line",
		"route":     "dev",
		"children": []map[string]interface{}{
			{
				"name":      "高级设置",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-tools-fill",
				"route":     "common",
				"children": []map[string]interface{}{
					{
						"name":      "后台权限管理",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/admin_rules/index",
						"children": []map[string]interface{}{
							{
								"name":      "新增权限",
								"type":      2,
								"is_compel": 0,
								"icon":      "",
								"route":     "/admin_rules/add",
								"children":  []map[string]interface{}{},
							}, {
								"name":      "修改权限",
								"type":      2,
								"is_compel": 0,
								"icon":      "",
								"route":     "/admin_rules/edit",
								"children":  []map[string]interface{}{},
							}, {
								"name":      "删除权限",
								"type":      2,
								"is_compel": 0,
								"icon":      "",
								"route":     "/admin_rules/delete",
								"children":  []map[string]interface{}{},
							}, {
								"name":      "禁用权限",
								"type":      2,
								"is_compel": 0,
								"icon":      "",
								"route":     "/admin_rules/status",
								"children":  []map[string]interface{}{},
							},
						},
					},
					{
						"name":      "设置项管理",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/configs/index",
						"children":  []map[string]interface{}{},
					},
					{
						"name":      "设置组管理",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/configs_group/index",
						"children":  []map[string]interface{}{},
					}, {
						"name":      "选项集",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/option_models/index",
						"children":  []map[string]interface{}{},
					},
				},
			},
			{
				"name":      "简易模型",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-stack-fill",
				"route":     "em",
				"children": []map[string]interface{}{
					{
						"name":      "后台模型",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/easy_models/index",
						"children":  []map[string]interface{}{},
					}, {
						"name":      "接口模型",
						"type":      1,
						"is_compel": 0,
						"icon":      "",
						"route":     "/easy_curd_models/index",
						"children":  []map[string]interface{}{},
					},
				},
			},
		},
	},
}
