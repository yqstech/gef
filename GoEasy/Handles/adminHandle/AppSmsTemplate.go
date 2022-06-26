/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AppSmsTemplate
 * @Version: 1.0.0
 * @Date: 2021/10/29 10:05 下午
 */

package adminHandle

import (
	"errors"
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
	"strings"
)

type AppSmsTemplate struct {
	Base
}

// NodeBegin 开始
func (that AppSmsTemplate) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	//同步模板信息
	templateList, err := db.New().Table("tb_sms_template").
		Where("is_delete", 0).
		Where("status", "1").Get()
	if err != nil {
		logger.Error(err.Error())
		return errors.New("查询信息出错！"), 0
	}
	for _, templateInfo := range templateList {
		Managertemplate, err := db.New().Table("tb_app_sms_template").
			Where("template_name", templateInfo["template_name"]).
			Where("is_delete", 0).
			First()
		if err != nil {
			logger.Error(err.Error())
			return errors.New("查询信息出错！"), 0
		}
		if Managertemplate == nil {
			db.New().Table("tb_app_sms_template").Insert(map[string]interface{}{
				"template_name":    templateInfo["template_name"],
				"template_content": templateInfo["default_content"],
				"create_time":      util.TimeNow(),
				"update_time":      util.TimeNow(),
				"status":           0,
			})
		}
	}
	pageData.SetTitle("短信模板管理")
	pageData.SetPageName("短信模板")
	pageData.SetTbName("tb_app_sms_template")
	return nil, 0
}

// NodeList 初始化列表
func (that AppSmsTemplate) NodeList(pageData *EasyApp.PageData) (error, int) {
	//清除列表顶部和右侧按钮
	pageData.ListRightBtnsClear()
	pageData.ListTopBtnsClear()
	//重置右侧按钮
	pageData.SetListRightBtns("edit", "disable", "enable")
	//获取列表
	templateOptions, err, code := Models.Model{}.SelectOptionsData("tb_sms_template", map[string]string{
		"template_name":  "value",
		"template_title": "name",
	}, "", "", "", "")
	if err != nil {
		return err, code
	}
	pageData.ListColumnAdd("template_name", "模板标识", "text", nil)
	pageData.ListColumnAdd("template_name", "模板名称", "array", templateOptions)
	pageData.ListColumnAdd("template_out_id", "外部模板ID", "text", nil)
	pageData.ListColumnAdd("template_content", "模板内容", "text", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	return nil, 0
}

// NodeForm 初始化表单
func (that AppSmsTemplate) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	//查询通道
	if id <= 0 {
		return errors.New("获取通道ID失败！"), 0
	}
	//查询应用短信模板
	managertemplate, err := db.New().Table("tb_app_sms_template").
		Where("id", id).
		Where("is_delete", 0).
		First()
	if err != nil {
		return err, 0
	}
	if managertemplate == nil {
		return errors.New("获取通道信息失败！"), 0
	}
	// 查询短信模板
	templateInfo, err := db.New().Table("tb_sms_template").
		Where("template_name", managertemplate["template_name"]).
		Where("is_delete", 0).
		Where("status", 1).First()
	if err != nil {
		return err, 0
	}
	vars := strings.Split(templateInfo["template_vars"].(string), ",")
	contentNotices := "当前模板支持的变量有："
	for _, v := range vars {
		contentNotices = contentNotices + "{{" + v + "}}"
	}

	pageData.FormFieldsAdd("template_out_id", "text", "外部模板ID", "例如阿里云模板ID SMS_123456", "", false, nil, "", nil)
	pageData.FormFieldsAdd("template_content", "textarea", "模板内容", "变量格式为变量+双花括号，例如验证码：{{code}}", "", false, nil, "", nil)
	pageData.FormFieldsAdd("", "notice", "", contentNotices, "", false, nil, "", nil)
	return nil, 0
}
