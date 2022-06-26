/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Sms
 * @Version: 1.0.0
 * @Date: 2022/2/7 10:48 下午
 */

package Models

import (
	"github.com/gef/GoEasy/Utils/pool"
	"time"
)

type Sms struct {
}

// SaveCode 保存验证码
func (s Sms) SaveCode(appId, tel, code string, second time.Duration) {
	pool.Gocache.Set(appId+"_"+tel+"_"+code, "1", time.Second*second)
}

// CheckCode 校验验证码
func (s Sms) CheckCode(appId, tel, code string) bool {
	value, ok := pool.Gocache.Get(appId + "_" + tel + "_" + code)
	if ok && value.(string) == "1" {
		return true
	} else {
		return false
	}
}

// HoldBackTel 阻止手机号
func (s Sms) HoldBackTel(appId, tel string, msg string, second time.Duration) {
	pool.Gocache.Set(appId+"_"+tel+"_SmsStop", msg, time.Second*second)
}

// HoldBackIp 阻止IP地址
func (s Sms) HoldBackIp(appId, ip string, msg string, second time.Duration) {
	pool.Gocache.Set(appId+"_"+ip+"_SmsStop", msg, time.Second*second)
}

// HoldBackAll 阻止IP地址
func (s Sms) HoldBackAll(appId string, second time.Duration) {
	pool.Gocache.Set(appId+"_"+"All_SmsStop", "短信发送频次超过限制！", time.Second*second)
}

// IsHoldBackTel 阻止手机号
func (s Sms) IsHoldBackTel(appId string, tel string) (bool, string) {
	stopMsg, ok := pool.Gocache.Get(appId + "_" + tel + "_SmsStop")
	if ok && stopMsg != nil {
		return true, stopMsg.(string)
	} else {
		return false, ""
	}
}

// IsHoldBackIp 阻挡Ip地址
func (s Sms) IsHoldBackIp(appId string, ip string) (bool, string) {
	stopMsg, ok := pool.Gocache.Get(appId + "_" + ip + "_SmsStop")
	if ok && stopMsg != nil {
		return true, stopMsg.(string)
	} else {
		return false, ""
	}
}

// IsHoldBackAll 阻挡全部的地址
func (s Sms) IsHoldBackAll(appId string) (bool, string) {
	stopMsg, ok := pool.Gocache.Get(appId + "_" + "All_SmsStop")
	if ok && stopMsg != nil {
		return true, stopMsg.(string)
	} else {
		return false, ""
	}
}
