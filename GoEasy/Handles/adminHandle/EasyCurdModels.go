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
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gef/config"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"net/http"
	"strings"
)

type EasyCurdModels struct {
	Base
}
// PageInit 初始化
func (that EasyCurdModels) PageInit(pageData *EasyApp.PageData) {
	//注册handle
	pageData.ActionAdd("export_insert_data", that.ExportInsertData)
}

// NodeBegin 开始
func (that EasyCurdModels) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("接口模型 EasyCurdModel")
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
	//导出结构
	pageData.SetButton("export_insert_data", EasyApp.Button{
		ButtonName: "",
		Action:     "/easy_curd_models/export_insert_data",
		ActionType: 2,
		LayerTitle: "后台模型导出成内置数据",
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
	pageData.SetListRightBtns("edit", "fields", "export_insert_data", "disable", "enable", "delete")
	pageData.SetListOrder("id asc")
	pageData.ListColumnAdd("model_key", "模型Key", "text", nil)
	pageData.ListColumnAdd("model_name", "模型名称", "text", nil)
	//pageData.ListColumnAdd("table_name", "数据表名", "text", nil)
	pageData.ListColumnAdd("allow_select", "查询", "switch::text=允许|禁止", nil)
	pageData.ListColumnAdd("allow_create", "新增", "switch::text=允许|禁止", nil)
	pageData.ListColumnAdd("allow_update", "更新", "switch::text=允许|禁止", nil)
	pageData.ListColumnAdd("allow_delete", "删除", "switch::text=允许|禁止", nil)
	//pageData.ListColumnAdd("soft_delete_disable", "删除方式", "switch::text=硬删除|软删除", nil)
	//pageData.ListColumnAdd("check_login", "校验登录", "switch::text=是|否", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ByKey("status", true))
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
	pageData.FormFieldsAdd("select_with_disabled", "radio", "禁用/取消可查", "status=0的数据是否可以查询", "0", true, Models.DefaultIsOrNot, "", nil)
	pageData.FormFieldsAdd("uk_name", "text", "用户字段名", "代表用户id的字段名称", "user_id", true, nil, "", nil)
	pageData.FormFieldsAdd("pk_name", "text", "主键字段名", "代表用主键的字段名称", "id", true, nil, "", nil)
	return nil, 0
}


// ExportInsertData 导出内置数据
func (that EasyCurdModels) ExportInsertData(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := util.GetValue(r, "id")
	
	s1 := "var InsideData = []gef.InsideData{"
	s2 := ",\r\n}\r\n"
	var items []string
	
	easyModel, err := db.New().Table("tb_easy_curd_models").Where("id", id).First()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if easyModel == nil {
		return
	}
	delete(easyModel, "create_time")
	delete(easyModel, "update_time")
	delete(easyModel, "is_delete")
	item := `
		//接口模型及模型字段-` + easyModel["model_name"].(string) + `
		{TableName: "tb_easy_curd_models", Condition: [][]interface{}{{"model_key","` + easyModel["model_key"].(string) + `"}},Data: map[string]interface{}` + util.JsonEncode(easyModel) + `}`
	items = append(items, item)
	
	Fields, err := db.New().Table("tb_easy_curd_models_fields").Where("is_delete", 0).Where("model_id", id).Order("id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, field := range Fields {
		delete(field, "id")
		delete(field, "create_time")
		delete(field, "update_time")
		delete(field, "is_delete")
		item = `{TableName: "tb_easy_curd_models_fields", Condition: [][]interface{}{{"model_id", "` + id + `"},{"field_key", "` + field["field_key"].(string) + `"}},
Data: map[string]interface{}` + util.JsonEncode(field) + `}`
		items = append(items, item)
	}
	fmt.Fprint(w, s1+strings.Join(items, ",\n")+s2)
}
