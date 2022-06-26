/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 应用配置项管理
 * @File: GroupConfigs
 * @Version: 1.0.0
 * @Date: 2021/10/22 10:57 下午
 */

package adminHandle

import (
	"errors"
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
	"html"
	"strings"
)

type Configs struct {
	Base
}

var ConfigFieldTypes = []map[string]interface{}{
	{
		"value": "text",
		"name":  "文本输入",
	},
	{
		"value": "textarea",
		"name":  "文本域",
	},
	{
		"value": "select",
		"name":  "下拉单选",
	},
	{
		"value": "image",
		"name":  "图片上传",
	},
}

// NodeBegin 开始
func (that Configs) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("应用设置项管理")
	pageData.SetPageName("应用设置项")
	pageData.SetTbName("tb_configs")
	return nil, 0
}

// NodeList 初始化列表
func (that Configs) NodeList(pageData *EasyApp.PageData) (error, int) {
	//获取配置组列表
	ConfigsGroups, err, code := Models.Model{}.SelectOptionsData("tb_configs_group", map[string]string{
		"id":         "value",
		"group_name": "name",
	}, "", "", "", "")
	if err != nil {
		return err, code
	}
	pageData.SetListOrder("group_id asc,index_num asc,id asc")
	pageData.ListColumnClear()
	pageData.ListColumnAdd("group_id", "分组", "array", ConfigsGroups)
	pageData.ListColumnAdd("name", "关键字", "text", nil)
	pageData.ListColumnAdd("title", "配置项名称", "text", nil)
	pageData.ListColumnAdd("value", "当前值/默认值", "html", nil)
	pageData.ListColumnAdd("notice", "说明", "text", nil)
	pageData.ListColumnAdd("index_num", "排序", "text", nil)
	//pageData.ListColumnAdd("status", "状态", "array", models.DefaultStatus)

	pageData.SetListColumnStyle("notice", "width:20%")
	pageData.SetListColumnStyle("value", "width:20%")
	pageData.SetListColumnStyle("action", "width:20%")

	//搜索表单
	pageData.ListSearchFieldAdd("group_id", "select", "选择分组", "", "", ConfigsGroups, "", nil)
	pageData.ListSearchFieldAdd("status", "select", "选择状态", "", "", Models.OptionModels{}.ById(2, false), "", nil)
	return nil, 0
}

// NodeListData 重写列表数据
func (that Configs) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	for key, value := range data {
		if value["field_type"].(string) == "image" || strings.Contains(value["value"].(string), ".png") {
			data[key]["value"] = "<img style=\"width:80px;max-hight:80px\" src=\"" + value["value"].(string) + "\"/>"
		}
		if value["field_type"].(string) == "select" {
			options := []map[string]interface{}{}
			util.JsonDecode(value["options"].(string), &options)
			for _, op := range options {
				if op["value"] == value["value"].(string) {
					data[key]["value"] = op["name"].(string)
					break
				}
			}
		}
	}
	return data, nil, 0
}

// NodeForm 初始化表单
func (that Configs) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	//获取配置组列表
	ConfigsGroups, err, code := Models.Model{}.SelectOptionsData("tb_configs_group", map[string]string{
		"id":         "value",
		"group_name": "name",
	}, "", "", "", "")
	if err != nil {
		return err, code
	}
	jsonDemo := html.EscapeString("[{\"name\":\"\",\"value\":\"\"}]")
	if id == 0 {
		pageData.FormFieldsAdd("group_id", "select", "选择分组", "", "1", true, ConfigsGroups, "", nil)
	}

	pageData.FormFieldsAdd("title", "text", "配置项名称", "中文名称", "", true, nil, "", nil)
	pageData.FormFieldsAdd("name", "text", "关键字", "推荐格式小写字母下划线拼接：app_domain", "", true, nil, "", nil)
	pageData.FormFieldsAdd("value", "text", "默认值", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("notice", "textarea", "配置项说明", "补充说明", "", false, nil, "", nil)
	pageData.FormFieldsAdd("field_type", "select", "表单类型", "", "text", true, ConfigFieldTypes, "", nil)
	pageData.FormFieldsAdd("options", "textarea", "下拉选项", "下拉单选类型的表单需要配置，JSON格式："+jsonDemo, "", false, nil, "", nil)
	pageData.FormFieldsAdd("index_num", "text", "排序", "", "200", true, nil, "", nil)
	pageData.FormFieldsAdd("if", "text", "联动展示", "例如：formFields.xxx>0", "", false, nil, "", nil)
	return nil, 0
}

// NodeFormData 表单显示前修改数据
func (that Configs) NodeFormData(pageData *EasyApp.PageData, data gorose.Data, id int64) (gorose.Data, error, int) {
	if id > 0 {
		//data["options"] = strings.ReplaceAll(data["options"].(string), "\"", "\\\"")
	}
	return data, nil, 0
}

// NodeSaveSuccess 保存成功后操作
func (that Configs) NodeSaveSuccess(pageData *EasyApp.PageData, postData map[string]interface{}, id int64) (bool, error, int) {
	if id > 0 {
		//查询配置项信息
		configInfo, err := db.New().Table("tb_configs").Where("id", id).First()
		if err != nil {
			logger.Error(err.Error())
			return false, errors.New("出错了！"), 500
		}
		//主动更新应用的配置信息
		if configInfo["group_id"].(int64) > 0 {
			cfg, err := db.New().Table("tb_app_configs").
				Where("is_delete", 0).
				Where("group_id", configInfo["group_id"]).
				Where("name", configInfo["name"].(string)).First()
			if err != nil {
				logger.Error(err.Error())
				return false, errors.New("系统运行出错"), 500
			}
			if cfg == nil {
				_, err = db.New().Table("tb_app_configs").Insert(map[string]interface{}{
					"group_id": configInfo["group_id"],
					"name":     configInfo["name"].(string),
					"value":    configInfo["value"],
				})
				if err != nil {
					logger.Error(err.Error())
					return false, errors.New("系统运行出错"), 500
				}
			}
		}
	}
	return true, nil, 0
}
