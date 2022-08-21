/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyModelsFields
 * @Version: 1.0.0
 * @Date: 2022/5/2 08:37
 */

package adminHandle

import (
	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/EasyApp"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/util"
	"github.com/yqstech/gef/config"
	"strings"
)

type EasyModelsFields struct {
	Base
}

//列表数据类型
var listDataType = []map[string]interface{}{
	{"name": "普通文字(text)", "value": "text"},
	{"name": "数组/匹配关联数据源(array)", "value": "array"},
	{"name": "输入框(input)", "value": "input"},
	{"name": "Html代码(html)", "value": "html"},
	{"name": "图标(icon)", "value": "icon"},
	{"name": "颜色(color)", "value": "color"},
	{"name": "开关(switch)", "value": "switch"},
	{"name": "图片(image)", "value": "image"},
	{"name": "图片100尺寸(image100)", "value": "image100"},
	{"name": "图片60尺寸(image60)", "value": "image60"},
	{"name": "图片50尺寸(image50)", "value": "image50"},
	{"name": "图片40尺寸(image60)", "value": "image40"},
	{"name": "图片30尺寸(image30)", "value": "image30"},
}

//表单数据类型
var formDataType = []map[string]interface{}{
	{"name": "文字输入框(text)", "value": "text"},
	{"name": "文字输入框(短)(text-sm)", "value": "text-sm"},
	{"name": "文字输入框(很短)(text-xs)", "value": "text-xs"},
	{"name": "文字输入框(超短)(text-xxs)", "value": "text-xxs"},
	{"name": "文字输入框(禁用状态)(text-disabled)", "value": "text-disabled"},
	{"name": "数字输入框(number)", "value": "number"},
	{"name": "数字输入框(短)(number-sm)", "value": "number-sm"},
	{"name": "数字输入框(很短)(number-xs)", "value": "number-xs"},
	{"name": "数字输入框(超短)(number-xxs)", "value": "number-xxs"},
	{"name": "密码输入框(password)", "value": "password"},
	{"name": "文本域(textarea)", "value": "textarea"},
	{"name": "文本域(禁用状态)(textarea-disabled)", "value": "textarea-disabled"},
	{"name": "下拉选项(select)", "value": "select"},
	{"name": "下拉选项(禁用状态)(select-disabled)", "value": "select-disabled"},
	{"name": "单选(radio)", "value": "radio"},
	{"name": "多选(checkbox)", "value": "checkbox"},
	{"name": "多级多选(checkbox_level)", "value": "checkbox_level"},
	{"name": "图片预览(imageview)", "value": "imageview"},
	{"name": "图片上传(image)", "value": "image"},
	{"name": "图片集上传(images)", "value": "images"},
	{"name": "视频上传(video)", "value": "video"},
	{"name": "音频上传(audio)", "value": "audio"},
	{"name": "文件上传(file)", "value": "file"},
	{"name": "Icon图标选择器(icon)", "value": "icon"},
	{"name": "Color颜色选择器(color)", "value": "color"},
	{"name": "标签Tags(tags)", "value": "tags"},
	{"name": "富文本编辑(wangEditor)", "value": "wangEditor"},
	{"name": "日期时间选择器(datetime)", "value": "datetime"},
	{"name": "日期选择器(date)", "value": "date"},
	{"name": "时间选择器(time)", "value": "time"},
	{"name": "百度地图位置选择器(lnglat)", "value": "lnglat"},
}

//数据库数据变换规则
//表单输入数据和数据库存储的数据，需要按照规则进行转换
var dataTransRulesForDB = []map[string]interface{}{
	{"name": "两位小数(元)转整数(分)", "value": "yuan2fen"},
}

