/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: ExpressJdly
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

type ExpressJdly struct {
	JdAppKey string
}

func (that *ExpressJdly) ExpressTrack(expressNumber, Tel string) (map[string]interface{}, error) {
	logger.Alert("京东万象-零羽科技")
	// 修改顺丰快递参数
	if expressNumber[0:2] == "SF" {
		Tel = Tel[len(Tel)-4:]
		expressNumber = expressNumber + "-" + Tel
	}
	Url := "https://way.jd.com/lywl/kd?nuo=" + expressNumber + "&appkey=" + that.JdAppKey
	content, err := util.FastHttpGet(Url)
	if err != nil {
		return nil, err
	}
	logger.Alert(content)
	//格式化json数据
	expressInfo := map[string]interface{}{}
	util.JsonDecode(content, &expressInfo)
	//查询错误返回码
	if expressInfo["code"].(string) != "10000" {
		return nil, errors.New(expressInfo["msg"].(string) + "；code=" + expressInfo["code"].(string))
	}
	//结果集
	result := expressInfo["result"].(map[string]interface{})
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
