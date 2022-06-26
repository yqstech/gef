/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsUpstreamParams
 * @Version: 1.0.0
 * @Date: 2021/10/28 10:16 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Utils/util"
)

type SmsUpstreamParams struct {
	Base
}

// NodeBegin 开始
func (that SmsUpstreamParams) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("短信通道配置项")
	pageData.SetPageName("通道配置项")
	pageData.SetTbName("tb_sms_upstream_params")
	return nil, 0
}

// NodeList 初始化列表
func (that SmsUpstreamParams) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.SetListOrder("id asc")
	//重设新增按钮增加
	ActionUrl := "add"
	upstreamId := util.GetValue(pageData.GetHttpRequest(), "id")
	if upstreamId != "" {
		ActionUrl = ActionUrl + "?upstream_id=" + upstreamId
		pageData.SetButton("add", EasyApp.Button{
			ButtonName: "新增通道配置项",
			Action:     "add",
			ActionType: 2,
			ActionUrl:  ActionUrl,
			Class:      "def",
			Icon:       "ri-play-list-add-fill",
			Display:    "",
			Expand: map[string]string{
				"w": "98%",
				"h": "98%",
			},
		})
	}

	//获取列表
	upstreamList := that.SmsUpstreamList()
	pageData.ListColumnAdd("upstream_id", "短信通道名称", "array", upstreamList)
	pageData.ListColumnAdd("param_name", "配置项", "text", nil)
	pageData.ListColumnAdd("param_title", "配置项名称", "text", nil)
	return nil, 0
}

// NodeListCondition 修改查询条件
func (that SmsUpstreamParams) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	upstreamID := 0
	upstreamId := util.GetValue(pageData.GetHttpRequest(), "id")
	if upstreamId != "" {
		upstreamID = util.String2Int(upstreamId)
		//追加查询条件
		condition = append(condition, []interface{}{
			"upstream_id", "=", upstreamID,
		})
	}

	return condition, nil, 0
}

// NodeForm 初始化表单
func (that SmsUpstreamParams) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	upstreamId := util.GetValue(pageData.GetHttpRequest(), "upstream_id")
	if upstreamId != "" {
		pageData.FormFieldsAdd("upstream_id", "hidden", "配置项", "", upstreamId, true, nil, "", nil)
	} else {
		//获取列表
		upstreamList := that.SmsUpstreamList()
		pageData.FormFieldsAdd("upstream_id", "select", "所属通道", "", "", true, upstreamList, "", nil)
	}

	pageData.FormFieldsAdd("param_name", "text", "配置项", "", "", true, nil, "", nil)
	pageData.FormFieldsAdd("param_title", "text", "配置项名称", "", "", true, nil, "", nil)
	return nil, 0
}
