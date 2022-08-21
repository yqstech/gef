/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: ExpressAm
 * @Version: 1.0.0
 * @Date: 2022/3/5 12:06 上午
 */

package Libs

import (
	"errors"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/GoEasy/Utils/util"
	"strings"
)

type ExpressAm struct {
	ApiUrl   string
	ApiToken string
}

func (that *ExpressAm) ExpressTrack(expressNumber, Tel string) (map[string]interface{}, error) {
	logger.Alert("云市场-快递查询")
	// 修改顺丰快递参数
	if expressNumber[0:2] == "SF" {
		Tel = Tel[len(Tel)-4:]
		expressNumber = expressNumber + "-" + Tel
	}
	Url := that.ApiUrl + "?nuo=" + expressNumber + "&api_token=" + that.ApiToken + "&order_id=" + util.Int642String(util.GetOnlyID())
	content, err := util.FastHttpGet(Url)
	if err != nil {
		return nil, err
	}
	logger.Alert(content)
	//格式化json数据
	expressInfo := map[string]interface{}{}
	util.JsonDecode(content, &expressInfo)
	//查询错误返回码
	if util.Interface2String(expressInfo["code"]) != "200" {
		return nil, errors.New(expressInfo["msg"].(string) + "；code=" + expressInfo["code"].(string))
	}
	//结果集
	result := expressInfo["data"].(map[string]interface{})
	//整理返回数据
	ExpressInfo := map[string]interface{}{
		"ExpressStatus": 0,
		"ExpressLast":   "",
		"ExpressMsg":    expressInfo["msg"].(string),
		"ExpressTracks": []map[string]string{},
	}
	//自主判断快递状态
	ExpressStatus := 0
	if strings.Contains(content, "签收") {
		ExpressStatus = 3
	} else if len(result["context"].([]interface{})) > 1 {
		ExpressStatus = 2
	} else if len(result["context"].([]interface{})) == 1 {
		ExpressStatus = 1
	}
	ExpressInfo["ExpressStatus"] = ExpressStatus

	//最后数据
	if v, ok := result["latest_progress"]; ok {
		ExpressInfo["ExpressLast"] = v.(string)
	}
	if len(result["context"].([]interface{})) > 0 {
		var ExpressTracks []map[string]string
		for _, item := range result["context"].([]interface{}) {
			ExpressTracks = append(ExpressTracks, map[string]string{
				"AcceptStation": item.(map[string]interface{})["desc"].(string),
				"AcceptTime":    util.UnixTimeFormat(int64(util.String2Int(util.Interface2String(item.(map[string]interface{})["time"]))), "2006-01-02 15:04:05"),
				"Location":      "",
			})
		}
		ExpressInfo["ExpressTracks"] = ExpressTracks
	}
	return ExpressInfo, nil
}
