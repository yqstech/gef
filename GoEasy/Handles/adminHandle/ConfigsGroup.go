/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: ConfigsGroup
 * @Version: 1.0.0
 * @Date: 2022/3/8 5:12 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
)

type ConfigsGroup struct {
	Base
}

// NodeBegin 开始
func (that ConfigsGroup) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("设置分组管理")
	pageData.SetPageName("设置分组")
	pageData.SetTbName("tb_configs_group")
	return nil, 0
}

// NodeList 初始化列表
func (that ConfigsGroup) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.SetListRightBtns("edit")
	pageData.SetListOrder("id asc")
	pageData.ListColumnAdd("group_name", "分组名称", "text", nil)
	pageData.ListColumnAdd("note", "分组备注", "text", nil)
	return nil, 0
}

// NodeForm 初始化表单
func (that ConfigsGroup) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("group_name", "text", "分组名称", "", "", true, nil, "", nil)
	pageData.FormFieldsAdd("note", "text", "分组备注", "", "", true, nil, "", nil)
	return nil, 0
}
