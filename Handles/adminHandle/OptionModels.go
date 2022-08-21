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
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/EasyApp"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	util2 "github.com/yqstech/gef/Utils/util"
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
	pageData.ActionAdd("export_insert_data", that.ExportInsertData)
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
	pageData.SetButton("export_insert_data", EasyApp.Button{
		ButtonName: "导出数据",
		Action:     "export_insert_data",
		ActionType: 2,
		LayerTitle: "后台模型导出成内置数据",
		ActionUrl:  "export_insert_data",
		Class:      "black",
		Icon:       "ri-braces-fill",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//!重置顶部按钮
	pageData.SetListTopBtns("add", "export_insert_data")

	pageData.SetListOrder("index_num,id asc")
	pageData.ListColumnAdd("unique_key", "标识符", "text", nil)
	pageData.ListColumnAdd("name", "名称", "text", nil)
	pageData.ListColumnAdd("data_type", "数据类型", "array", dataTypes)
	pageData.ListColumnAdd("static_data", "静态数据", "html", nil)
	pageData.ListColumnAdd("table_name", "数据表", "text", nil)
	pageData.ListColumnAdd("value_field", "值字段", "text", nil)
	pageData.ListColumnAdd("name_field", "名称字段", "text", nil)
	pageData.ListColumnAdd("parent_field", "上级字段", "text", nil)
	pageData.ListColumnAdd("children_option_model_key", "下级选项集", "array", that.ChildrenOptionModelsList(""))
	pageData.ListColumnAdd("index_num", "排序", "input::type=number&width=50px", nil)
	return nil, 0
}

// NodeListData 重写列表数据
func (that OptionModels) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	for k, v := range data {
		staticData := []map[string]interface{}{}
		if v["static_data"].(string) != "" {
			util2.JsonDecode(v["static_data"].(string), &staticData)
			transData := []string{}
			for _, opt := range staticData {
				transData = append(transData, opt["name"].(string)+" : "+util2.Interface2String(opt["value"]))
			}
			data[k]["static_data"] = strings.Join(transData, "<br>")
		}
	}
	return data, nil, 0
}

