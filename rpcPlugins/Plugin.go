/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 使用基于RPC 的 go-plugin实现插件
 * @File: Plugin.go
 * @Version: 1.0.0
 * @Date: 2022/3/5 10:08 上午
 */

package rpcPlugins

import (
	"github.com/gef/config"
	"github.com/hashicorp/go-plugin"
	"github.com/wonderivan/logger"
	"io/ioutil"
)

// PluginsLookup 检索项目目录的插件
func PluginsLookup(prefix string) ([]string, error) {
	var pluginNames []string
	//查找Rpc插件目录所有文件
	plugins, err := ioutil.ReadDir(config.WorkPath + "/plugins/" + config.GOOS)
	if err != nil {
		plugins, err = ioutil.ReadDir(config.AppPath + "/plugins/" + config.GOOS)
		if err != nil {
			return pluginNames, nil
		}
	}
	//查找所有的非目录文件
	for _, file := range plugins {
		if !file.IsDir() {
			pluginName := file.Name()
			//按前缀过滤
			if prefix != "" {
				if pluginName[0:len(prefix)] != prefix {
					continue
				}
			}
			pluginNames = append(pluginNames, pluginName)
		}
	}
	return pluginNames, nil
}

// HandshakeConfig 插件默认的握手验证信息
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PLUGIN_KEY",
	MagicCookieValue: "PLUGIN_VALUE",
}

// RpcPluginClients 存储Rpc插件客户端（进程连接客户端指针），用户最后的统一销毁
var RpcPluginClients = map[string]*plugin.Client{}

// RpcPluginClientsKill RPC连接全部销毁
func RpcPluginClientsKill() {
	for name, client := range RpcPluginClients {
		logger.Info("销毁Rpc插件:" + name)
		client.Kill()
	}
}
