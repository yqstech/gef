/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AppSmsBlackWhite
 * @Version: 1.0.0
 * @Date: 2022/2/8 6:20 下午
 */

package Models

type AppSmsHoldBack struct {
}

var SmsHoldBackRuleTypes = []map[string]interface{}{
	{
		"name":  "单个手机号",
		"value": "0",
	},
	{
		"name":  "单个IP地址",
		"value": "1",
	},
	{
		"name":  "所有短信",
		"value": "2",
	},
}
var SmsHoldBackActions = []map[string]interface{}{
	{
		"name":  "暂停发短信",
		"value": "1",
	},
	{
		"name":  "加入黑名单",
		"value": "2",
	},
}