// NodeForm 初始化表单
func (that OptionModels) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("unique_key", "text-sm", "标识符", "选项集唯一的识别关键字", "", true, nil, "", nil)
	pageData.FormFieldsAdd("name", "text-sm", "名称", "", "", true, nil, "", nil)
	pageData.FormFieldsAdd("data_type", "radio", "数据类型", "", "0", true, dataTypes, "", nil)
	pageData.FormFieldsAdd("static_data", "textarea", "静态数据", "", "[{\"name\":\"是\",\"value\":\"1\"},{\"name\":\"否\",\"value\":\"0\"}]", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==0",
	})
	pageData.FormFieldsAdd("default_data", "textarea", "默认数据", "和静态数据格式相同", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	//配置数据表
	pageData.FormFieldsAdd("", "block", "配置数据表", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("table_name", "text-xs", "数据表名称", "tb_", "tb_", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("value_field", "text-xs", "值字段", "查询到的数据作为值", "id", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("name_field", "text-xs", "名称字段", "查询到的数据作为名称", "name", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("select_order", "text-xs", "查询排序", "数据表查询的排序方式", "id asc", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageData.FormFieldsAdd("select_where", "text-sm", "补充查询条件", "补充数据表查询条件", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})

	//联动设置
	pageData.FormFieldsAdd("", "block", "选项集联动", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})

	pageData.FormFieldsAdd("dynamic_params", "textarea", "联动配置", "用来做数据联动的参数设置，程序根据设置的字段，查询post参数\n格式为 监听参数:选项集数据表字段:默认值，例如：group_id:group_id:0\n默认值为空自动忽略，每行一个转换规则", "", false,
		nil, "", map[string]interface{}{
			"if": "formFields.data_type==1",
		})

	//多级、数据转换
	pageData.FormFieldsAdd("", "block", "多级、数据转换", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("parent_field", "text-xs", "上级字段", "设置上级，选项集中会多pid一项数据，值就是上级字段的值。", "", false, nil, "", nil)
	pageData.FormFieldsAdd("to_tree_array", "radio", "选项集转多维", "将选项集根据pid转为树形结构（多维数组）", "0", false, Models.OptionModels{}.ByKey("is", false), "", map[string]interface{}{
		"if": "formFields.parent_field!=''",
	})

	//!选择下级选项集，要排除自己
	exceptOptionModelKey := ""
	if id > 0 {
		//查找自己的 unique_key
		my, err := db.New().Table("tb_option_models").Where("is_delete", 0).Where("id", id).First()
		if err != nil {
			return err, 500
		}
		if my != nil {
			exceptOptionModelKey = my["unique_key"].(string)
		}
	}
	pageData.FormFieldsAdd("children_option_model_key", "select", "下级选项集", "下级选项集必须查询数据表类型，必须设置上级字段", "", false, that.ChildrenOptionModelsList(exceptOptionModelKey), "", nil)
	//pageData.FormFieldsAdd("options_disable", "radio", "当前选项禁选", "设置此项后，选项只能选择下级的选项集。当前选项集自动添加字段disabled=disabled", "0", false, Models.OptionModels{}.ByKey("is", false), "",
	//	map[string]interface{}{
	//		"if": "!!formFields.children_option_model_key",
	//	})
	//
	//选项图标、颜色、缩进
	pageData.FormFieldsAdd("", "block", "选项图标、颜色、缩进", "", "", false, nil, "", map[string]interface{}{
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
	//自动匹配与排序
	pageData.FormFieldsAdd("", "block", "自动匹配与排序", "", "", false, nil, "", nil)
	//默认排序值
	indexNum := 1
	if id == 0 {
		num, err := db.New().Table("tb_option_models").Where("is_delete", 0).Count()
		if err != nil {
			logger.Error(err.Error())
			return err, 500
		}
		indexNum = util2.Int642Int(num) + 1
	}
	pageData.FormFieldsAdd("match_fields", "textarea", "自动匹配字段", "每个字段占一行，支持全匹配字段和半匹配字段,例如is_*", "", false, nil, "", nil)
	pageData.FormFieldsAdd("index_num", "text-xs", "排序", "", util2.Int2String(indexNum), true, nil, "", nil)

	return nil, 0
}

// NodeFormData 表单显示前修改数据
func (that OptionModels) NodeFormData(pageData *EasyApp.PageData, data gorose.Data, id int64) (gorose.Data, error, int) {
	if id > 0 {

	}
	return data, nil, 0
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
	optionModelKey := util2.PostValue(r, "_dynamic_option_model_key")
	if optionModelKey == "" {
		that.ApiResult(w, 201, "参数不全", nil)
		return
	}
	DynamicParams := Models.OptionModels{}.DynamicParams(optionModelKey)
	var wheres []string
	for _, dp := range DynamicParams {
		v := util2.PostValue(r, dp.ParamKey)
		if v == "" {
			wheres = append(wheres, dp.FieldKey+" = '"+dp.DefValue+"'")
		} else {
			wheres = append(wheres, dp.FieldKey+" = '"+v+"'")
		}
	}
	wheres = append(wheres, "unique_key = '"+optionModelKey+"'")

	ops := Models.OptionModels{}.Select(0, strings.Join(wheres, " and "), false)
	that.ApiResult(w, 200, "success", ops)
}

func (that OptionModels) ExportInsertData(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//查找页面列表
	mainList, err := db.New().Table("tb_option_models").Where("is_delete", 0).Order("index_num,id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Fprint(w, "//! 选项集")
	//遍历列表
	for index, Item := range mainList {
		//删除部分字段
		delete(Item, "id")
		delete(Item, "create_time")
		delete(Item, "update_time")
		delete(Item, "is_delete")
		//按顺序添加排序字段
		Item["index_num"] = index + 1

		content := `
{
	TableName: "tb_option_models",
	Condition: [][]interface{}{{"unique_key", "` + Item["unique_key"].(string) + `"}},
	Data: map[string]interface{}` + util2.JsonEncode(Item) + `,
},
`
		fmt.Fprint(w, content)
	}
}
