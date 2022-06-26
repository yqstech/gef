/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 短信通道管理
 * @File: SmsUpstream
 * @Version: 1.0.0
 * @Date: 2021/10/28 10:16 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/config"
	"github.com/gef/rpcPlugins"
)

type SmsUpstream struct {
	Base
}

// NodeBegin 开始
func (that SmsUpstream) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("短信通道管理")
	pageData.SetPageName("短信通道")
	pageData.SetTbName("tb_sms_upstream")
	return nil, 0
}

// NodeList 初始化列表
func (that SmsUpstream) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.SetListOrder("id asc")
	//新增右侧日志开启关闭按钮
	pageData.SetButton("params", EasyApp.Button{
		ButtonName: "配置项",
		Action:     "/sms_upstream_params/index",
		ActionType: 2,
		ActionUrl:  config.AdminPath + "/sms_upstream_params/index",
		Class:      "",
		Icon:       "ri-list-settings-fill",
		Display:    "(!item.btn_params || item.btn_params!='hide')",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	//!重置右侧按钮
	pageData.SetListRightBtns("edit", "disable", "enable", "params", "delete")

	pageData.ListColumnAdd("upstream_name", "通道名称", "text", nil)
	pageData.ListColumnAdd("plugin_name", "调用插件", "text", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	return nil, 0
}

// NodeForm 初始化表单
func (that SmsUpstream) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pluginNames, err := rpcPlugins.PluginsLookup("Sms")
	if err != nil {
		return err, 500
	}
	pluginList := []map[string]interface{}{}
	for _, pluginName := range pluginNames {
		pluginList = append(pluginList, map[string]interface{}{
			"name":  pluginName,
			"value": pluginName,
		})
	}
	pageData.FormFieldsAdd("upstream_name", "text", "通道名称", "", "", true, nil, "", nil)
	pageData.FormFieldsAdd("plugin_name", "select", "选择插件", "", "", true, pluginList, "", nil)
	return nil, 0
}
