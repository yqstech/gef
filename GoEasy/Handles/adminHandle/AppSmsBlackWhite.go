/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AppSmsBlackWhite
 * @Version: 1.0.0
 * @Date: 2022/3/8 9:42 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
)

type AppSmsBlackWhite struct {
	Base
}

// NodeBegin 开始
func (that AppSmsBlackWhite) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("短信黑名单和白名单")
	pageData.SetPageName("黑白名单")
	pageData.SetTbName("tb_app_sms_black_white")
	return nil, 0
}

// NodeList 初始化列表
func (that AppSmsBlackWhite) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.ListColumnClear()
	pageData.ListColumnAdd("type", "类型", "array", Models.SmsBlackWhiteTypes)
	pageData.ListColumnAdd("rule_type", "匹配类型", "array", Models.SmsBlackWhiteRuleTypes)
	pageData.ListColumnAdd("rule", "手机号或ip地址", "text", nil)
	pageData.ListColumnAdd("note", "备注", "text", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))

	return nil, 0
}

// NodeForm 初始化表单
func (that AppSmsBlackWhite) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("type", "radio", "类型", "", "0", true, Models.SmsBlackWhiteTypes, "", nil)
	pageData.FormFieldsAdd("rule_type", "radio", "匹配类型", "", "0", true, Models.SmsBlackWhiteRuleTypes, "", nil)
	pageData.FormFieldsAdd("rule", "text", "匹配内容", "手机号或ip地址", "", true, nil, "", nil)
	pageData.FormFieldsAdd("note", "text", "备注", "", "", false, nil, "", nil)
	return nil, 0
}
