/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsAli
 * @Version: 1.0.0
 * @Date: 2022/3/5 12:23 下午
 */

package main

import (
	"encoding/gob"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gef/rpcPlugins"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

type SmsAli struct {
	logger hclog.Logger
}

func (that *SmsAli) SendSms(programs map[string]interface{}) (int64, string) {
	for k, v := range programs {
		that.logger.Info(k, v)
	}
	//todo:======>通过短信内容，过滤一下参数，排除掉无用的参数
	client, err := dysmsapi.NewClientWithAccessKey("cn-zhangjiakou", programs["accessKeyId"].(string), programs["accessKeySecret"].(string))
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = programs["tel"].(string)
	request.SignName = programs["SignName"].(string)
	request.TemplateCode = programs["template_out_id"].(string)
	request.TemplateParam = util.JsonEncode(programs["params"])

	response, err := client.SendSms(request)
	if err != nil {
		return 500, err.Error()
	}
	if response.Code == "OK" {
		//短信回执
		return 200, response.BizId
	} else {
		return 500, response.Message + ";code=" + response.Code
	}
}

func init() {
	//Rpc传递消息需要注册此数据类型
	gob.Register(map[string]interface{}{})
}
func main() {
	//初始化插件
	smsPlugin := SmsAli{
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
