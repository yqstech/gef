package adminHandle

import (
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"
)

type EasyModelsSearchForm struct {
	Base
}

// 搜索表单类型
var searchDataType = []map[string]interface{}{
	{"name": "文字输入框(text)", "value": "text"},
	{"name": "文字输入框(短)(text-sm)", "value": "text-sm"},
	{"name": "数字输入框(number)", "value": "number"},
	{"name": "数字输入框(短)(number-sm)", "value": "number-sm"},
	{"name": "下拉单选框(select)", "value": "select"},
	//{"name": "勾选框(checkbox)", "value": "checkbox"},
	//{"name": "日期选择器(date)", "value": "date"},
	//{"name": "时间选择器(time)", "value": "time"},
	{"name": "日期时间选择器(datetime)", "value": "datetime"},
}

// 搜索表单类型
var searchMatchType = []map[string]interface{}{
	{"name": "=", "value": "="},
	{"name": ">=", "value": ">="},
	{"name": "<=", "value": "<="},
	{"name": "模糊查询(like)", "value": "like"},
	//{"name": "自定义", "value": "-"}, //checkbox勾选后自定义sql查询语句
}

// NodeBegin 开始
func (that EasyModelsSearchForm) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("搜索表单项")
	pageBuilder.SetPageName("表单项")
	pageBuilder.SetTbName("tb_easy_models_search_form")
	return nil, 0
}

// NodeList 初始化列表
func (that EasyModelsSearchForm) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	modelId := util.GetValue(pageBuilder.GetHttpRequest(), "id")
	pageBuilder.SetButtonActionUrl("add", "model_id="+modelId, true)
	pageBuilder.SetButtonActionUrl("edit", "model_id="+modelId, true)
	//!设置tab列表
	pageBuilder.SetListOrder("id asc")
	pageBuilder.ListColumnClear()
	pageBuilder.ListColumnAdd("search_key", "搜索项标识", "text", nil)
	pageBuilder.ListColumnAdd("search_name", "搜索项名称", "text", nil)
	pageBuilder.ListColumnAdd("data_type", "表单组件", "array", searchDataType)
	pageBuilder.ListColumnAdd("option_models_key", "选项集", "array", that.OptionModelsList())
	pageBuilder.ListColumnAdd("search_fields", "搜索字段", "text", nil)
	pageBuilder.ListColumnAdd("match_type", "匹配规则", "array", searchMatchType)
	pageBuilder.ListColumnAdd("default_value", "默认值", "input::width=60px", nil)
	pageBuilder.ListColumnAdd("index_num", "排序", "input::type=number&width=50px", nil)
	pageBuilder.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ByKey("status", true))
	return nil, 0
}

// NodeListCondition 修改查询条件
func (that EasyModelsSearchForm) NodeListCondition(pageBuilder *builder.PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {
	//追加查询条件
	modelId := util.GetValue(pageBuilder.GetHttpRequest(), "id")
	if modelId != "" {
		modelID := util.String2Int(modelId)
		condition = append(condition, []interface{}{
			"model_id", "=", modelID,
		})
	}
	return condition, nil, 0
}

// NodeForm 初始化表单
func (that EasyModelsSearchForm) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
	modelId := util.GetValue(pageBuilder.GetHttpRequest(), "model_id")
	if modelId != "" {
		pageBuilder.FormFieldsAdd("model_id", "hidden", "", "", modelId, true, nil, "", nil)
		pageBuilder.FormFieldsAdd("search_key", "text-sm", "搜索项标识", "模型内唯一搜索项", "", true, nil, "", nil)
		pageBuilder.FormFieldsAdd("search_name", "text-sm", "搜索项名称", "输入框前面的显示的文字", "", true, nil, "", nil)
		pageBuilder.FormFieldsAdd("placeholder", "text-sm", "提示信息", "输入框内显示的提示信息", "", false, nil, "", nil)
		pageBuilder.FormFieldsAdd("data_type", "select-sm", "表单组件", "设置数据类型和使用的组件", "text", true, searchDataType, "", nil)
		pageBuilder.FormFieldsAdd("option_models_key", "select-xs", "关联选项集", "下拉选项关联选项集", "", false, that.OptionModelsList(), "", map[string]interface{}{
			"if": "formFields.data_type=='select' || formFields.data_type=='select-sm'",
		})
		pageBuilder.FormFieldsAdd("search_fields", "checkbox", "搜索字段", "搜索关联的字段", "", true, that.GetEasyModelsFields(modelId), "", nil)
		pageBuilder.FormFieldsAdd("match_type", "radio", "匹配类型", "查询数据匹配类型", "=", true, searchMatchType, "", nil)
		pageBuilder.FormFieldsAdd("default_value", "text-xxs", "默认值", "搜索填充默认值", "", false, nil, "", nil)
		pageBuilder.FormFieldsAdd("style", "text", "附加样式", "输入框、下拉框等组件调整样式", "", false, nil, "", nil)
	}

	return nil, 0
}
