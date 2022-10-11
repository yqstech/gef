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
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"
	"net/http"
	"strings"
)

type OptionModels struct {
	Base
}

var dataTypes = []map[string]interface{}{
	{"name": "静态", "value": "0", "color": "#0066CC", "icon": "ri-braces-line"},
	{"name": "查询", "value": "1", "color": "#FF6600", "icon": "ri-database-2-line"},
}

// NodeInit 初始化
func (that *OptionModels) NodeInit(pageBuilder *builder.PageBuilder) {
	//注册handle
	that.NodePageActions["dynamic"] = that.Dynamic
	that.NodePageActions["export_insert_data"] = that.ExportInsertData
}

// NodeBegin 开始
func (that OptionModels) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("选项集管理")
	pageBuilder.SetPageName("选项集")
	pageBuilder.SetTbName("tb_option_models")
	return nil, 0
}

// NodeList 初始化列表
func (that OptionModels) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetButton("export_insert_data", builder.Button{
		ButtonName: "导出数据",
		Action:     "export_insert_data",
		ActionType: 2,
		LayerTitle: "选项集导出成内置数据",
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
	pageBuilder.SetListTopBtns("add", "export_insert_data")
	pageBuilder.SetListPageSize(50)
	pageBuilder.SetListOrder("index_num,id asc")
	pageBuilder.ListColumnAdd("unique_key", "标识符", "html", nil)
	pageBuilder.ListColumnAdd("name", "名称", "text", nil)
	pageBuilder.ListColumnAdd("data_type", "数据类型", "array", Models.OptionModels{}.Beautify(dataTypes))
	pageBuilder.ListColumnAdd("static_data", "静态数据/数据表字段", "html", nil)
	//pageBuilder.ListColumnAdd("table_name", "数据表", "text", nil)
	//pageBuilder.ListColumnAdd("value_field", "值字段", "text", nil)
	//pageBuilder.ListColumnAdd("name_field", "名称字段", "text", nil)
	//pageBuilder.ListColumnAdd("parent_field", "上级字段", "text", nil)
	pageBuilder.ListColumnAdd("children_option_model_key", "下级选项集", "array", that.ChildrenOptionModelsList(""))
	pageBuilder.ListColumnAdd("index_num", "排序", "input::type=number&width=50px", nil)
	pageBuilder.SetListColumnStyle("unique_key", "width:10%")
	pageBuilder.SetListColumnStyle("name", "width:10%")
	return nil, 0
}

// NodeListData 重写列表数据
func (that OptionModels) NodeListData(pageBuilder *builder.PageBuilder, data []gorose.Data) ([]gorose.Data, error, int) {
	for k, v := range data {

		if v["data_type"].(int64) == 1 {

			data[k]["unique_key"] = v["unique_key"].(string) + "<br><span style=\"color:#FF6600\">" + v["table_name"].(string) + "</span>"

			data[k]["static_data"] = ""
			if v["value_field"].(string) != "" || v["name_field"].(string) != "" {
				data[k]["static_data"] = data[k]["static_data"].(string) + "<span style=\"color:#FF6600\">值字段：" + v["value_field"].(string) + "</span><br>" + "<span style=\"color:#FF6600\">名字段：" + v["name_field"].(string) + "</span><br>"
			}
			if v["parent_field"].(string) != "" {
				data[k]["static_data"] = data[k]["static_data"].(string) + "<span style=\"color:#FF6600\">上级字段：" + v["parent_field"].(string) + "</span>"
			}
		} else {
			if v["static_data"].(string) != "" {
				var staticData []map[string]interface{}
				util.JsonDecode(v["static_data"].(string), &staticData)
				var transData []string
				for _, opt := range staticData {
					transData = append(transData, opt["name"].(string)+" : "+util.Interface2String(opt["value"]))
				}
				data[k]["static_data"] = strings.Join(transData, "<br>")
			}
		}
	}
	return data, nil, 0
}

// NodeForm 初始化表单
func (that OptionModels) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
	pageBuilder.FormFieldsAdd("unique_key", "text-sm", "标识符", "选项集唯一的识别关键字", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("name", "text-sm", "名称", "", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("data_type", "radio", "数据类型", "", "0", true, dataTypes, "", nil)
	pageBuilder.FormFieldsAdd("static_data", "textarea", "静态数据", "", "[{\"name\":\"是\",\"value\":\"1\"},{\"name\":\"否\",\"value\":\"0\"}]", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==0",
	})
	pageBuilder.FormFieldsAdd("default_data", "textarea", "默认数据", "和静态数据格式相同", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	//配置数据表
	pageBuilder.FormFieldsAdd("", "block", "配置数据表", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("table_name", "text-xs", "数据表名称", "tb_", "tb_", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("value_field", "text-xs", "值字段", "查询到的数据作为值", "id", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("name_field", "text-xs", "名称字段", "查询到的数据作为名称", "name", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("select_order", "text-xs", "查询排序", "数据表查询的排序方式", "id asc", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("select_where", "text-sm", "补充查询条件", "补充数据表查询条件", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})

	//联动设置
	pageBuilder.FormFieldsAdd("", "block", "选项集联动", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})

	pageBuilder.FormFieldsAdd("dynamic_params", "textarea", "联动配置", "用来做数据联动的参数设置，程序根据设置的字段，查询post参数\n格式为 监听参数:选项集数据表字段:默认值，例如：group_id:group_id:0\n默认值为空自动忽略，每行一个转换规则", "", false,
		nil, "", map[string]interface{}{
			"if": "formFields.data_type==1",
		})

	//多级、数据转换
	pageBuilder.FormFieldsAdd("", "block", "多级、数据转换", "", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("parent_field", "text-xs", "上级字段", "设置上级，选项集中会多pid一项数据，值就是上级字段的值。", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("to_tree_array", "radio", "选项集转多维", "将选项集根据pid转为树形结构（多维数组）", "0", false, Models.OptionModels{}.ByKey("is", false), "", map[string]interface{}{
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
	pageBuilder.FormFieldsAdd("children_option_model_key", "select", "下级选项集", "下级选项集必须查询数据表类型，必须设置上级字段", "", false, that.ChildrenOptionModelsList(exceptOptionModelKey), "", nil)
	//pageBuilder.FormFieldsAdd("options_disable", "radio", "当前选项禁选", "设置此项后，选项只能选择下级的选项集。当前选项集自动添加字段disabled=disabled", "0", false, Models.OptionModels{}.ByKey("is", false), "",
	//	map[string]interface{}{
	//		"if": "!!formFields.children_option_model_key",
	//	})
	//
	//选项图标、颜色、缩进
	pageBuilder.FormFieldsAdd("", "block", "选项图标、颜色、缩进", "", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("color_field", "text", "颜色字段", "设置某字段值作为颜色值", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("icon_field", "text", "Icon字段", "设置某字段值作为Icon图标值", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("color_array", "textarea", "颜色组", "选项列表对应的颜色列表，使用英文逗号分割，例如:#FFFFFF,#FFF660", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	pageBuilder.FormFieldsAdd("icon_array", "textarea", "Icon图标组", "选项列表对应的图标列表，使用英文逗号分割，例如:ri-send-plane-line,ri-chat-3-line", "", false, nil, "", map[string]interface{}{
		"if": "formFields.data_type==1",
	})
	//自动匹配与排序
	pageBuilder.FormFieldsAdd("", "block", "自动匹配与排序", "", "", false, nil, "", nil)
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
	pageBuilder.FormFieldsAdd("match_fields", "textarea", "自动匹配字段", "每个字段占一行，支持全匹配字段和半匹配字段,例如is_*", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("index_num", "text-xs", "排序", "", util.Int2String(indexNum), true, nil, "", nil)

	return nil, 0
}

// NodeFormData 表单显示前修改数据
func (that OptionModels) NodeFormData(pageBuilder *builder.PageBuilder, data gorose.Data, id int64) (gorose.Data, error, int) {
	if id > 0 {

	}
	return data, nil, 0
}

// NodeSaveData 表单保存数据前使用
func (that OptionModels) NodeSaveData(pageBuilder *builder.PageBuilder, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
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
func (that OptionModels) Dynamic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	optionModelKey := util.PostValue(r, "_dynamic_option_model_key")
	if optionModelKey == "" {
		that.ApiResult(w, 201, "参数不全", nil)
		return
	}
	DynamicParams := Models.OptionModels{}.DynamicParams(optionModelKey)
	var wheres []string
	for _, dp := range DynamicParams {
		v := util.PostValue(r, dp.ParamKey)
		if v == "" {
			wheres = append(wheres, dp.FieldKey+" = '"+dp.DefValue+"'")
		} else {
			wheres = append(wheres, dp.FieldKey+" = '"+v+"'")
		}
	}
	ops := Models.OptionModels{}.Select(0, "unique_key = '"+optionModelKey+"'", strings.Join(wheres, " and "), false)
	//选项缩进
	var FieldOptions []map[string]interface{}
	//深度拷贝，直接复制的话，多次刷新会造成多次缩进
	for _, v := range ops {
		x := map[string]interface{}{}
		for k1, v1 := range v {
			x[k1] = v1
		}
		FieldOptions = append(FieldOptions, x)
	}
	//转一下类型
	var FieldOptionsCopy []gorose.Data
	for _, v := range FieldOptions {
		v["id"] = int64(util.String2Int(util.Interface2String(v["id"])))
		v["pid"] = int64(util.String2Int(util.Interface2String(v["pid"])))
		FieldOptionsCopy = append(FieldOptionsCopy, v)
	}
	//执行缩进
	FieldOptionsCopy = Models.Model{}.GoroseDataLevelOrder(FieldOptionsCopy, "id", "pid", 0, 0)
	for k, v := range FieldOptionsCopy {
		if v["level"] == int64(0) {
		} else if v["level"] == int64(1) {
			FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		} else if v["level"] == int64(2) {
			FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		} else if v["level"] == int64(3) {
			FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		} else if v["level"] == int64(4) {
			FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		}
	}
	//类型再回来
	FieldOptions = []map[string]interface{}{}
	for _, v := range FieldOptionsCopy {
		FieldOptions = append(FieldOptions, v)
	}
	that.ApiResult(w, 200, "success", FieldOptions)
}

func (that OptionModels) ExportInsertData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	Data: map[string]interface{}` + util.JsonEncode(Item) + `,
},
`
		fmt.Fprint(w, content)
	}
}
