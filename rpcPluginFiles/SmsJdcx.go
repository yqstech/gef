/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsJdcx
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

type SmsJdcx struct {
	logger hclog.Logger
}

func (that *SmsJdcx) SendSms(programs map[string]interface{}) (int64, string) {
	//for k,v := range programs {
	//	that.logger.Info(k,v)
	//}
	appKey := programs["appkey"].(string)
	Content := programs["content"].(string)
	tel := programs["tel"].(string)
	sign := programs["sign"].(string)
	url := "https://way.jd.com/chuangxin/dxjk?mobile=" + tel + "&content=【" + sign + "】" + Content + "&appkey=" + appKey
	content, err := util.FastHttpGet(url)
	if err != nil {
		that.logger.Error(err.Error())
		return 500, err.Error()
	}
	that.logger.Debug(content)

	expressInfo := map[string]interface{}{}
	util.JsonDecode(content, &expressInfo)

	if expressInfo["code"].(string) != "10000" {
		return 201, expressInfo["msg"].(string) + "；code=" + expressInfo["code"].(string)
	}
	result := expressInfo["result"].(map[string]interface{})
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
	smsPlugin := SmsJdcx{
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
