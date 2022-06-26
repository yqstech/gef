/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsAm
 * @Version: 1.0.0
 * @Date: 2022/3/5 12:23 下午
 */

package main

import (
	"encoding/gob"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gef/rpcPlugins"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

type SmsAm struct {
	logger hclog.Logger
}

func (that *SmsAm) SendSms(programs map[string]interface{}) (int64, string) {
	recordId := programs["record_id"].(int64)
	smsUrl := programs["sms_url"].(string)
	apiToken := programs["api_token"].(string)
	Content := programs["content"].(string)
	tel := programs["tel"].(string)
	sign := programs["sign"].(string)
	url := smsUrl + "?mobile=" + tel + "&content=【" + sign + "】" + Content + "&api_token=" + apiToken + "&order_id=" + util.Int642String(recordId)
	content, err := util.FastHttpGet(url)
	if err != nil {
		that.logger.Error(err.Error())
		return 500, err.Error()
	}
	that.logger.Debug(content)

	expressInfo := map[string]interface{}{}
	util.JsonDecode(content, &expressInfo)

	if util.Interface2String(expressInfo["code"]) != "200" {
		return 201, expressInfo["msg"].(string) + "；code=" + expressInfo["code"].(string)
	}
	result := expressInfo["data"].(map[string]interface{})
	if result["ReturnStatus"].(string) == "Success" {
		return 200, ""
	} else {
		return 202, result["Message"].(string)
	}
}

func init() {
	//Rpc传递消息需要注册此数据类型
	gob.Register(map[string]interface{}{})
}
func main() {
	//初始化插件
	smsPlugin := SmsAm{
		logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		}),
	}
	//插件封装进容器
	SmsPluginPack := rpcPlugins.SmsPluginPack{Impl: &smsPlugin}
	//插件容器放入插件集内
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &SmsPluginPack,
	}
	//启动插件监听服务
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: rpcPlugins.HandshakeConfig, //Rpc插件的连接校验信息
		Plugins:         pluginMap,                  //监听的插件集合
	})
}
