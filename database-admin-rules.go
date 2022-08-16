/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: database-admin-rules
 * @Version: 1.0.0
 * @Date: 2022/8/16 22:19
 */

package gef

var adminRules = []map[string]interface{}{
	{
		"name":      "后台管理",
		"type":      1,
		"is_compel": 0,
		"icon":      "icon-slideshow-3-line",
		"route":     "#home",
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
				"route":     "#admin",
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
				"route":     "#log",
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
				"route":     "#attachment",
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
		"route":     "#app",
		"children":  []map[string]interface{}{},
	},
	{
		"name":      "设置",
		"type":      1,
		"is_compel": 0,
		"icon":      "icon-settings-6-line",
		"route":     "#config",
		"children": []map[string]interface{}{
			{
				"name":      "通用设置",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-settings-4-line",
				"route":     "#common",
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
		"route":     "#dev",
		"children": []map[string]interface{}{
			{
				"name":      "高级设置",
				"type":      1,
				"is_compel": 0,
				"icon":      "ri-tools-fill",
				"route":     "#common",
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
				"route":     "#em",
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
