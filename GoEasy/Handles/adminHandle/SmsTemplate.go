/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsTemplate
 * @Version: 1.0.0
 * @Date: 2021/10/29 12:59 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
)

type SmsTemplate struct {
	Base
}

// NodeBegin 开始
func (that SmsTemplate) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("短信模板管理")
	pageData.SetPageName("短信模板")
	pageData.SetTbName("tb_sms_template")
	return nil, 0
}

// NodeList 初始化列表
func (that SmsTemplate) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.ListColumnAdd("template_name", "模板标识", "text", nil)
	pageData.ListColumnAdd("template_title", "模板名称", "text", nil)
	pageData.ListColumnAdd("template_vars", "模板变量", "text", nil)
	pageData.ListColumnAdd("default_content", "默认模板内容", "text", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	return nil, 0
}

// NodeForm 初始化表单
func (that SmsTemplate) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("template_name", "text", "模板标识", "和程序保持一致，例如:login_code", "", true, nil, "", nil)
	pageData.FormFieldsAdd("template_title", "text", "模板名称", "模板名称，例如:登录验证码", "", true, nil, "", nil)
	pageData.FormFieldsAdd("template_vars", "text", "模板变量", "程序支持的模板变量，例如:code,other", "", true, nil, "", nil)
	pageData.FormFieldsAdd("default_content", "textarea", "短信模板示例", "变量格式：{{code}}，例如：您正在进行登录验证，验证码为{{code}}，千万不要告诉别人！", "", true, nil, "", nil)
	return nil, 0
}
