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
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Handles/commHandle"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/builder"
)

type Base struct {
	commHandle.Base
}

// NodeCheckAuth 重写校验权限节点
func (b Base) NodeCheckAuth(_ *builder.PageBuilder, btnRule string, accountID int) (bool, error) {
	return Models.Admin{}.CheckAuth(btnRule, accountID), nil
}

// OptionModelsList 选项集列表
func (b Base) OptionModelsList() []map[string]interface{} {
	//获取列表
	OptionModelsList, err, _ := Models.Model{}.SelectOptionsData("tb_option_models", map[string]string{
		"unique_key": "value",
		"name":       "name",
	}, "", "", "", "")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return OptionModelsList
}

// DynamicOptionModelsList 动态选项集列表
func (b Base) DynamicOptionModelsList() []map[string]interface{} {
	//获取列表
	OptionModelsList, err, _ := Models.Model{}.SelectOptionsData("tb_option_models", map[string]string{
		"unique_key": "value",
		"name":       "name",
	}, "", "", "data_type=1 and dynamic_params!=''", "")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return OptionModelsList
}

// ChildrenOptionModelsList 下级选项集列表
func (b Base) ChildrenOptionModelsList(uniqueKey string) []map[string]interface{} {
	//获取列表
	ChildrenOptionModels, err, _ := Models.Model{}.SelectOptionsData("tb_option_models", map[string]string{
		"unique_key": "value",
		"name":       "name",
	}, "", "", "data_type=1 and parent_field!='' and unique_key!='"+uniqueKey+"'", "")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return ChildrenOptionModels
}
