/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: OptionModels
 * @Version: 1.0.0
 * @Date: 2022/5/3 09:55
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"net/http"
	"strings"
)

type OptionModels struct {
	Base
}

var dataTypes = []map[string]interface{}{
	{"name": "静态json数据", "value": "0"},
	{"name": "查询数据表", "value": "1"},
}

// PageInit 初始化
func (that OptionModels) PageInit(pageData *EasyApp.PageData) {
	//注册handle
	pageData.ActionAdd("dynamic", that.Dynamic)
}

// NodeBegin 开始
func (that OptionModels) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("选项集管理")
	pageData.SetPageName("选项集")
	pageData.SetTbName("tb_option_models")
	return nil, 0
}

// NodeList 初始化列表
func (that OptionModels) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.SetListOrder("index_num,id asc")
	pageData.ListColumnAdd("name", "名称", "text", nil)
	pageData.ListColumnAdd("data_type", "数据类型", "array", dataTypes)
	pageData.ListColumnAdd("static_data", "静态数据", "html", nil)
	pageData.ListColumnAdd("table_name", "数据表", "text", nil)
	pageData.ListColumnAdd("value_field", "值字段", "text", nil)
	pageData.ListColumnAdd("name_field", "名称字段", "text", nil)
	//pageData.ListColumnAdd("color_field", "颜色字段", "text", nil)
	//pageData.ListColumnAdd("icon_field", "图标字段", "text", nil)
	pageData.ListColumnAdd("index_num", "排序", "input::type=number&width=50px", nil)
	return nil, 0
}

// NodeListData 重写列表数据
func (that OptionModels) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	for k, v := range data {
		staticData := []map[string]interface{}{}
		if v["static_data"].(string) != "" {
			util.JsonDecode(v["static_data"].(string), &staticData)
			transData := []string{}
			for _, opt := range staticData {
				transData = append(transData, opt["name"].(string)+" : "+util.Interface2String(opt["value"]))
			}
			data[k]["static_data"] = strings.Join(transData, "<br>")
		}
	}
	return data, nil, 0
}

// NodeForm 初始化表单
func (that OptionModels) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("name", "text", "名称", "", "", true, nil, "", nil)
	pageData.FormFieldsAdd("data_type", "radio", "数据类型", "", "0", true, dataTypes, "", nil)
	pageData.FormFieldsAdd("static_data", "textarea", "静态数据", "", "[{\"name\":\"是\",\"value\":\"1\"},{\"name\":\"否\",\"value\":\"0\"}]", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==0",
	})
	pageData.FormFieldsAdd("default_data", "textarea", "默认数据", "和静态数据格式相同", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("", "block", "配置数据表", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("table_name", "text", "数据表名称", "tb_", "tb_", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("value_field", "text", "值字段", "查询到的数据作为值", "id", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("name_field", "text", "名称字段", "查询到的数据作为名称", "name", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("select_order", "text", "查询排序", "数据表查询的排序方式", "id asc", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("select_where", "text", "补充查询条件", "补充数据表查询条件", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("dynamic_params", "textarea", "支持的动态参数", "用来做数据联动的参数设置，程序根据设置的字段，查询post参数\n格式为 监听参数:选项集数据表字段:默认值，例如：group_id:group_id:0\n默认值为空自动忽略，每行一个转换规则", "", false, nil,
		"",
		map[string]interface{}{
			"if": "formFields.data_type==1",
		})
	pageData.FormFieldsAdd("", "block", "选项图标、颜色、缩进", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("parent_field", "text", "上级字段", "设置上级字段会查询原数据库上级字段信息，并在数据中增加pid一项", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("color_field", "text", "颜色字段", "设置某字段值作为颜色值", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("icon_field", "text", "Icon字段", "设置某字段值作为Icon图标值", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("color_array", "textarea", "颜色组", "选项列表对应的颜色列表，使用英文逗号分割，例如:#FFFFFF,#FFF660", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("icon_array", "textarea", "Icon图标组", "选项列表对应的图标列表，使用英文逗号分割，例如:ri-send-plane-line,ri-chat-3-line", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("", "block", "自动匹配与排序", "", "", false, nil, "", nil)
	//默认排序值
	indexNum := 1
	if id == 0 {
		num, err := db.New().Table("tb_option_models").Where("is_delete", 0).Count()
		if err != nil {
			logger.Error(err.Error())
			return err, 500
		}
		indexNum = util.Int642Int(num) + 1
	}
	pageData.FormFieldsAdd("match_fields", "textarea", "自动匹配字段", "每个字段占一行，支持全匹配字段和半匹配字段,例如is_*", "", false, nil, "", nil)
	pageData.FormFieldsAdd("index_num", "text", "排序", "", util.Int2String(indexNum), true, nil, "", nil)
	
	return nil, 0
}

// NodeSaveData 表单保存数据前使用
func (that OptionModels) NodeSaveData(pageData *EasyApp.PageData, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	if postData["data_type"].(string) == "0" {
		postData["table_name"] = ""
		postData["value_field"] = ""
		postData["name_field"] = ""
	} else {
		postData["static_data"] = ""
	}
	return postData, nil, 0
}

// Dynamic 动态获取选项集
func (that OptionModels) Dynamic(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	optionModelId := util.PostValue(r, "_dynamic_option_model_id")
	if optionModelId == "" {
		that.ApiResult(w, 201, "参数不全", nil)
		return
	}
	DynamicParams := Models.OptionModels{}.DynamicParams(util.String2Int(optionModelId))
	var wheres []string
	for _, dp := range DynamicParams {
		v := util.PostValue(r, dp.ParamKey)
		if v == "" {
			wheres = append(wheres, dp.FieldKey+" = '"+dp.DefValue+"'")
		} else {
			wheres = append(wheres, dp.FieldKey+" = '"+v+"'")
		}
	}
	ops := Models.OptionModels{}.Select(util.String2Int(optionModelId), strings.Join(wheres, " and "), false)
	that.ApiResult(w, 200, "success", ops)
}
