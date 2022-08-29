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
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"
	"net/http"
)

type ConfigsGroup struct {
	Base
}

// NodeInit 初始化
func (that *ConfigsGroup) NodeInit(pageBuilder *builder.PageBuilder) {
	//注册handle
	that.NodePageActions["export_insert_data"] = that.ExportInsertData
}

// NodeBegin 开始
func (that ConfigsGroup) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("设置分组管理")
	pageBuilder.SetPageName("设置分组")
	pageBuilder.SetTbName("tb_configs_group")
	return nil, 0
}

// NodeList 初始化列表
func (that ConfigsGroup) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	//导出结构
	pageBuilder.SetButton("export_insert_data", builder.Button{
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
	pageBuilder.SetListRightBtns("edit", "export_insert_data")
	pageBuilder.SetListOrder("id asc")
	pageBuilder.ListColumnAdd("group_key", "分组标识", "text", nil)
	pageBuilder.ListColumnAdd("group_name", "分组名称", "text", nil)
	pageBuilder.ListColumnAdd("note", "分组备注", "text", nil)
	return nil, 0
}

// NodeForm 初始化表单
func (that ConfigsGroup) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
	pageBuilder.FormFieldsAdd("group_key", "text-xs", "分组标识", "配置项分组的唯一标识", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("group_name", "text-xs", "分组名称", "配置项分组的名称", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("note", "text", "分组备注", "", "", true, nil, "", nil)
	return nil, 0
}

// ExportInsertData 导出内置数据
func (that ConfigsGroup) ExportInsertData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := util.GetValue(r, "id")

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
	Data: map[string]interface{}` + util.JsonEncode(configGroup) + `,
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
	Data: map[string]interface{}` + util.JsonEncode(Item) + `,
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
		delete(Item, "is_inside")
		//标记上级ID
		Item["group_id"] = "__PID__"
		content = `
gef.InsideData{
	TableName: "tb_app_configs",
	Condition: [][]interface{}{{"group_id", "__PID__"},{"name", "` + Item["name"].(string) + `"}},
	Data: map[string]interface{}` + util.JsonEncode(Item) + `,
},
`
		fmt.Fprint(w, content)
	}

	fmt.Fprint(w, "\n\n\n\n\n")
}
