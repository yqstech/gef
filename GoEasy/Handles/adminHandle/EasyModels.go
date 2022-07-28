/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyModels
 * @Version: 1.0.0
 * @Date: 2022/5/1 22:16
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/config"
	"github.com/wonderivan/logger"
)

type EasyModels struct {
	Base
}

// NodeBegin 开始
func (that EasyModels) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("EasyModel 后台模型")
	pageData.SetPageName("模型")
	pageData.SetTbName("tb_easy_models")
	return nil, 0
}

// NodeList 初始化列表
func (that EasyModels) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.SetButton("buttons", EasyApp.Button{
		ButtonName: "自定义按钮",
		Action:     "easy_models_buttons",
		ActionType: 2,
		LayerTitle: "模型页面按钮管理",
		ActionUrl:  config.AdminPath + "/easy_models_buttons/index",
		Class:      "brown",
		Icon:       "ri-radio-button-line",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})

	//!重置顶部按钮
	pageData.SetListTopBtns("add", "buttons")
	//!重置右侧按钮
	//重新设置编辑按钮
	pageData.SetButton("edit", EasyApp.Button{
		ButtonName: "编辑模型",
		Action:     "edit",
		ActionType: 2,
		ActionUrl:  "edit",
		Class:      "layui-btn-normal",
		Icon:       "ri-edit-box-line",
		Display:    "!item.btn_edit || item.btn_edit!='hide'",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//新增右侧字段管理按钮
	pageData.SetButton("fields", EasyApp.Button{
		ButtonName: "模型字段",
		Action:     "/easy_models_fields/index",
		ActionType: 2,
		LayerTitle: "模型字段管理",
		ActionUrl:  config.AdminPath + "/easy_models_fields/index",
		Class:      "rose",
		Icon:       "ri-table-line",
		Display:    "!item.btn_field || item.btn_field!='hide'",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//!重置右侧按钮
	pageData.SetListRightBtns("edit", "fields", "disable", "enable", "delete")

	pageData.ListColumnAdd("model_key", "模型Key", "text", nil)
	pageData.ListColumnAdd("model_name", "模型名称", "text", nil)
	//pageData.ListColumnAdd("table_name", "数据表名", "text", nil)
	pageData.ListColumnAdd("note", "备注", "text", nil)
	pageData.ListColumnAdd("allow_create", "新增按钮", "switch::text=显示|隐藏", nil)
	pageData.ListColumnAdd("allow_update", "修改按钮", "switch::text=显示|隐藏", nil)
	pageData.ListColumnAdd("allow_status", "状态按钮", "switch::text=显示|隐藏", nil)
	pageData.ListColumnAdd("allow_delete", "删除按钮", "switch::text=显示|隐藏", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	return nil, 0
}

// NodeForm 初始化表单
func (that EasyModels) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	//设置默认按钮tags列表
	buttonList := []map[string]interface{}{
		{"text": "add", "class": "tag-3"},
		{"text": "edit", "class": "tag-1"},
		{"text": "disable", "class": "tag-1"},
		{"text": "enable", "class": "tag-1"},
		{"text": "delete", "class": "tag-1"},
	}
	//设置自定义按钮的tag列表
	buttons, err := db.New().Table("tb_easy_models_buttons").
		Where("is_delete", 0).
		Where("status", 1).
		Fields("id,button_key").
		Get()
	if err != nil {
		logger.Error(err.Error())
		return err, 500
	}
	for _, btn := range buttons {
		buttonList = append(buttonList, map[string]interface{}{"text": btn["button_key"], "class": "tag-2"})
	}
	pageData.FormFieldsAdd("", "block", "基础信息", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("model_key", "text", "模型Key", "支持英文、数字、下划线，不能以下划线开头", "", true, nil, "", nil)
	pageData.FormFieldsAdd("model_name", "text", "模型名称", "模型名称即是管理页面的关键字", "", true, nil, "", nil)
	pageData.FormFieldsAdd("table_name", "text", "关联数据表名", "关联操作的数据表名称", "", true, nil, "", nil)
	pageData.FormFieldsAdd("order_type", "text", "排序方式", "列表页默认排序方式", "id desc", true, nil, "", nil)
	pageData.FormFieldsAdd("page_size", "number", "分页大小", "列表页每页的数据条数，最小值为1", "20", true, nil, "", nil)
	pageData.FormFieldsAdd("batch_action", "radio", "支持批量操作", "", "0", true, Models.OptionModels{}.ById(1,false), "", nil)
	pageData.FormFieldsAdd("note", "text", "模型备注", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("", "block", "页面元素", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("top_buttons", "tags", "顶部按钮", "", "[{'classes':'tag-3','text':'add'}]", false, buttonList, "", nil)
	rightBtn := "[{'classes':'tag-1','text':'edit'},{'classes':'tag-1','text':'disable'},{'classes':'tag-1','text':'enable'},{'classes':'tag-1','text':'delete'}]"
	pageData.FormFieldsAdd("right_buttons", "tags", "操作按钮", "", rightBtn, false, buttonList, "", nil)
	pageData.FormFieldsAdd("page_notice", "textarea", "页面公告", "列表页面顶部的提示信息", "", false, nil, "", nil)
	pageData.FormFieldsAdd("tabs_for_list", "textarea", "列表选项卡", "格式:tab名称|查询条件，每行一个", "", false, nil, "", nil)
	pageData.FormFieldsAdd("", "block", "高级用法", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("level_indent", "text", "字段按级缩进", "列表页支持按字段1的上下级关系缩进字段2，格式为:级别字段key:缩进字段key，例如pid:name", "", false, nil, "", nil)
	pageData.FormFieldsAdd("url_params", "textarea", "Url传参", "Url参数转为列表查询条件 并 透传顶部按钮链接 \n格式为 参数:数据库字段:默认值，例如：id:model_id:0\n默认值为空自动忽略，每行一个转换规则", "", false, nil, "", nil)
	return nil, 0
}
