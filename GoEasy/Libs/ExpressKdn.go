/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 快递鸟
 * @File: Express
 * @Version: 1.0.0
 * @Date: 2022/3/3 9:29 下午
 */

package Libs

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ExpressKdn struct {
	RequestType string
	AppID       string
	ApiKey      string
}

// GetParams 快递鸟 格式化参数
func (that *ExpressKdn) GetParams(expressNo, expressNumber, Tel string) (v map[string]string) {
	// 组装应用级参数
	if expressNo != "SF" {
		Tel = ""
	} else {
		//顺丰快递取后四位
		Tel = Tel[len(Tel)-4:]
	}
	RequestData := map[string]string{
		"CustomerName": Tel,
		"OrderCode":    "",
		"ShipperCode":  expressNo,
		"LogisticCode": expressNumber,
	}
	RequestDataJson := util.JsonEncode(RequestData)
	var DataSign string = that.getSign(RequestDataJson, that.ApiKey)
	// 组装系统级参数
	v = map[string]string{
		"RequestType": that.RequestType,
		"EBusinessID": that.AppID,
		"DataType":    "2",
		"RequestData": RequestDataJson,
		"DataSign":    DataSign,
	}
	return v
}

func (that *ExpressKdn) getSign(RequestData, ApiKey string) (data string) {
	str := RequestData + ApiKey
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	debyte := []byte(base64.StdEncoding.EncodeToString([]byte(md5str)))
	data = fmt.Sprintf("%s", debyte)
	return data
}

// Post 发送post信息
func (that *ExpressKdn) Post(params map[string]string) string {
	//正式发送地址
	postURL := "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"

	var values []string
	for k, v := range params {
		values = append(values, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	fmt.Println(values)
	resp, err := http.Post(postURL, "application/x-www-form-urlencoded", strings.NewReader(strings.Join(values, "&")))
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err.Error())
	}
	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(contentBytes)
}

// ExpressTrack 快递追踪信息
func (that *ExpressKdn) ExpressTrack(expressNo, expressNumber, Tel string) (map[string]interface{}, error) {
	params := that.GetParams(expressNo, expressNumber, Tel)
	content := that.Post(params)
	logger.Alert(content)

	expressInfo := map[string]interface{}{}
	util.JsonDecode(content, &expressInfo)

	if !expressInfo["Success"].(bool) {
		return nil, errors.New(expressInfo["Reason"].(string))
	}
	//格式化数据
	ExpressStatus := util.String2Int(expressInfo["State"].(string))
	ExpressMsg := ""
	if ExpressStatus == 4 {
		if StateEx, ok := expressInfo["StateEx"]; ok {
			switch StateEx {
			case "401":
				ExpressMsg = "发货无信息"
			case "402":
				ExpressMsg = "超时未签收"
			case "403":
				ExpressMsg = "超时未更新"
			case "404":
				ExpressMsg = "拒收(退件),"
			case "405":
				ExpressMsg = "派件异常"
			case "406":
				ExpressMsg = "退货签收"
			case "407":
				ExpressMsg = "退货未签收"
			case "412":
				ExpressMsg = "快递柜或驿站超时未取"
			default:
				ExpressMsg = "快件异常！"
			}
		}

	}
	ExpressInfo := map[string]interface{}{
		"ExpressStatus": ExpressStatus,
		"ExpressLast":   "",
		"ExpressMsg":    ExpressMsg,
		"ExpressTracks": []map[string]string{},
	}
	if len(expressInfo["Traces"].([]interface{})) > 0 {
		var ExpressTracks []map[string]string
		for _, item := range expressInfo["Traces"].([]interface{}) {
			ExpressTracks = append(ExpressTracks, map[string]string{
				"AcceptStation": item.(map[string]interface{})["AcceptStation"].(string),
				"AcceptTime":    item.(map[string]interface{})["AcceptTime"].(string),
				"Location":      util.Interface2String(item.(map[string]interface{})["Location"]),
			})
			ExpressInfo["ExpressLast"] = item.(map[string]interface{})["AcceptStation"].(string)
		}
		ExpressInfo["ExpressTracks"] = ExpressTracks
	}
	return ExpressInfo, nil
}
