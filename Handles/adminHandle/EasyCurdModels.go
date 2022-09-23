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

type EasyCurdModels struct {
	Base
}

// NodeInit 初始化
func (that *EasyCurdModels) NodeInit(pageBuilder *builder.PageBuilder) {
	//注册handle
	that.NodePageActions["export_insert_data"] = that.ExportInsertData
	that.NodePageActions["update_all_models_fields"] = that.UpdateAllFields
	that.NodePageActions["export_all_data"] = that.ExportAllData
}

// NodeBegin 开始
func (that EasyCurdModels) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("接口模型 EasyCurdModel")
	pageBuilder.SetPageName("接口模型")
	pageBuilder.SetTbName("tb_easy_curd_models")
	return nil, 0
}

// NodeList 初始化列表
func (that EasyCurdModels) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	//新增顶部按钮
	pageBuilder.SetButton("update_all_models_fields", builder.Button{
		ButtonName: "刷新字段",
		Action:     "/easy_curd_models/update_all_models_fields",
		ActionType: 2,
		LayerTitle: "刷新全部模型字段",
		ActionUrl:  config.AdminPath + "/easy_curd_models/update_all_models_fields",
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
		Action:     "/easy_curd_models/export_all_data",
		ActionType: 2,
		LayerTitle: "导出全部接口模型",
		ActionUrl:  config.AdminPath + "/easy_curd_models/export_all_data",
		Class:      "black",
		Icon:       "ri-braces-fill",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//!重置顶部按钮
	pageBuilder.SetListTopBtns("add", "update_all_models_fields", "export_all_data")

	//新增右侧字段管理按钮
	pageBuilder.SetButton("fields", builder.Button{
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
	//导出结构
	pageBuilder.SetButton("export_insert_data", builder.Button{
		ButtonName: "",
		Action:     "/easy_curd_models/export_insert_data",
		ActionType: 2,
		LayerTitle: "接口模型导出成内置数据",
		ActionUrl:  config.AdminPath + "/easy_curd_models/export_insert_data",
		Class:      "black",
		Icon:       "ri-braces-fill",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//!重置右侧按钮
	pageBuilder.SetListRightBtns("edit", "fields", "export_insert_data", "disable", "enable", "delete")
	pageBuilder.SetListOrder("id asc")
	pageBuilder.ListColumnAdd("model_key", "模型Key", "text", nil)
	pageBuilder.ListColumnAdd("model_name", "模型名称", "text", nil)
	//pageBuilder.ListColumnAdd("table_name", "数据表名", "text", nil)
	pageBuilder.ListColumnAdd("allow_select", "查询", "switch::text=允许|禁止", nil)
	pageBuilder.ListColumnAdd("allow_create", "新增", "switch::text=允许|禁止", nil)
	pageBuilder.ListColumnAdd("allow_update", "更新", "switch::text=允许|禁止", nil)
	pageBuilder.ListColumnAdd("allow_delete", "删除", "switch::text=允许|禁止", nil)
	//pageBuilder.ListColumnAdd("soft_delete_disable", "删除方式", "switch::text=硬删除|软删除", nil)
	//pageBuilder.ListColumnAdd("check_login", "校验登录", "switch::text=是|否", nil)
	pageBuilder.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ByKey("status", true))
	return nil, 0
}

// NodeForm 初始化表单
func (that EasyCurdModels) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
	pageBuilder.FormFieldsAdd("model_key", "text", "模型Key", "模型唯一识别码", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("model_name", "text", "模型名称", "模型自定义名称", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("table_name", "text", "关联数据表", "关联的数据表完整名称", "tb_", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("select_order", "text", "查询排序", "查询排序方式", "id desc", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("allow_select", "radio", "允许查询", "", "1", true, Models.DefaultIsOrNot, "", nil)
	pageBuilder.FormFieldsAdd("allow_create", "radio", "允许新增", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageBuilder.FormFieldsAdd("allow_update", "radio", "允许更新", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageBuilder.FormFieldsAdd("allow_delete", "radio", "允许删除", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageBuilder.FormFieldsAdd("soft_delete_disable", "radio", "禁用软删除", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageBuilder.FormFieldsAdd("check_login", "radio", "校验登录", "", "0", true, Models.DefaultIsOrNot, "", nil)
	pageBuilder.FormFieldsAdd("select_with_disabled", "radio", "禁用/取消可查", "status=0的数据是否可以查询", "0", true, Models.DefaultIsOrNot, "", nil)
	pageBuilder.FormFieldsAdd("uk_name", "text", "用户字段名", "代表用户id的字段名称", "user_id", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("pk_name", "text", "主键字段名", "代表用主键的字段名称", "id", true, nil, "", nil)
	return nil, 0
}

// UpdateAllFields 更新所有字段
func (that EasyCurdModels) UpdateAllFields(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	models, err := db.New().Table("tb_easy_curd_models").Where("is_delete", 0).Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ModelsFields := EasyCurdModelsFields{}
	for _, m := range models {
		ModelsFields.syncModelFields(util.Int642Int(m["id"].(int64)))
		fmt.Fprint(w, "刷新模型【"+m["model_name"].(string)+"】成功！\n")
	}
	fmt.Fprint(w, "操作成功！")
}

// ExportAllData 导出所有数据
func (that EasyCurdModels) ExportAllData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	models, err := db.New().
		Table("tb_easy_curd_models").
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
func (that EasyCurdModels) ExportInsertData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := util.GetValue(r, "id")

	easyModel, err := db.New().Table("tb_easy_curd_models").Where("id", id).First()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if easyModel == nil {
		return
	}
	that.ExportModelsAndFileds(w, easyModel)
}

func (that EasyCurdModels) ExportModelsAndFileds(w http.ResponseWriter, easyModel gorose.Data) {
	fmt.Fprint(w, "//! 接口模型【"+easyModel["model_name"].(string)+"】")
	fieldsContent := that.ExportModelFileds(easyModel["id"].(int64))
	delete(easyModel, "id")
	delete(easyModel, "create_time")
	delete(easyModel, "update_time")
	delete(easyModel, "is_delete")

	content := `
{
	TableName: "tb_easy_curd_models",
	Condition: [][]interface{}{{"model_key", "` + easyModel["model_key"].(string) + `"}},
	Data: map[string]interface{}` + util.JsonEncode(easyModel) + `,
	Children:[]interface{}{
		//!接口模型的字段
		` + fieldsContent + `
	},
},
`
	fmt.Fprint(w, content)

}

// ExportModelFileds 导出字段
func (that EasyCurdModels) ExportModelFileds(id int64) string {
	Fields, err := db.New().
		Table("tb_easy_curd_models_fields").
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
			TableName: "tb_easy_curd_models_fields",
			Condition: [][]interface{}{{"model_id", "__PID__"},{"field_key", "` + Item["field_key"].(string) + `"}},
			Data: map[string]interface{}` + util.JsonEncode(Item) + `,
		},
`
	}
	return content
}