// NodeBegin 开始
func (that EasyModelsFields) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("模型字段管理")
	pageData.SetPageName("模型字段")
	pageData.SetTbName("tb_easy_models_fields")
	if pageData.GetHttpRequest().Method == "POST" {
		//自动同步数据库字段
		id := util.GetValue(pageData.GetHttpRequest(), "id")
		that.syncModelFields(util.String2Int(id))
		//不在表单的项目，都设置成非必填项
		db.New().Table("tb_easy_models_fields").
			Where("is_delete", 0).
			Where("allow_create", 0).
			Where("allow_update", 0).
			Update(map[string]interface{}{
				"is_must": 0,
			})
	}
	return nil, 0
}

// NodeList 初始化列表
func (that EasyModelsFields) NodeList(pageData *EasyApp.PageData) (error, int) {
	//隐藏新增按钮
	pageData.SetListTopBtns()
	//删除ID字段
	pageData.ListColumnClear()

	//!设置tab列表
	//获取页面地址，允许参数有参数id
	validUrl := util.UrlScreenParam(pageData.GetHttpRequest(), []string{"id"}, false, true)
	pageData.PageTabAdd("全部字段", validUrl)
	pageData.PageTabAdd("列表页预览", validUrl+"tab=1")
	pageData.PageTabAdd("新增页", validUrl+"tab=2")
	pageData.PageTabAdd("编辑页", validUrl+"tab=3")
	//获取第几页
	tabIndex := that.GetTabIndex(pageData, "tab")
	pageData.SetPageTabSelect(tabIndex)

	if tabIndex == 0 {
		//重新设置排序
		pageData.SetListOrder("index_num asc,id asc")

		//pageData.ListColumnAdd("field_key", "字段关键字", "text", nil)
		pageData.ListColumnAdd("field_name", "字段和提示", "html", nil)
		//pageData.ListColumnAdd("field_notice", "字段提示", "text", nil)
		pageData.ListColumnAdd("is_show_on_list", "列表页", "switch::text=显示|隐藏", nil)
		pageData.ListColumnAdd("allow_create", "新增页", "switch::text=显示|隐藏", nil)
		pageData.ListColumnAdd("allow_update", "修改页", "switch::text=显示|隐藏", nil)
		pageData.ListColumnAdd("is_must", "必填项", "switch::text=是|否", nil)
		pageData.ListColumnAdd("option_models_key", "选项集", "array", that.OptionModelsList())
		pageData.ListColumnAdd("dynamic_option_models_key", "联动选项集", "array", that.DynamicOptionModelsList())
		pageData.ListColumnAdd("default_value", "默认值", "input::width=60px", nil)
		pageData.ListColumnAdd("index_num", "排序", "input::type=number&width=50px", nil)
	} else if tabIndex == 1 {
		//重新设置排序
		pageData.SetListOrder("is_show_on_list desc,index_num asc,id asc")
		pageData.ListColumnAdd("field_name", "字段和提示", "html", nil)
		pageData.ListColumnAdd("is_show_on_list", "是否显示", "switch::text=显示|隐藏", nil)
		pageData.ListColumnAdd("data_type_on_list", "列表页数据类型", "array", listDataType)
		pageData.ListColumnAdd("data_type_command_on_list", "数据指令", "text", nil)
		pageData.ListColumnAdd("option_models_key", "选项集", "array", that.OptionModelsList())
		pageData.ListColumnAdd("index_num", "排序值", "input::type=number&width=50px", nil)
	} else if tabIndex == 2 {
		//重新设置排序
		pageData.SetListOrder("allow_create desc,index_num asc,id asc")

		pageData.ListColumnAdd("field_name", "字段和提示", "html", nil)
		pageData.ListColumnAdd("allow_create", "是否显示", "switch::text=显示|隐藏", nil)
		pageData.ListColumnAdd("is_must", "必填项", "switch::text=是|否", nil)
		pageData.ListColumnAdd("data_type_on_create", "数据类型", "array", formDataType)
		pageData.ListColumnAdd("option_models_key", "选项集", "array", that.OptionModelsList())
		pageData.ListColumnAdd("dynamic_option_models_key", "联动选项集", "array", that.DynamicOptionModelsList())
		pageData.ListColumnAdd("default_value", "默认值", "input::width=60px", nil)
		pageData.ListColumnAdd("index_num", "排序值", "input::type=number&width=50px", nil)

	} else if tabIndex == 3 {
		//重新设置排序
		pageData.SetListOrder("allow_update desc,index_num asc,id asc")

		pageData.ListColumnAdd("field_name", "字段和提示", "html", nil)
		pageData.ListColumnAdd("allow_update", "是否显示", "switch::text=显示|隐藏", nil)
		pageData.ListColumnAdd("is_must", "必填项", "switch::text=是|否", nil)
		pageData.ListColumnAdd("data_type_on_update", "数据类型", "array", formDataType)
		pageData.ListColumnAdd("option_models_key", "选项集", "array", that.OptionModelsList())
		pageData.ListColumnAdd("dynamic_option_models_key", "联动选项集", "array", that.DynamicOptionModelsList())
		pageData.ListColumnAdd("default_value", "默认值", "input::width=60px", nil)
		pageData.ListColumnAdd("index_num", "排序值", "input::type=number&width=50px", nil)
	}
	return nil, 0

}

