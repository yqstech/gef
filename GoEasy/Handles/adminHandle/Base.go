/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Base.go
 * @Version: 1.0.0
 * @Date: 2021/10/21 6:05 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Handles/commHandle"
	"github.com/gef/GoEasy/Models"
	"github.com/wonderivan/logger"
)

type Base struct {
	commHandle.Base
}

// NodeCheckAuth 重写校验权限节点
func (b Base) NodeCheckAuth(pageData *EasyApp.PageData, btnRule string, accountID int) (bool, error) {
	return Models.Admin{}.CheckAuth(btnRule, accountID), nil
}

func (b Base) SmsUpstreamList() []map[string]interface{} {
	//获取列表
	//获取列表
	upstreamOptions, err, _ := Models.Model{}.SelectOptionsData("tb_sms_upstream", map[string]string{
		"id":            "value",
		"upstream_name": "name",
	}, "", "", "", "")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return upstreamOptions
}
func (b Base) EasyModels() []map[string]interface{} {
	data, err, _ := Models.Model{}.SelectOptionsData("tb_easy_models", map[string]string{
		"id":         "value",
		"model_name": "name",
	}, "", "", "", "")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return data
}

func (b Base) OptionModelsList() []map[string]interface{} {
	//获取列表
	SelectData, err, _ := Models.Model{}.SelectOptionsData("tb_option_models", map[string]string{
		"id":   "value",
		"name": "name",
	}, "", "", "", "")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return SelectData
}

func (b Base) DynamicOptionModelsList() []map[string]interface{} {
	//获取列表
	SelectData, err, _ := Models.Model{}.SelectOptionsData("tb_option_models", map[string]string{
		"id":   "value",
		"name": "name",
	}, "", "", "data_type=1 and dynamic_params!=''", "")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return SelectData
}
