/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyModelsButtons
 * @Version: 1.0.0
 * @Date: 2022/5/8 10:35
 */

package adminHandle

import (
	"fmt"
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"net/http"
	"strings"
)

type EasyModelsButtons struct {
	Base
}

var btnActionTypes = []map[string]interface{}{
	{"name": "ajax请求", "value": "1"},
	{"name": "打开弹窗", "value": "2"},
	//{"name": "执行js代码", "value": "3"},
}
var btnClasses = []map[string]interface{}{
	{"name": "原始按钮", "value": "layui-btn-primary"},
	{"name": "默认青绿", "value": "layui-btn"},
	{"name": "主题紫色", "value": "def"},
	{"name": "蓝色信息", "value": "layui-btn-normal"},
	{"name": "黄色警告", "value": "layui-btn-warm"},
	{"name": "红色危险", "value": "layui-btn-danger"},
	{"name": "绿色", "value": "green"},
	{"name": "青色", "value": "cyan"},
	{"name": "棕色", "value": "brown"},
	{"name": "玫瑰", "value": "rose"},
	{"name": "黑色", "value": "black"},
}

// PageInit 初始化
func (that EasyModelsButtons) PageInit(pageData *EasyApp.PageData) {
	//注册handle
	pageData.ActionAdd("export_insert_data", that.ExportInsertData)
}

// NodeBegin 开始
func (that EasyModelsButtons) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("模型按钮管理")
	pageData.SetPageName("按钮")
	pageData.SetTbName("tb_easy_models_buttons")
	return nil, 0
}

// NodeList 初始化列表
func (that EasyModelsButtons) NodeList(pageData *EasyApp.PageData) (error, int) {
	//导出结构
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
	
	pageData.ListColumnAdd("button_key", "按钮关键字", "text", nil)
	pageData.ListColumnAdd("button_name", "按钮名称", "text", nil)
	pageData.ListColumnAdd("button_note", "按钮备注", "text", nil)
	pageData.ListColumnAdd("button_icon", "按钮图标", "icon", nil)
	pageData.ListColumnAdd("class_name", "按钮样式", "array", btnClasses)
	pageData.ListColumnAdd("action_type", "按钮类型", "array", btnActionTypes)
	pageData.ListColumnAdd("status", "按钮状态", "array", Models.OptionModels{}.ByKey("status", true))
	return nil, 0
}

// NodeForm 初始化表单
func (that EasyModelsButtons) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("button_key", "text", "按钮关键字", "系统内的唯一按钮识别标记", "btn_", true, nil, "", nil)
	pageData.FormFieldsAdd("button_note", "text", "按钮备注信息", "按钮的备注信息", "", false, nil, "", nil)
	pageData.FormFieldsAdd("button_name", "text", "按钮名称", "按钮显示的文字", "", true, nil, "", nil)
	pageData.FormFieldsAdd("button_icon", "icon", "按钮图标", "按钮显示的图标 ri-xxx", "ri-radio-button-line", false, nil, "", nil)
	pageData.FormFieldsAdd("class_name", "select", "按钮样式", "按钮样式类", "", false, btnClasses, "", nil)
	pageData.FormFieldsAdd("display", "text", "按钮显示条件", "按钮显示隐藏条件，根据列表项的字段信息自动判断", "", false, nil, "", nil)
	pageData.FormFieldsAdd("action", "text", "链接权限", "用来校验权限，不要以/开头，例如： add 或 user/add", "", true, nil, "", nil)
	pageData.FormFieldsAdd("action_url", "text", "链接地址", "相对路径或绝对路径,/开头会自动补全后台地址", "", true, nil, "", nil)
	pageData.FormFieldsAdd("action_type", "radio", "按钮页面类型", "", "1", true, btnActionTypes, "", nil)
	pageData.FormFieldsAdd("confirm_msg", "text", "确认信息", "填入值，则点击按钮，弹出确认对话框", "", false, nil, "", map[string]interface{}{
		"if": "formFields.action_type==1",
	})
	pageData.FormFieldsAdd("batch_action", "radio", "支持批量操作", "", "0", true, Models.OptionModels{}.ByKey("is", false), "", map[string]interface{}{
		"if": "formFields.action_type==1",
	})
	pageData.FormFieldsAdd("layer_title", "text", "弹窗标题", "弹窗类按钮，在此定义弹窗标题", "", false, nil, "", map[string]interface{}{
		"if": "formFields.action_type==2",
	})
	pageData.FormFieldsAdd("layer_width", "text", "弹窗宽度", "设置弹窗尺寸，支持px和%", "90%", false, nil, "", map[string]interface{}{
		"if": "formFields.action_type==2",
	})
	pageData.FormFieldsAdd("layer_height", "text", "弹窗高度", "设置弹窗尺寸，支持px和%", "86%", false, nil, "", map[string]interface{}{
		"if": "formFields.action_type==2",
	})
	
	return nil, 0
}

func (that EasyModelsButtons) ExportInsertData(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var items []string
	buttons, err := db.New().Table("tb_easy_models_buttons").Where("is_delete", 0).Order("id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	
	item := ""
	for index, button := range buttons {
		delete(button, "id")
		delete(button, "create_time")
		delete(button, "update_time")
		delete(button, "is_delete")
		button["index_num"] = index + 1
		item = `{TableName: "tb_easy_models_buttons", Condition: [][]interface{}{{"button_key", "` + button["button_key"].(string) + `"}},
Data: map[string]interface{}` + util.JsonEncode(button) + `}`
		items = append(items, item)
	}
	fmt.Fprint(w, "\r\n//后台模型自定义按钮\r\n"+strings.Join(items, ",\n"))
}
