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
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/config"
	"github.com/yqstech/gef/util"
	"net/http"
)

type EasyModels struct {
	Base
}

// NodeInit 初始化
func (that *EasyModels) NodeInit(pageBuilder *builder.PageBuilder) {
	//注册handle
	that.NodePageActions["export_insert_data"] = that.ExportInsertData
	that.NodePageActions["update_all_models_fields"] = that.UpdateAllFields
	that.NodePageActions["export_all_data"] = that.ExportAllData
}

// NodeBegin 开始
func (that EasyModels) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("后台模型 EasyModelHandle")
	pageBuilder.SetPageName("模型")
	pageBuilder.SetTbName("tb_easy_models")
	return nil, 0
}

// NodeList 初始化列表
func (that EasyModels) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetButton("buttons", builder.Button{
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
	//新增顶部按钮
	pageBuilder.SetButton("update_all_models_fields", builder.Button{
		ButtonName: "刷新字段",
		Action:     "/easy_models/update_all_models_fields",
		ActionType: 2,
		LayerTitle: "刷新全部模型字段",
		ActionUrl:  config.AdminPath + "/easy_models/update_all_models_fields",
		Class:      "rose",
		Icon:       "ri-refresh-fill",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	pageBuilder.SetButton("export_all_data", builder.Button{
		ButtonName: "导出未禁用",
		Action:     "/easy_models/export_all_data",
		ActionType: 2,
		LayerTitle: "导出全部接口模型",
		ActionUrl:  config.AdminPath + "/easy_models/export_all_data",
		Class:      "black",
		Icon:       "ri-braces-fill",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//!重置顶部按钮
	pageBuilder.SetListTopBtns("add", "buttons", "update_all_models_fields", "export_all_data")

	//!重置右侧按钮
	//重新设置编辑按钮
	pageBuilder.SetButton("edit", builder.Button{
		ButtonName: "",
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
	pageBuilder.SetButton("fields", builder.Button{
		ButtonName: "字段",
		Action:     "/easy_models_fields/index",
		ActionType: 2,
		LayerTitle: "模型字段管理",
		ActionUrl:  config.AdminPath + "/easy_models_fields/index",
		Class:      "rose",
		Icon:       "ri-table-line",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//新增右侧字段管理按钮
	pageBuilder.SetButton("search_form", builder.Button{
		ButtonName: "搜索",
		Action:     "/easy_models_search_form/index",
		ActionType: 2,
		LayerTitle: "模型搜索表单管理",
		ActionUrl:  config.AdminPath + "/easy_models_search_form/index",
		Class:      "cyan",
		Icon:       "ri-search-2-line",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//导出结构
	pageBuilder.SetButton("export_insert_data", builder.Button{
		ButtonName: "",
		Action:     "/easy_models/export_insert_data",
		ActionType: 2,
		LayerTitle: "后台模型导出成内置数据",
		ActionUrl:  config.AdminPath + "/easy_models/export_insert_data",
		Class:      "black",
		Icon:       "ri-braces-fill",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	pageBuilder.SetButtonName("disable", "")
	pageBuilder.SetButtonName("enable", "")
	pageBuilder.SetButtonName("delete", "")
	//!重置右侧按钮
	pageBuilder.SetListRightBtns("edit", "fields", "search_form", "export_insert_data", "disable", "enable", "delete")

	pageBuilder.SetListOrder("id asc")
	pageBuilder.ListColumnAdd("model_key", "模型Key", "text", nil)
	pageBuilder.ListColumnAdd("model_name", "模型名称", "text", nil)
	//pageBuilder.ListColumnAdd("table_name", "数据表名", "text", nil)
	pageBuilder.ListColumnAdd("note", "备注", "text", nil)
	pageBuilder.ListColumnAdd("allow_create", "新增按钮", "switch::text=显示|隐藏", nil)
	pageBuilder.ListColumnAdd("allow_update", "修改按钮", "switch::text=显示|隐藏", nil)
	pageBuilder.ListColumnAdd("allow_status", "状态按钮", "switch::text=显示|隐藏", nil)
	pageBuilder.ListColumnAdd("allow_delete", "删除按钮", "switch::text=显示|隐藏", nil)
	pageBuilder.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ByKey("status", true))
	return nil, 0
}

// NodeForm 初始化表单
func (that EasyModels) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
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
	pageBuilder.FormFieldsAdd("", "block", "基础信息", "", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("model_key", "text", "模型Key", "支持英文、数字、下划线，不能以下划线开头", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("model_name", "text", "模型名称", "模型名称即是管理页面的关键字", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("table_name", "text", "关联数据表名", "关联操作的数据表名称", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("order_type", "text", "排序方式", "列表页默认排序方式", "id desc", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("page_size", "number", "分页大小", "列表页每页的数据条数，最小值为1", "20", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("batch_action", "radio", "支持批量操作", "", "0", true, Models.OptionModels{}.ByKey("is", false), "", nil)
	pageBuilder.FormFieldsAdd("note", "text", "模型备注", "", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("", "block", "页面元素", "", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("top_buttons", "tags", "顶部按钮", "", "[{'classes':'tag-3','text':'add'}]", false, buttonList, "", nil)
	rightBtn := "[{'classes':'tag-1','text':'edit'},{'classes':'tag-1','text':'disable'},{'classes':'tag-1','text':'enable'},{'classes':'tag-1','text':'delete'}]"
	pageBuilder.FormFieldsAdd("right_buttons", "tags", "操作按钮", "", rightBtn, false, buttonList, "", nil)
	pageBuilder.FormFieldsAdd("page_notice", "textarea", "页面公告", "列表页面顶部的提示信息", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("tabs_for_list", "textarea", "列表选项卡", "格式:tab名称|查询条件|搜索表单联动参数对(例如:a=1&b=2)，每行一个", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("", "block", "高级用法", "", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("level_indent", "text", "字段按级缩进", "列表页支持按字段1的上下级关系缩进字段2，格式为:级别字段key:缩进字段key，例如pid:name", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("url_params", "textarea", "Url传参", "Url参数转为列表查询条件 并 透传顶部按钮链接 \n格式为 参数:数据库字段:默认值，例如：id:model_id:0\n默认值为空自动忽略，每行一个转换规则", "", false, nil, "", nil)
	return nil, 0
}

// UpdateAllFields 更新所有字段
func (that EasyModels) UpdateAllFields(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	models, err := db.New().Table("tb_easy_models").Where("is_delete", 0).Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ModelsFields := EasyModelsFields{}
	for _, m := range models {
		ModelsFields.syncModelFields(util.Int642Int(m["id"].(int64)))
		fmt.Fprint(w, "刷新模型【"+m["model_name"].(string)+"】成功！\n")
	}
	fmt.Fprint(w, "操作成功！")
}

// ExportAllData 导出所有数据
func (that EasyModels) ExportAllData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	models, err := db.New().
		Table("tb_easy_models").
		Where("is_delete", 0).
		Where("status", 1).
		Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, m := range models {
		that.ExportModelsAndFileds(w, m)
	}
}

// ExportInsertData 导出内置数据
func (that EasyModels) ExportInsertData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := util.GetValue(r, "id")

	easyModel, err := db.New().Table("tb_easy_models").Where("id", id).First()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if easyModel == nil {
		return
	}
	that.ExportModelsAndFileds(w, easyModel)
}

func (that EasyModels) ExportModelsAndFileds(w http.ResponseWriter, easyModel gorose.Data) {

	fmt.Fprint(w, "//! 后台模型【"+easyModel["model_name"].(string)+"】")

	fieldsContent := that.ExportModelFileds(easyModel["id"].(int64))
	searchFormContent := that.ExportModelSearchForm(easyModel["id"].(int64))
	delete(easyModel, "id")
	delete(easyModel, "create_time")
	delete(easyModel, "update_time")
	delete(easyModel, "is_delete")
	content := `
{
	TableName: "tb_easy_models",
	Condition: [][]interface{}{{"model_key", "` + easyModel["model_key"].(string) + `"}},
	Data: map[string]interface{}` + util.JsonEncode(easyModel) + `,
	Children:[]interface{}{
		//!后台模型的字段
		` + fieldsContent + `
		//!后台模型的搜索表单
		` + searchFormContent + `
	},
},
`
	fmt.Fprint(w, content)
}

// ExportModelFileds 导出字段
func (that EasyModels) ExportModelFileds(id int64) string {
	Fields, err := db.New().
		Table("tb_easy_models_fields").
		Where("is_delete", 0).
		Where("model_id", id).
		Order("index_num,id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	content := ""
	for index, Item := range Fields {
		//删除部分字段
		delete(Item, "id")
		delete(Item, "create_time")
		delete(Item, "update_time")
		delete(Item, "is_delete")
		//按顺序添加排序字段
		Item["index_num"] = index + 1
		//标记上级ID
		Item["model_id"] = "__PID__"
		content += `
		gef.InsideData{
			TableName: "tb_easy_models_fields",
			Condition: [][]interface{}{{"model_id", "__PID__"},{"field_key", "` + Item["field_key"].(string) + `"}},
			Data: map[string]interface{}` + util.JsonEncode(Item) + `,
		},
`
	}
	return content
}

// ExportModelSearchForm 导出搜索表单
func (that EasyModels) ExportModelSearchForm(id int64) string {
	SearchForm, err := db.New().
		Table("tb_easy_models_search_form").
		Where("is_delete", 0).
		Where("model_id", id).
		Order("index_num,id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	content := ""
	for index, Item := range SearchForm {
		//删除部分字段
		delete(Item, "id")
		delete(Item, "create_time")
		delete(Item, "update_time")
		delete(Item, "is_delete")
		//按顺序添加排序字段
		Item["index_num"] = index + 1
		//标记上级ID
		Item["model_id"] = "__PID__"
		content += `
		gef.InsideData{
			TableName: "tb_easy_models_search_form",
			Condition: [][]interface{}{{"model_id", "__PID__"},{"search_key", "` + Item["search_key"].(string) + `"}},
			Data: map[string]interface{}` + util.JsonEncode(Item) + `,
		},
`
	}
	return content
}
