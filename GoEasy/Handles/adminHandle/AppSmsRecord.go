/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AppSmsRecord
 * @Version: 1.0.0
 * @Date: 2022/3/8 9:43 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
)

type AppSmsRecord struct {
	Base
}

// NodeBegin 开始
func (that AppSmsRecord) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("短信记录")
	pageData.SetPageName("短信记录")
	pageData.SetTbName("tb_app_sms_record")
	return nil, 0
}

// NodeList 初始化列表
func (that AppSmsRecord) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.ListColumnClear()
	pageData.SetListOrder("id desc")
	pageData.ListTopBtnsClear()
	pageData.ListRightBtnsClear()
	//获取列表
	upstreamOptions := that.SmsUpstreamList()
	pageData.ListColumnAdd("template_name", "短信模板", "text", nil)
	pageData.ListColumnAdd("upstream_id", "短信通道", "array", upstreamOptions)
	pageData.ListColumnAdd("tel", "手机号码", "text", nil)
	pageData.ListColumnAdd("ip", "ip地址", "text", nil)
	pageData.ListColumnAdd("content", "短信内容", "text", nil)
	pageData.ListColumnAdd("msg", "短信凭据/错误信息", "text", nil)
	pageData.ListColumnAdd("status", "状态", "array", []map[string]interface{}{
		{"value": "1", "name": "未发送"},
		{"value": "2", "name": "发送中"},
		{"value": "3", "name": "成功"},
		{"value": "4", "name": "失败"},
	})
	pageData.ListColumnAdd("create_time", "发送时间", "text", nil)
	pageData.SetListColumnStyle("msg", "max-width:200px;overflow-x:auto")
	pageData.PageStyle = "table td{word-break: keep-all;white-space:nowrap;}"
	return nil, 0
}
