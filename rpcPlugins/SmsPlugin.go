/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 定义短信插件
 * @File: SmsPlugin
 * @Version: 1.0.0
 * @Date: 2022/3/5 10:13 上午
 */

package rpcPlugins

import (
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/wonderivan/logger"
	"net/rpc"
	"os"
	"os/exec"
)

//! 定义插件接口
//! 定义客户端插件，封装一个RPC客户端，客户端插件通过RPC协议调用服务端插件
//! 定义服务端插件，封装进一个目标插件，服务端插件调用目标插件实现功能
//! 插件容器，调用客户端插件和服务端插件，封装目标插件，初始化服务端插件时传递给服务端插件
//!【RPC客户端类】===call===> 【RPC服务端插件】===内部调用===> 真正的插件

//! 定义短信插件接口

// SmsPlugin 短信插件接口
type SmsPlugin interface {
	SendSms(programs map[string]interface{}) (int64, string) // SendSms 发短信方法
}

//!客户端的短信插件类，封装进一个Rpc客户端对象

// SmsPluginClient 客户端短信插件
type SmsPluginClient struct {
	client *rpc.Client
}

func (that *SmsPluginClient) SendSms(programs map[string]interface{}) (int64, string) {
	//将 map[string]interface{} 转为 interface{}类型
	var args interface{} = programs
	var resp map[string]interface{}
	//! rpc 远程调用 SmsPluginServer 的 SendSms方法
	err := that.client.Call("Plugin.SendSms", &args, &resp)
	if err != nil {
		logger.Error(err.Error())
		return 500, "系统运行错误！"
	}
	if code, ok := resp["code"]; ok {
		return code.(int64), resp["error"].(string)
	}
	return 500, "系统运行错误！"
}

//! 服务端的短信插件类，插件文件初始化时封装进真正的插件对象

// SmsPluginServer 服务端的短信插件类
type SmsPluginServer struct {
	Impl SmsPlugin
}

// SendSms 用于通信的同名方法，内部调用真正的插件实现
func (that *SmsPluginServer) SendSms(args interface{}, resp *map[string]interface{}) error {
	//调用初始化插件时传入的插件对象方法
	code, errMsg := that.Impl.SendSms(args.(map[string]interface{}))
	*resp = map[string]interface{}{
		"code": code, "error": errMsg,
	}
	return nil
}

//!定义插件容器，内部封装两个方法，通过Server方法获取服务端插件，通过Client获取客户端插件

// SmsPluginPack 插件容器
type SmsPluginPack struct {
	Impl SmsPlugin //!内嵌插件对象
}

// Server 插件进程调用
func (that *SmsPluginPack) Server(*plugin.MuxBroker) (interface{}, error) {
	//返回服务端插件，并传入真正的插件实现
	return &SmsPluginServer{Impl: that.Impl}, nil
}

// Client 此方法由宿主进程调用，返回客户端插件用来Rpc通信
func (SmsPluginPack) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &SmsPluginClient{client: c}, nil
}

// GetSmsPlugin
//!根据名称，获取Rpc插件的客户端实例，即实际返回的是（SmsPluginClient的实例）
func GetSmsPlugin(pluginName string) (SmsPlugin, error) {
	if _, ok := RpcPluginClients[pluginName]; !ok {
		//!不存在则加载插件
		SmsPluginLoad(pluginName)
	}

	if client, ok := RpcPluginClients[pluginName]; ok {
		//!通过【进程客户端】获取【Rpc协议客户端】，用于之后的通信
		rpcClient, err := client.Client()
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("获取协议客户端失败！")
		}
		//!根据插件集名称分配【插件实例】
		raw, err := rpcClient.Dispense("greeter")
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("获取客户端实例失败！")
		}
		//!返回的是插件客户端实例
		//!由于实现了短信插件接口的全部方法，所以可断言成插件实例
		return raw.(SmsPlugin), nil
	} else {
		//加载一遍还没有就是加载失败了
		return nil, errors.New("加载插件失败！")
	}
}

// SmsPluginLoad
//!通过文件名，加载插件文件并初始化进程连接对象
func SmsPluginLoad(pluginName string) {
	//!定义插件的日志对象
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	//!定义插件集（仅有一个）
	var pluginMap = map[string]plugin.Plugin{
		//插件名称到插件对象的映射关系
		"greeter": &SmsPluginPack{},
	}
	//! 创建进程之间的连接，返回【进程客户端】
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig, //默认插件验证方法
		Plugins:         pluginMap,
		Cmd:             exec.Command("./rpcPluginFiles/" + pluginName),
		Logger:          logger,
	})
	//! 保存到进程连接对象中，后期可一并销毁
	RpcPluginClients[pluginName] = client
}
