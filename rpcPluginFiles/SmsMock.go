/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: SmsMock
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

type SmsMock struct {
	logger hclog.Logger
}

func (that *SmsMock) SendSms(programs map[string]interface{}) (int64, string) {
	that.logger.Error("rpcSmsMock插件接收到数据:" + util.JsonEncode(programs))
	return 200, ""
}

func init() {
	//Rpc传递消息需要注册此数据类型
	gob.Register(map[string]interface{}{})
}
func main() {
	//初始化插件
	smsPlugin := SmsMock{
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
