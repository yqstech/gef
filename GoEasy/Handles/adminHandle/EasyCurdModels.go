/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyCurdModels
 * @Version: 1.0.0
 * @Date: 2022/5/3 21:22
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/config"
)

type EasyCurdModels struct {
	Base
}

// NodeBegin 开始
func (that EasyCurdModels) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("EasyCurdModel 接口模型")
	pageData.SetPageName("接口模型")
	pageData.SetTbName("tb_easy_curd_models")
	return nil, 0
}

// NodeList 初始化列表
func (that EasyCurdModels) NodeList(pageData *EasyApp.PageData) (error, int) {
	//!重置顶部按钮
	//pageData.SetListTopBtns("add", "select_data")

	//新增右侧字段管理按钮
	pageData.SetButton("fields", EasyApp.Button{
		ButtonName: "字段管理",
		Action:     "/easy_curd_models_fields/index",
		ActionType: 2,
		LayerTitle: "模型字段管理",
		ActionUrl:  config.AdminPath + "/easy_curd_models_fields/index",
		Class:      "def",
		Icon:       "ri-table-line",
		Display:    "!item.btn_field || item.btn_field!='hide'",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//!重置右侧按钮
	pageData.SetListRightBtns("edit", "disable", "enable", "fields", "delete")
	pageData.ListColumnAdd("model_key", "模型Key", "text", nil)
	pageData.ListColumnAdd("model_name", "模型名称", "text", nil)
	//pageData.ListColumnAdd("table_name", "数据表名", "text", nil)
	pageData.ListColumnAdd("allow_select", "查询", "switch::text=允许|禁止", nil)
	pageData.ListColumnAdd("allow_create", "新增", "switch::text=允许|禁止", nil)
	pageData.ListColumnAdd("allow_update", "更新", "switch::text=允许|禁止", nil)
	pageData.ListColumnAdd("allow_delete", "删除", "switch::text=允许|禁止", nil)
	//pageData.ListColumnAdd("soft_delete_disable", "删除方式", "switch::text=硬删除|软删除", nil)
	//pageData.ListColumnAdd("check_login", "校验登录", "switch::text=是|否", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	return nil, 0
}

// NodeForm 初始化表单
func (that EasyCurdModels) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("model_key", "text", "模型Key", "模型唯一识别码", "", true, nil, "", nil)
	pageData.FormFieldsAdd("model_name", "text", "模型名称", "模型自定义名称", "", true, nil, "", nil)
	pageData.FormFieldsAdd("table_name", "text", "关联数据表", "关联的数据表完整名称", "tb_", true, nil, "", nil)
	pageData.FormFieldsAdd("select_order", "text", "查询排序", "查询排序方式", "id desc", true, nil, "", nil)
	pageData.FormFieldsAdd("allow_select", "radio", "允许查询", "", "1", true, Models.DefaultIsOrNot, "", nil)
	pageData.FormFieldsAdd("allow_create", "radio", "允许新增", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageData.FormFieldsAdd("allow_update", "radio", "允许更新", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageData.FormFieldsAdd("allow_delete", "radio", "允许删除", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageData.FormFieldsAdd("soft_delete_disable", "radio", "禁用软删除", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageData.FormFieldsAdd("check_login", "radio", "校验登录", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageData.FormFieldsAdd("uk_name", "text", "用户字段名", "代表用户id的字段名称", "user_id", true, nil, "", nil)
	pageData.FormFieldsAdd("pk_name", "text", "主键字段名", "代表用主键的字段名称", "id", true, nil, "", nil)
	return nil, 0
}
