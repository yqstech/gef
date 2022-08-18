/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: database
 * @Version: 1.0.0
 * @Date: 2022/8/16 21:59
 */

package gef

import (
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/gdb"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
)

// DbManager 数据库管理
type DbManager struct {
}

// AutoTable 自动维护数据库结构
func (that DbManager) AutoTable(tables []interface{}) {
	err := gdb.New().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(tables...)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

// AutoAdminRules 后台权限表
func (that DbManager) AutoAdminRules(rules []map[string]interface{}) {
	setAdminRules(0, rules)
}
func setAdminRules(pid int64, rules []map[string]interface{}) {
	for index, rule := range rules {
		//route参数必填
		if rule["route"].(string) == "" {
			continue
		}
		ruleId := int64(0)
		//查找已经插入的数据
		ruleInfo, err := db.New().Table("tb_admin_rules").
			Where("is_delete", 0).
			Where("pid", pid). //同一个pid下不能有重复的route
			Where("route", rule["route"].(string)).
			First()
		if err != nil {
			panic("权限更新失败！" + err.Error())
			return
		}
		
		//否则新增这个数据
		if ruleInfo == nil {
			IndexNum := index + 1
			if indexNum, ok := rule["index_num"]; ok {
				IndexNum = indexNum.(int)
			}
			newData := map[string]interface{}{
				"pid":         pid,
				"name":        rule["name"],
				"type":        rule["type"],
				"is_compel":   rule["is_compel"],
				"icon":        rule["icon"],
				"route":       rule["route"],
				"index_num":   IndexNum,
				"create_time": util.TimeNow(),
				"update_time": util.TimeNow(),
			}
			//存在状态字段，则设置状态，否则是默认的1
			if ruleStatus, ok := rule["status"]; !ok {
				newData["status"] = ruleStatus
			}
			insertId, err := db.New().Table("tb_admin_rules").InsertGetId(newData)
			if err != nil {
				panic("权限更新失败！" + err.Error())
				return
			}
			ruleId = insertId
		} else {
			ruleId = ruleInfo["id"].(int64)
			updateData := map[string]interface{}{
				"name":        rule["name"],
				"type":        rule["type"],
				"is_compel":   rule["is_compel"],
				"icon":        rule["icon"],
				"update_time": util.TimeNow(),
			}
			//存在状态字段，则更新状态
			if ruleStatus, ok := rule["status"]; !ok {
				updateData["status"] = ruleStatus
			}
			//值不全，不更新
			if _, ok := rule["name"]; !ok {
				goto EndUpdate
			}
			if _, ok := rule["type"]; !ok {
				goto EndUpdate
			}
			if _, ok := rule["is_compel"]; !ok {
				goto EndUpdate
			}
			if _, ok := rule["icon"]; !ok {
				goto EndUpdate
			}
			//锁定后不更改
			if ruleInfo["is_lock"].(int64) == 1 {
				goto EndUpdate
			}
			//更新一下
			_, err = db.New().Table("tb_admin_rules").Where("id", ruleId).Update(updateData)
			logger.Alert("update " + rule["name"].(string))
			if err != nil {
				panic(err.Error())
				return
			}
		EndUpdate:
		}
		if children, ok := rule["children"]; ok && len(children.([]map[string]interface{})) > 0 {
			setAdminRules(ruleId, children.([]map[string]interface{}))
		}
	}
}

// AutoInsideData 内置数据维护
func (that DbManager) AutoInsideData(data []InsideData) {
	for _, d := range data {
		if d.TableName != "" && len(d.Condition) > 0 {
			conn := db.New().Table(d.TableName)
			for _, c := range d.Condition {
				conn.Where(c...)
			}
			first, err := conn.First()
			if err != nil {
				panic(err.Error())
				return
			}
			if first == nil {
				_, err := db.New().Table(d.TableName).Insert(d.Data)
				if err != nil {
					panic(err.Error())
					return
				}
			}
		}
	}
}

type InsideData struct {
	TableName string                 //数据表
	Condition [][]interface{}        //查询条件
	Data      map[string]interface{} //存储数据
}
