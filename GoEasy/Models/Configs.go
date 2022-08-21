/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: GroupConfigs
 * @Version: 1.0.0
 * @Date: 2021/10/24 11:56 上午
 */

package Models

import (
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/GoEasy/Utils/db"
	"github.com/yqstech/gef/GoEasy/Utils/util"
)

type Configs struct {
}

// GroupConfigs
// 按分组获取配置项
// 自动格式化options数据
func (ac Configs) GroupConfigs(groupId int) []map[string]interface{} {
	conn := db.New().Table("tb_configs")
	configInfo, err := conn.Where("is_delete", 0).
		Where("status", 1).
		Where("group_id", groupId).
		Order("index_num asc,id asc").
		Get()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	result := []map[string]interface{}{}
	for _, item := range configInfo {
		row := map[string]interface{}{
			"name":       item["name"].(string),
			"value":      item["value"].(string),
			"title":      item["title"].(string),
			"notice":     item["notice"].(string),
			"if":         item["if"].(string),
			"field_type": item["field_type"].(string),
			"options":    []map[string]interface{}{},
		}
		if item["options"].(string) != "" {
			options := []map[string]interface{}{}
			util.JsonDecode(item["options"].(string), &options)
			row["options"] = options
		}
		result = append(result, row)
	}
	return result
}
