/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: ConfigsGroup
 * @Version: 1.0.0
 * @Date: 2022/3/8 5:12 下午
 */

package adminHandle

import (
	"fmt"
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"net/http"
	"strings"
)

type ConfigsGroup struct {
	Base
}
// PageInit 初始化
func (that ConfigsGroup) PageInit(pageData *EasyApp.PageData) {
	//注册handle
	pageData.ActionAdd("export_insert_data", that.ExportInsertData)
}

// NodeBegin 开始
func (that ConfigsGroup) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("设置分组管理")
	pageData.SetPageName("设置分组")
	pageData.SetTbName("tb_configs_group")
	return nil, 0
}

// NodeList 初始化列表
func (that ConfigsGroup) NodeList(pageData *EasyApp.PageData) (error, int) {
	//导出结构
	pageData.SetButton("export_insert_data", EasyApp.Button{
		ButtonName: "导出分组和设置项",
		Action:     "export_insert_data",
		ActionType: 2,
		LayerTitle: "数据导出",
		ActionUrl:  "export_insert_data",
		Class:      "black",
		Icon:       "ri-braces-fill",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	pageData.SetListRightBtns("edit", "export_insert_data")
	pageData.SetListOrder("id asc")
	pageData.ListColumnAdd("group_key", "分组标识", "text", nil)
	pageData.ListColumnAdd("group_name", "分组名称", "text", nil)
	pageData.ListColumnAdd("note", "分组备注", "text", nil)
	return nil, 0
}

// NodeForm 初始化表单
func (that ConfigsGroup) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("group_key", "text-xs", "分组标识", "配置项分组的唯一标识", "", true, nil, "", nil)
	pageData.FormFieldsAdd("group_name", "text-xs", "分组名称", "配置项分组的名称", "", true, nil, "", nil)
	pageData.FormFieldsAdd("note", "text", "分组备注", "", "", true, nil, "", nil)
	return nil, 0
}



// ExportInsertData 导出内置数据
func (that ConfigsGroup) ExportInsertData(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := util.GetValue(r, "id")
	
	var items []string
	
	groups, err := db.New().Table("tb_configs_group").Where("id", id).First()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if groups == nil {
		return
	}
	delete(groups, "id")
	delete(groups, "create_time")
	delete(groups, "update_time")
	delete(groups, "is_delete")
	item := `
		//设置分组` + groups["group_name"].(string) + `和其设置项
		{TableName: "tb_configs_group", Condition: [][]interface{}{{"group_key","` + groups["group_key"].(string) + `"}},
Data: map[string]interface{}` + util.JsonEncode(groups) + `}`
	items = append(items, item)
	
	configList, err := db.New().Table("tb_configs").Where("is_delete", 0).Where("group_id", id).Order("index_num,id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for index, configItem := range configList {
		delete(configItem, "id")
		delete(configItem, "create_time")
		delete(configItem, "update_time")
		delete(configItem, "is_delete")
		configItem["index_num"] = index + 1
		item = `{TableName: "tb_configs", Condition: [][]interface{}{{"group_id", "` + id + `"},{"name","`+configItem["name"].(string)+`"}},
Data: map[string]interface{}` + util.JsonEncode(configItem) + `}`
		items = append(items, item)
	}
	fmt.Fprint(w, strings.Join(items, ",\n"))
}
