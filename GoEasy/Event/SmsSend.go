/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsSend
 * @Version: 1.0.0
 * @Date: 2022/2/7 10:38 下午
 */

package Event

import (
	"errors"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gef/rpcPlugins"
	"github.com/wonderivan/logger"
)

type SmsSend struct {
}

func (that SmsSend) Do(eventName string, data ...interface{}) (error, int) {
	ps := data[0].(map[string]interface{})

	SmsUpstreams, err := db.New().Table("tb_sms_upstream").
		Where("is_delete", 0).
		Where("status", 1).Get()
	if err != nil {
		logger.Error(err.Error())
		return errors.New("系统运行错误"), 500
	}
	if len(SmsUpstreams) == 0 {
		return errors.New("系统未开启短信通道！"), 501
	}
	var UpstreamIds []interface{}
	UpstreamPlugins := map[int64]string{}
	for _, v := range SmsUpstreams {
		UpstreamPlugins[v["id"].(int64)] = v["plugin_name"].(string)
		UpstreamIds = append(UpstreamIds, v["id"])
	}

	//查找可用的通道列表
	appSmsUpstreams, err := db.New().Table("tb_app_sms_upstream").
		Where("is_delete", 0).
		Where("status", 1).
		WhereIn("upstream_id", UpstreamIds).
		Order("index_num asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return errors.New("系统运行错误"), 500
	}
	if len(appSmsUpstreams) == 0 {
		return errors.New("应用未开启短信通道！"), 503
	}
	for _, item := range appSmsUpstreams {
		//短信设置项
		upstreamId := item["upstream_id"].(int64)
		configs := map[string]interface{}{}
		util.JsonDecode(item["configs"].(string), &configs)
		//短信内容
		configs["content"] = ps["content"].(string)
		//外部模板ID
		configs["template_out_id"] = ps["template_out_id"].(string)
		//短信模板里的全部参数
		configs["params"] = ps["params"]
		//手机号码
		configs["tel"] = ps["tel"]

		//保存短信记录
		recordId, err := db.New().Table("tb_app_sms_record").
			InsertGetId(map[string]interface{}{
				"template_name":   ps["template_name"].(string),
				"upstream_id":     upstreamId,
				"tel":             ps["tel"],
				"ip":              ps["ip"],
				"template_out_id": ps["template_out_id"].(string),
				"content":         ps["content"].(string),
				"create_time":     util.TimeNow(),
				"update_time":     util.TimeNow(),
				"status":          2,
			})
		if err != nil {
			logger.Error(err.Error())
			return errors.New("系统运行错误"), 500
		}

		//短信记录ID
		configs["record_id"] = recordId

		//查找通道插件
		if pluginName, ok := UpstreamPlugins[upstreamId]; ok {
			if pluginName == "" {
				return errors.New("短信通道未指定插件！"), 503
			}
			//获取插件对象
			Plugin, err := rpcPlugins.GetSmsPlugin(pluginName)
			if err != nil {
				logger.Error(err.Error())
				//!短信发送标记失败
				db.New().Table("tb_app_sms_record").
					Where("id", recordId).
					Update(map[string]interface{}{
						"status": 4, "msg": "加载短信插件失败！",
					})
				return errors.New("加载短信插件失败！"), 503
			}
			//code, err := Plugin.SendFunc.(func(map[string]interface{}) (int64, error))(configs)
			code, errMsg := Plugin.SendSms(configs)
			if code == 200 {
				//!标记成功!
				db.New().Table("tb_app_sms_record").
					Where("id", recordId).
					Update(map[string]interface{}{
						"status": 3,
						"msg":    errMsg,
					})
				return nil, 200
			}
			if errMsg != "" {
				logger.Error(errMsg)
			}
			//!其他错误标记失败
			db.New().Table("tb_app_sms_record").
				Where("id", recordId).
				Update(map[string]interface{}{
					"status": 4, "msg": errMsg,
				})
		}
	}
	//短信发送记录
	return errors.New("短信发送失败，请稍后再试！"), 504
}