// NodeListCondition 修改查询条件
func (that EasyModelsFields) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	modelID := 0
	modelId := util.GetValue(pageData.GetHttpRequest(), "id")
	if modelId != "" {
		modelID = util.String2Int(modelId)
		//追加查询条件
		condition = append(condition, []interface{}{
			"model_id", "=", modelID,
		})
	}

	return condition, nil, 0
}

// NodeListData 重写列表数据
func (that EasyModelsFields) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	for k, v := range data {
		data[k]["field_name"] = "<strong style='color:#FF9900; font-size:13px'>" + v["field_key"].(string) + "</strong><br>" +
			v["field_name"].(string) +
			"<br><span style='color:#0099CC; font-size:13px'>" + v["field_notice"].(string) + "</span>"
	}
	return data, nil, 0
}

// NodeForm 初始化表单
func (that EasyModelsFields) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	//pageData.FormFieldsAdd("model_id", "select-disabled", "所属模型", "", "", false, that.EasyModels(), "", nil)
	//pageData.FormFieldsAdd("field_key", "text-disabled", "数据表字段", "必须和数据表内的字段一致", "", false, nil, "", nil)
	//pageData.FormFieldsAdd("field_name", "text-disabled", "字段名称", "字段自定义名称", "", false, nil, "", nil)
	//pageData.FormFieldsAdd("field_notice", "text-disabled", "提示信息", "表单数据项的提示信息", "", false, nil, "", nil)
	//pageData.FormFieldsAdd("is_show_on_list", "radio", "列表页显示", "是否在列表页显示此字段", "1", true, Models.DefaultIsOrNot, "", nil)

	//数据类型和选项
	pageData.FormFieldsAdd("", "block", "字段基础信息", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("data_type_on_list", "select", "列表数据类型", "列表页显示的组件", "text", false, listDataType, "", nil)
	pageData.FormFieldsAdd("data_type_command_on_list", "text", "数据指令", "列表页字段组件的配置信息，switch组件/input组件需要", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type_on_list=='switch' || formFields.data_type_on_list=='input'",
	})
	pageData.FormFieldsAdd("data_type_on_create", "select", "新增页数据类型", "新增页的数据类型", "text", false, formDataType, "", nil)
	pageData.FormFieldsAdd("data_type_on_update", "select", "编辑页数据类型", "修改页的数据类型", "text", false, formDataType, "", nil)
	pageData.FormFieldsAdd("option_models_key", "select", "关联选项集", "", "", false, that.OptionModelsList(), "", nil)
	pageData.FormFieldsAdd("option_beautify", "radio", "选项集美化", "", "1", false, []map[string]interface{}{
		{"name": "是", "value": "1"},
		{"name": "否", "value": "0"},
		{"name": "仅列表美化", "value": "2"},
	}, "", map[string]interface{}{
		"if": "formFields.option_models_key!=''",
	})
	pageData.FormFieldsAdd("default_value", "text", "字段默认值", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("save_trans_rule", "select", "存储格式转换", "", "", false, dataTransRulesForDB, "", nil)

	pageData.FormFieldsAdd("", "block", "联动设置", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type_on_create=='select' || formFields.data_type_on_update=='select' || formFields.data_type_on_create=='radio' || formFields.data_type_on_update=='radio' ",
	})
	pageData.FormFieldsAdd("watch_fields", "text", "监听字段", "联动监听的字段，多个字段用英文逗号,分割", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type_on_create=='select' || formFields.data_type_on_update=='select' || formFields.data_type_on_create=='radio' || formFields.data_type_on_update=='radio' ",
	})
	pageData.FormFieldsAdd("dynamic_option_models_key", "select", "关联选项集", "", "", false, that.DynamicOptionModelsList(), "", map[string]interface{}{
		"if": "formFields.data_type_on_create=='select' || formFields.data_type_on_update=='select' || formFields.data_type_on_create=='radio' || formFields.data_type_on_update=='radio' ",
	})

	//列表列装饰
	pageData.FormFieldsAdd("", "block", "列表页装饰", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("set_as_tabs", "radio", "选项集设为Tabs", "", "0", false, Models.OptionModels{}.ByKey("is", false), "", map[string]interface{}{
		"if": "formFields.option_models_key!=''",
	})

	pageData.FormFieldsAdd("field_name_reset", "text", "重置列标题", "重新设置列表页中此字段的标题", "", false, nil, "", nil)
	pageData.FormFieldsAdd("field_style_reset", "text", "重置列样式", "设置列表页此列的样式，例如：width:20%", "", false, nil, "", nil)
	pageData.FormFieldsAdd("field_augment", "textarea", "美化原始数据", "支持html代码，列表数据类型需要改为html，{{this}}代表原数据", "", false, nil, "", nil)
	pageData.FormFieldsAdd("attach_to_field", "text", "多字段合并显示", "将此字段数据合并到其他字段", "", false, nil, "", nil)

	//数据分组
	pageData.FormFieldsAdd("", "block", "表单页装饰", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("option_indent", "radio", "选项按上下级缩进", "", "0", false, Models.OptionModels{}.ByKey("is", false), "", map[string]interface{}{
		"if": "formFields.option_models_key!=''",
	})
	pageData.FormFieldsAdd("group_title", "text", "创建一个分组", "填写分组名称，会从当前项的前面创建一个新分组", "", false, nil, "", nil)
	pageData.FormFieldsAdd("expand_if", "text", "联动if条件", "例如：formFields.data_type==1", "", false, nil, "", nil)
	return nil, 0
}

// NodeFormData 表单显示前修改数据
func (that EasyModelsFields) NodeFormData(pageData *EasyApp.PageData, data gorose.Data, id int64) (gorose.Data, error, int) {
	if id > 0 {
	}
	return data, nil, 0
}

// NodeSaveData 表单保存数据前使用
func (that EasyModelsFields) NodeSaveData(pageData *EasyApp.PageData, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	if postData["set_as_tabs"] == "1" {
		_, err := db.New().Table("tb_easy_models_fields").
			Where("model_id", oldData["model_id"]).
			Where("id", "!=", oldData["id"]).
			Where("is_delete", 0).
			Update(map[string]interface{}{
				"set_as_tabs": 0,
			})
		if err != nil {
			return nil, err, 500
		}
	}
	return postData, nil, 0
}

//默认列表需要显示的字段
var defaultListShowFields = []interface{}{"id", "status"}

//禁止新增的字段——新增页不显示
var defaultNotAllowCreateFields = []interface{}{"is_delete", "update_time", "create_time", "id", "status"}

//默认禁止修改的字段——编辑页不显示
var defaultNotAllowUpdateFields = []interface{}{"is_delete", "update_time", "create_time", "id", "status"}

//同步模型字段
func (that EasyModelsFields) syncModelFields(easyModelId int) {
	//查询模型信息
	easyModelInfo, err := db.New().Table("tb_easy_models").
		Where("id", easyModelId).
		Where("is_delete", 0).First()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if easyModelInfo == nil {
		return
	}
	if easyModelInfo["table_name"].(string) == "" {
		return
	}
	//得到模型数据表名称
	tableName := easyModelInfo["table_name"].(string)

	//查询数据表实时字段信息
	query, err := db.New().Query("select COLUMN_NAME,COLUMN_COMMENT,COLUMN_DEFAULT,"+
		"COLUMN_TYPE from information_schema.COLUMNS where table_name = ? and table_schema = ? order by ordinal_position",
		tableName,
		config.DbName)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	//查询模型字段列表
	fields, err := db.New().Table("tb_easy_models_fields").
		Where("model_id", easyModelId).
		Where("is_delete", 0).
		Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//字段列表转成以fieldKey为键的map数据
	fieldsMap := map[string]gorose.Data{}
	for _, field := range fields {
		fieldsMap[field["field_key"].(string)] = field
	}
	timeNow := util.TimeNow()
	//对比更新模型字段信息
	for index, fieldInfo := range query {
		//字段key
		fieldKey := fieldInfo["COLUMN_NAME"].(string)
		//字段名称和备注，直接同步数据库的备注信息
		fieldName := ""
		fieldNotice := ""
		commands := strings.Split(fieldInfo["COLUMN_COMMENT"].(string), "|")
		if len(commands) > 0 {
			fieldName = commands[0]
		}
		if len(commands) > 1 {
			fieldNotice = commands[1]
		}
		//数据类型
		dataTypeOnList := "text"
		dataTypeOnEdit := "text"
		//关联数据
		OptionModelKey := Models.FieldMatchOptionModelsKey(fieldKey)
		if OptionModelKey != "" {
			//匹配到关联数据，自动设置类型
			dataTypeOnList = "array"
			dataTypeOnEdit = "select"
		}
		//是否在列表显示
		isShowOnList := 0
		if util.IsInArray(fieldKey, defaultListShowFields) {
			isShowOnList = 1
		}
		//是否默认设置成锁定
		allowCreate := 1
		if util.IsInArray(fieldKey, defaultNotAllowCreateFields) {
			allowCreate = 0
		}
		//是否默认设置成锁定
		allowUpdate := 1
		if util.IsInArray(fieldKey, defaultNotAllowUpdateFields) {
			allowUpdate = 0
		}
		//是否存在字段
		if field, ok := fieldsMap[fieldKey]; ok {
			db.New().Table("tb_easy_models_fields").
				Where("id", field["id"]).
				Update(map[string]interface{}{
					"field_name":   fieldName,
					"field_notice": fieldNotice,
					"update_time":  timeNow,
				})
			fieldsMap[fieldKey]["sync_tag"] = true
		} else {
			db.New().Table("tb_easy_models_fields").Insert(map[string]interface{}{
				"model_id":            easyModelId,
				"field_key":           fieldKey,
				"field_name":          fieldName,
				"field_notice":        fieldNotice,
				"index_num":           index + 1,
				"option_models_key":   OptionModelKey, //自动匹配选择数据源
				"is_show_on_list":     isShowOnList,
				"data_type_on_list":   dataTypeOnList, //设置默认列表显示组件
				"allow_create":        allowCreate,
				"allow_update":        allowUpdate,
				"data_type_on_create": dataTypeOnEdit, //设置默认表单显示组件
				"data_type_on_update": dataTypeOnEdit, //设置默认表单显示组件
			})
		}
	}

	for _, field := range fieldsMap {
		//未标记的都删除
		if _, ok2 := field["sync_tag"]; !ok2 {
			db.New().Table("tb_easy_models_fields").
				Where("id", field["id"]).
				Update(map[string]interface{}{
					"update_time": timeNow,
					"is_delete":   1,
				})
		}
	}
}
