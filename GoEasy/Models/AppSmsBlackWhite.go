/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AppSmsBlackWhite
 * @Version: 1.0.0
 * @Date: 2022/2/8 6:20 下午
 */

package Models

import (
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
)

type AppSmsBlackWhite struct {
}

var SmsBlackWhiteTypes = []map[string]interface{}{
	{
		"name":  "黑名单",
		"value": "0",
	},
	{
		"name":  "白名单",
		"value": "1",
	},
}
var SmsBlackWhiteRuleTypes = []map[string]interface{}{
	{
		"name":  "手机号",
		"value": "0",
	},
	{
		"name":  "IP地址",
		"value": "1",
	},
}

// AddBlack 加入黑名单
func (that AppSmsBlackWhite) AddBlack(ruleType int, rule string, note string) {
	_, err := db.New().Table("tb_app_sms_black_white").
		Insert(map[string]interface{}{
			"type":        0,
			"rule_type":   ruleType,
			"rule":        rule,
			"note":        note,
			"create_time": util.TimeNow(),
			"update_time": util.TimeNow(),
		})
	if err != nil {
		logger.Error(err.Error())
	}
}
