/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsHoldBack
 * @Version: 1.0.0
 * @Date: 2022/2/8 3:34 下午
 */

package Event

import (
	"errors"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
	"time"
)

type SmsHoldBack struct {
}

// Do 短信发送过滤，防止恶意发送短信
func (that SmsHoldBack) Do(eventName string, data ...interface{}) (error, int) {
	ps := data[0].(map[string]interface{})
	
	//手机号临时阻止
	holdBack, holdBackMsg := Models.Sms{}.IsHoldBackTel(ps["tel"].(string))
	if holdBack {
		return errors.New(holdBackMsg), 502
	}
	//ip地址临时阻止
	holdBack, holdBackMsg = Models.Sms{}.IsHoldBackIp(ps["ip"].(string))
	if holdBack {
		return errors.New(holdBackMsg), 502
	}
	//短信功能暂停
	holdBack, holdBackMsg = Models.Sms{}.IsHoldBackAll()
	if holdBack {
		return errors.New(holdBackMsg), 502
	}
	
	//查找黑白名单
	blackWhite, err := db.New().Table("tb_app_sms_black_white").
		Where("is_delete", 0).
		Where("status", 1).
		Get()
	if err != nil {
		logger.Error(err.Error())
		return errors.New("系统运行错误！"), 500
	}
	for _, item := range blackWhite {
		if item["type"].(int64) == 0 {
			//黑名单直接拦截
			if item["rule"].(string) == ps["tel"].(string) || item["rule"].(string) == ps["ip"].(string) {
				Models.Sms{}.HoldBackIp(ps["ip"].(string), "短信发送失败!", 120)
				Models.Sms{}.HoldBackTel(ps["tel"].(string), "短信发送失败!", 120)
				return errors.New("短信发送失败！"), 502
			}
		} else {
			//白名单直接放行
			if item["rule"].(string) == ps["tel"].(string) || item["rule"].(string) == ps["ip"].(string) {
				return nil, 200
			}
		}
	}
	//手机号和IP不在黑白名单
	//检查ip和手机号是否达到阻止条件
	//规则类型：按手机号，按IP地址，限制频次
	//规定时长：多少秒内
	//达到次数：短信达到2次
	//执行操作：临时冻结、加入黑名单
	//冻结分钟数
	holdBackRules, err := db.New().Table("tb_app_sms_hold_back").
		Where("is_delete", 0).
		Where("status", 1).
		Order("sms_max asc").
		Get()
	if err != nil {
		logger.Error(err.Error())
		return errors.New("系统运行错误！"), 500
	}
	t := time.Now().Unix()
	var returnErr error
	for _, rule := range holdBackRules {
		
		//#规则一
		if rule["rule_type"].(int64) == 0 {
			//按手机号统计多少秒内的次数
			if rule["range_second"].(int64) <= 0 || rule["sms_max"].(int64) <= 0 {
				//规则设置错误
				continue
			}
			formTime := util.UnixTimeFormat(t-rule["range_second"].(int64), "2006-01-02 15:04:05")
			sendCount, err := db.New().Table("tb_app_sms_record").
				Where("is_delete", 0).
				Where("tel", ps["tel"].(string)).
				Where("status", "<", 4).
				Where("create_time", ">=", formTime).
				Count()
			if err != nil {
				logger.Error(err.Error())
				return errors.New("系统运行错误！"), 500
			}
			if sendCount >= rule["sms_max"].(int64) {
				//满足条件
				if rule["action"].(int64) == 1 && rule["frozen_second"].(int64) > 0 {
					//临时冻结
					Models.Sms{}.HoldBackTel(ps["tel"].(string), "您的操作过于频繁！", time.Duration(rule["frozen_second"].(int64)))
					returnErr = errors.New("您的操作过于频繁！")
				} else if rule["action"].(int64) == 2 {
					//手机号拉黑
					Models.AppSmsBlackWhite{}.AddBlack(0, ps["tel"].(string), "[自动拉黑]"+rule["note"].(string))
					returnErr = errors.New("您的操作过于频繁！")
				}
			}
		} else if rule["rule_type"].(int64) == 1 { //#规则二
			//按ip统计多少秒内的次数
			if rule["range_second"].(int64) <= 0 || rule["sms_max"].(int64) <= 0 {
				//规则设置错误
				continue
			}
			formTime := util.UnixTimeFormat(t-rule["range_second"].(int64), "2006-01-02 15:04:05")
			sendCount, err := db.New().Table("tb_app_sms_record").
				Where("is_delete", 0).
				Where("ip", ps["ip"].(string)).
				Where("status", "<", 4).
				Where("create_time", ">=", formTime).
				Count()
			if err != nil {
				logger.Error(err.Error())
				return errors.New("系统运行错误！"), 500
			}
			if sendCount >= rule["sms_max"].(int64) {
				//满足条件
				if rule["action"].(int64) == 1 && rule["frozen_second"].(int64) > 0 {
					//临时冻结
					Models.Sms{}.HoldBackIp(ps["ip"].(string), "您的操作过于频繁！", time.Duration(rule["frozen_second"].(int64)))
					returnErr = errors.New("您的操作过于频繁！")
				} else if rule["action"].(int64) == 2 {
					//IP地址拉黑
					Models.AppSmsBlackWhite{}.AddBlack(1, ps["ip"].(string), "[自动拉黑]"+rule["note"].(string))
					returnErr = errors.New("您的操作过于频繁！")
				}
			}
		} else { //# 规则三
			//统计全部的
			if rule["range_second"].(int64) <= 0 || rule["sms_max"].(int64) <= 0 || rule["frozen_second"].(int64) <= 0 {
				//规则设置错误
				continue
			}
			formTime := util.UnixTimeFormat(t-rule["range_second"].(int64), "2006-01-02 15:04:05")
			sendCount, err := db.New().Table("tb_app_sms_record").
				Where("is_delete", 0).
				Where("status", "<", 4).
				Where("create_time", ">=", formTime).
				Count()
			if err != nil {
				logger.Error(err.Error())
				return errors.New("系统运行错误！"), 500
			}
			if sendCount >= rule["sms_max"].(int64) {
				//临时冻结
				Models.Sms{}.HoldBackAll(time.Duration(rule["frozen_second"].(int64)))
				returnErr = errors.New("短信发送频次超过限制！")
			}
		}
	}
	if returnErr != nil {
		return returnErr, 502
	}
	return nil, 0
}
