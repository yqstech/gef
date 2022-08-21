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
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/EasyApp"
	"github.com/yqstech/gef/Utils/db"
	util2 "github.com/yqstech/gef/Utils/util"
	"net/http"
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
	id := util2.GetValue(r, "id")

	configGroup, err := db.New().Table("tb_configs_group").Where("id", id).First()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if configGroup == nil {
		return
	}
	fmt.Fprint(w, "//! 后台设置项分组【"+configGroup["group_name"].(string)+"】")
	delete(configGroup, "id")
	delete(configGroup, "create_time")
	delete(configGroup, "update_time")
	delete(configGroup, "is_delete")
	content := `
{
	TableName: "tb_configs_group",
	Condition: [][]interface{}{{"group_key", "` + configGroup["group_key"].(string) + `"}},
	Data: map[string]interface{}` + util2.JsonEncode(configGroup) + `,
	Children:[]interface{}{
		//!当前分组内的设置项
		
	},
},
`
	fmt.Fprint(w, content)

	fmt.Fprint(w, "\n\n//=============================================>\n")
	fmt.Fprint(w, "//=============================================>\n")
	fmt.Fprint(w, "//=============================================>\n")
	fmt.Fprint(w, "//=============================================>\n\n\n")
	fmt.Fprint(w, "//!下边是当前分组内的设置项\n\n")

	List, err := db.New().
		Table("tb_configs").
		Where("is_delete", 0).
		Where("group_id", id).
		Order("index_num,id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for index, Item := range List {
		//删除部分字段
		delete(Item, "id")
		delete(Item, "create_time")
		delete(Item, "update_time")
		delete(Item, "is_delete")
		//按顺序添加排序字段
		Item["index_num"] = index + 1
		//标记上级ID
		Item["group_id"] = "__PID__"
		content = `
gef.InsideData{
	TableName: "tb_configs",
	Condition: [][]interface{}{{"group_id", "__PID__"},{"name", "` + Item["name"].(string) + `"}},
	Data: map[string]interface{}` + util2.JsonEncode(Item) + `,
},
`
		fmt.Fprint(w, content)
	}

	fmt.Fprint(w, "\n\n//=============================================>\n")
	fmt.Fprint(w, "//=============================================>\n")
	fmt.Fprint(w, "//=============================================>\n")
	fmt.Fprint(w, "//=============================================>\n\n\n")
	fmt.Fprint(w, "//!下边是当前分组内的应用设置记录\n\n")

	List, err = db.New().
		Table("tb_app_configs").
		Where("is_delete", 0).
		Where("group_id", id).
		Order("id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, Item := range List {
		//删除部分字段
		delete(Item, "id")
		delete(Item, "create_time")
		delete(Item, "update_time")
		delete(Item, "is_delete")
		//标记上级ID
		Item["group_id"] = "__PID__"
		content = `
gef.InsideData{
	TableName: "tb_app_configs",
	Condition: [][]interface{}{{"group_id", "__PID__"},{"name", "` + Item["name"].(string) + `"}},
	Data: map[string]interface{}` + util2.JsonEncode(Item) + `,
},
`
		fmt.Fprint(w, content)
	}

	fmt.Fprint(w, "\n\n\n\n\n")
}
