/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-10-21 14:09:42
 * @LastEditTime: 2021-10-21 14:37:20
 * @Description: 后台日志
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/util"
)

type AdminLog struct {
	Base
}

func (a AdminLog) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("后台操作日志")
	pageData.SetPageName("操作认知")
	pageData.SetTbName("tb_admin_log")
	return nil, 0
}

func (a AdminLog) NodeList(pageData *EasyApp.PageData) (error, int) {
	//清除顶部和右侧按钮
	pageData.ListRightBtnsClear()
	pageData.ListTopBtnsClear()

	adminList, err, code := Models.Model{}.SelectOptionsData("tb_admin", map[string]string{
		"id":   "value",
		"name": "name",
	}, "", "", "", "")
	if err != nil {
		return err, code
	}
	pageData.ListSearchFieldAdd("account_id", "select", "操作人", "", "", adminList, "", nil)

	pageData.ListColumnAdd("account_id", "操作人ID", "text", nil)
	pageData.ListColumnAdd("account_name", "操作人名称", "text", nil)
	pageData.ListColumnAdd("account", "操作人账号", "text", nil)
	pageData.ListColumnAdd("url", "链接地址", "text", nil)
	pageData.ListColumnAdd("rule_name", "执行操作", "text", nil)
	pageData.ListColumnAdd("notice", "备注", "text", nil)
	return nil, 0
}

func (a AdminLog) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {

	for k, v := range condition {
		if v[0].(string) == "account_id" {
			if util.Interface2String(v[1]) == "" {
				condition = append(condition[0:k], condition[k+1:]...)
			}
		}
	}
	return condition, nil, 0
}
