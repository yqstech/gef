/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AppConfigs
 * @Version: 1.0.0
 * @Date: 2022/3/9 8:24 上午
 */

package Models

import (
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Utils/db"
)

// AppConfigs 应用配置信息
type AppConfigs struct {
}

// Value 根据配置名获取配置信息
func (that AppConfigs) Value(name string) string {
	config, err := db.New().Table("tb_app_configs").Where("name", name).Where("is_delete", 0).First()
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	if config == nil {
		return ""
	}
	return config["value"].(string)
}

// Values 获取多个配置信息
func (that AppConfigs) Values(groupId int) map[string]string {
	configMap := map[string]string{}
	conn := db.New().Table("tb_app_configs")
	if groupId > 0 {
		conn = conn.Where("group_id", groupId)
	}
	configs, err := conn.Where("is_delete", 0).Get()
	if err != nil {
		logger.Error(err.Error())
		return configMap
	}
	if len(configs) > 0 {
		for _, cfg := range configs {
			configMap[cfg["name"].(string)] = cfg["value"].(string)
		}
	}
	return configMap
}
