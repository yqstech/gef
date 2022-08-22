/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: FmtRuleFuncs
 * @Version: 1.0.0
 * @Date: 2022/3/31 3:13 下午
 */

package EasyCurd

import (
	"github.com/gohouse/gorose/v2"
	"github.com/yqstech/gef/util"
)

// CompleteDataFunc 定义完善数据方法
type CompleteDataFunc func(gorose.Data) gorose.Data

// FmtRuleFunc 定义类型，数据转换方法
type FmtRuleFunc func(interface{}) interface{}

// CompleteResultDataFunc 修改返回的数据
type CompleteResultDataFunc func(map[string]interface{}, string) interface{}

// JsonDecode 字符串转json
func JsonDecode(json interface{}) interface{} {
	var data interface{}
	util.JsonDecode(json.(string), &data)
	return data
}

// Interface2String 任意值转字符串
func Interface2String(data interface{}) interface{} {
	return util.Interface2String(data)
}
