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
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
)

type SmsAli struct {
}

func (that SmsAli) Do(eventName string, data ...interface{}) (error, int) {
	programs := data[0].(map[string]interface{})
	for k, v := range programs {
		logger.Info(k, v)
	}
	client, err := dysmsapi.NewClientWithAccessKey("cn-zhangjiakou", programs["accessKeyId"].(string), programs["accessKeySecret"].(string))
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = programs["tel"].(string)
	request.SignName = programs["SignName"].(string)
	request.TemplateCode = programs["template_out_id"].(string)
	request.TemplateParam = util.JsonEncode(programs["params"])
	
	response, err := client.SendSms(request)
	if err != nil {
		return err, 500
	}
	if response.Code == "OK" {
		//短信回执 response.BizId
		return nil, 200
	} else {
		return errors.New(response.Message + ";code=" + response.Code), 500
	}
}
