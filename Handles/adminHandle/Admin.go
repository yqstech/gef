/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-05-13 15:23:10
 * @LastEditTime: 2021-10-20 17:28:18
 * @Description: 管理员账户列表
 */

package adminHandle

import (
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/util"
	"github.com/yqstech/gef/builder"

	"github.com/gohouse/gorose/v2"
)

type Admin struct {
	Base
}

func (a Admin) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("后台账户管理")
	pageBuilder.SetPageName("后台账户")
	pageBuilder.SetTbName("tb_admin")
	return nil, 0
}

func (a Admin) NodeList(pageBuilder *builder.PageBuilder) (error, int) {

	//获取角色列表
	groupOptions, err, code := Models.Model{}.SelectOptionsData("tb_admin_group", map[string]string{
		"id":         "value",
		"group_name": "name",
	}, "0", "未分配", "", "")
	if err != nil {
		return err, code
	}
	pageBuilder.ListColumnAdd("group_id", "权限组/角色", "array", groupOptions)
	pageBuilder.ListColumnAdd("name", "名称", "text", nil)
	pageBuilder.ListColumnAdd("account", "登录账户", "text", nil)
	pageBuilder.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ByKey("status", true))
	return nil, 0
}

func (a Admin) NodeListCondition(pageBuilder *builder.PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {
	//修改查询条件
	condition = append(condition, []interface{}{
		"id", ">", 1,
	})
	return condition, nil, 0
}

func (a Admin) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
	groupOptions, err, code := Models.Model{}.SelectOptionsData("tb_admin_group", map[string]string{
		"id":         "value",
		"group_name": "name",
	}, "0", "暂不分配", "", "")
	if err != nil {
		return err, code
	}
	pageBuilder.FormFieldsAdd("group_id", "select", "分配角色", "", "0", false, groupOptions, "", nil)
	pageBuilder.FormFieldsAdd("name", "text", "管理员名称", "", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("account", "text", "登录账户名", "字母或数字", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("password", "password", "密码", "不修改密码请留空", "", false, nil, "", nil)
	return nil, 0
}

// NodeFormData 表单显示前修改数据
func (a Admin) NodeFormData(pageBuilder *builder.PageBuilder, data gorose.Data, id int64) (gorose.Data, error, int) {
	data["password"] = ""
	return data, nil, 0
}

// NodeSaveData 表单保存数据前使用
func (a Admin) NodeSaveData(pageBuilder *builder.PageBuilder, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	if postData["password"] != "" {
		postData["password"] = util.GetPassword(postData["password"].(string))
	}
	return postData, nil, 0
}
