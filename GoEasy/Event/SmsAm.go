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

type SmsAm struct {
}

func (that SmsAm) Do(eventName string, data ...interface{}) (error, int) {
	programs := data[0].(map[string]interface{})
	recordId := programs["record_id"].(int64)
	smsUrl := programs["sms_url"].(string)
	apiToken := programs["api_token"].(string)
	Content := programs["content"].(string)
	tel := programs["tel"].(string)
	sign := programs["sign"].(string)
	url := smsUrl + "?mobile=" + tel + "&content=【" + sign + "】" + Content + "&api_token=" + apiToken + "&order_id=" + util.Int642String(recordId)
	content, err := util.FastHttpGet(url)
	if err != nil {
		logger.Error(err.Error())
		return err, 500
	}
	
	expressInfo := map[string]interface{}{}
	util.JsonDecode(content, &expressInfo)
	
	if util.Interface2String(expressInfo["code"]) != "200" {
		return errors.New(expressInfo["msg"].(string) + "；code=" + expressInfo["code"].(string)), 201
	}
	result := expressInfo["data"].(map[string]interface{})
	if result["ReturnStatus"].(string) == "Success" {
		return nil, 200
	} else {
		return errors.New(result["Message"].(string)), 202
	}
}
