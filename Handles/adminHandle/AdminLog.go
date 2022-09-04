/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-10-21 14:09:42
 * @LastEditTime: 2021-10-21 14:37:20
 * @Description: 后台日志
 */

package adminHandle

import (
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"
)

type AdminLog struct {
	Base
}

func (a AdminLog) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("后台操作日志")
	pageBuilder.SetPageName("操作认知")
	pageBuilder.SetTbName("tb_admin_log")
	return nil, 0
}

func (a AdminLog) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	//清除顶部和右侧按钮
	pageBuilder.ListRightBtnsClear()
	pageBuilder.ListTopBtnsClear()

	adminList, err, code := Models.Model{}.SelectOptionsData("tb_admin", map[string]string{
		"id":   "value",
		"name": "name",
	}, "", "", "", "")
	if err != nil {
		return err, code
	}
	pageBuilder.ListSearchFieldAdd("account_id", "select", "人员", "", "", adminList, "", map[string]interface{}{
		"placeholder": "选择管理员",
	})

	pageBuilder.ListColumnAdd("account_id", "操作人ID", "text", nil)
	pageBuilder.ListColumnAdd("account_name", "操作人名称", "text", nil)
	pageBuilder.ListColumnAdd("account", "操作人账号", "text", nil)
	pageBuilder.ListColumnAdd("url", "链接地址", "text", nil)
	pageBuilder.ListColumnAdd("rule_name", "执行操作", "text", nil)
	pageBuilder.ListColumnAdd("notice", "备注", "text", nil)
	return nil, 0
}

func (a AdminLog) NodeListCondition(pageBuilder *builder.PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {

	for k, v := range condition {
		if v[0].(string) == "account_id" {
			if util.Interface2String(v[1]) == "" {
				condition = append(condition[0:k], condition[k+1:]...)
			}
		}
	}
	return condition, nil, 0
}
