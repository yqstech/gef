/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsAli
 * @Version: 1.0.0
 * @Date: 2022/7/29 23:03
 */

package Event

import (
	"errors"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
)

type SmsJdcx struct {
}

func (that SmsJdcx) Do(eventName string, data ...interface{}) (error, int) {
	programs := data[0].(map[string]interface{})
	appKey := programs["appkey"].(string)
	Content := programs["content"].(string)
	tel := programs["tel"].(string)
	sign := programs["sign"].(string)
	url := "https://way.jd.com/chuangxin/dxjk?mobile=" + tel + "&content=【" + sign + "】" + Content + "&appkey=" + appKey
	content, err := util.FastHttpGet(url)
	if err != nil {
		logger.Error(err.Error())
		return err, 500
	}
	expressInfo := map[string]interface{}{}
	util.JsonDecode(content, &expressInfo)
	
	if expressInfo["code"].(string) != "10000" {
		return errors.New(expressInfo["msg"].(string) + "；code=" + expressInfo["code"].(string)), 201
	}
	result := expressInfo["result"].(map[string]interface{})
	if result["ReturnStatus"].(string) == "Success" {
		return nil, 200
	} else {
		return errors.New(result["Message"].(string)), 202
	}
}
