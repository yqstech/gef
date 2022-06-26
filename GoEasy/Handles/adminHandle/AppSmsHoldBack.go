/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AppSmsHoldBack
 * @Version: 1.0.0
 * @Date: 2022/3/8 9:43 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gohouse/gorose/v2"
)

type AppSmsHoldBack struct {
	Base
}

// NodeBegin 开始
func (that AppSmsHoldBack) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("短信防火墙规则")
	pageData.SetPageName("规则")
	pageData.SetTbName("tb_app_sms_hold_back")
	return nil, 0
}

// NodeList 初始化列表
func (that AppSmsHoldBack) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.ListColumnClear()
	pageData.ListColumnAdd("rule_type", "规则类型", "array", Models.SmsHoldBackRuleTypes)
	pageData.ListColumnAdd("range_second", "几秒钟内", "text", nil)
	pageData.ListColumnAdd("sms_max", "短信超过几条", "text", nil)
	pageData.ListColumnAdd("action", "执行操作", "array", Models.SmsHoldBackActions)
	pageData.ListColumnAdd("frozen_second", "暂停秒数", "text", nil)
	pageData.ListColumnAdd("note", "备注", "text", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	return nil, 0
}

// NodeListData 重写列表数据
func (that AppSmsHoldBack) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	for k, v := range data {
		if v["frozen_second"].(int64) < 0 {
			data[k]["frozen_second"] = "#"
		}
	}
	return data, nil, 0
}

// NodeForm 初始化表单
func (that AppSmsHoldBack) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("rule_type", "radio", "规则类型", "", "0", true, Models.SmsHoldBackRuleTypes, "", nil)
	pageData.FormFieldsAdd("range_second", "text", "几秒钟内", "填写大于0的整数", "", false, nil, "", nil)
	pageData.FormFieldsAdd("sms_max", "text", "短信超过几条", "填写大于0的整数", "", false, nil, "", nil)
	pageData.FormFieldsAdd("action", "radio", "执行操作", "", "1", false, Models.SmsHoldBackActions, "", map[string]interface{}{
		"if": "formFields.rule_type<2",
	})
	pageData.FormFieldsAdd("action2", "radio", "执行操作", "", "1", false, []map[string]interface{}{{"name": "暂停短信功能(白名单除外)", "value": "1"}}, "", map[string]interface{}{
		"if": "formFields.rule_type==2",
	})
	pageData.FormFieldsAdd("frozen_second", "text", "暂停秒数", "填写大于0的整数", "", false, nil, "", map[string]interface{}{
		"if": "!(formFields.rule_type<2 && formFields.action==2)",
	})
	pageData.FormFieldsAdd("note", "text", "备注", "", "", false, nil, "", nil)
	return nil, 0
}

// NodeSaveData 表单保存数据前使用
func (that AppSmsHoldBack) NodeSaveData(pageData *EasyApp.PageData, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	if postData["rule_type"] == "2" {
		postData["action"] = postData["action2"]
	} else {
		if postData["action"] == "2" {
			postData["frozen_second"] = -1
		}
	}
	delete(postData, "action2")
	return postData, nil, 0
}
