/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-04-01 16:19:26
 * @LastEditTime: 2021-06-21 11:11:36
 * @Description: 基础控制器
 */

package commHandle

import (
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/util"
	"github.com/yqstech/gef/builder"
	"net/http"

	"github.com/gohouse/gorose/v2"
)

type Base struct {
	builder.NodePage
}

func (base Base) GetTabIndex(pageBuilder *builder.PageBuilder, keyName string) int {
	//获取tab序列号并设置选中序列
	tab := util.GetValue(pageBuilder.GetHttpRequest(), keyName)
	tabIndex := 0
	if tab != "" {
		tabIndex = util.String2Int(tab)
	}
	return tabIndex
}

// DbModel 快速获取数据模型
func (base Base) DbModel(tbname string) gorose.IOrm {
	DB := db.New()
	obj := DB.Table(tbname)
	return obj
}

func StringValue(r *http.Request, name string) string {
	value := util.PostValue(r, name)
	if value == "" {
		value = util.GetValue(r, name)
	}
	return value
}
func IntValue(r *http.Request, name string) int {
	value := StringValue(r, name)
	if value == "" {
		return 0
	}
	return util.String2Int(value)
}
func RankData(rules []map[string]interface{}, pid int64) []map[string]interface{} {
	var rulesArr []map[string]interface{}
	for _, rule := range rules {
		if rule["pid"].(int64) == pid {
			rulesArr = append(rulesArr, map[string]interface{}{
				"id":          rule["id"].(int64),
				"name":        rule["name"],
				"selected":    false,
				"subs":        RankData(rules, rule["id"].(int64)),
				"rule_values": "",
			})
		}
	}
	return rulesArr
}
