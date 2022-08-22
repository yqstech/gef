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
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/gdb"
	"github.com/yqstech/gef/util"
	"strings"
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
		IndexNum := index + 1
		if indexNum, ok := rule["index_num"]; ok {
			IndexNum = indexNum.(int)
		}
		//否则新增这个数据
		if ruleInfo == nil {
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
			if ruleStatus, ok := rule["status"]; ok {
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
				"index_num":   IndexNum,
				"update_time": util.TimeNow(),
			}
			//存在状态字段，则更新状态
			if ruleStatus, ok := rule["status"]; ok {
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
			if ruleInfo["is_inside"].(int64) == 0 {
				goto EndUpdate
			}
			//更新一下
			_, err = db.New().Table("tb_admin_rules").Where("id", ruleId).Update(updateData)
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
	//遍历内部数据组
	for _, d := range data {
		//必须设置表名和查询条件
		if d.TableName != "" && len(d.Condition) > 0 {
			//查询一条数据
			conn := db.New().Table(d.TableName)
			for _, c := range d.Condition {
				conn = conn.Where(c...)
			}
			conn = conn.Where("is_delete", 0)
			first, err := conn.First()
			if err != nil {
				panic(err.Error())
				return
			}
			//logger.Debug(conn.LastSql())
			DataId := int64(0)
			//不存在则添加
			if first == nil {
				insertId, err := db.New().Table(d.TableName).InsertGetId(d.Data)
				if err != nil {
					panic(err.Error())
					return
				}
				DataId = insertId
			} else {
				//如果存在is_inside 且 is_inside=1时，更新数据
				if _, ok := first["is_inside"]; ok {
					if first["is_inside"].(int64) == 1 {
						db.New().Table(d.TableName).Where("id", first["id"]).Update(d.Data)
					}
				}
				DataId = first["id"].(int64)
			}
			//插入或查找成功，有下级，处理下级
			if DataId > 0 && len(d.Children) > 0 {
				//遍历下级数据,替换掉查询条件里的__PID__ 和 数据里的__PID__
				var newInsideDataList []InsideData
				for childIndex, childData := range d.Children {
					//处理下级的查询条件
					ccJson := util.JsonEncode(childData.(InsideData).Condition)
					ccJson = strings.Replace(ccJson, "__PID__", util.Int642String(DataId), -1)
					var ccArr [][]interface{}
					util.JsonDecode(ccJson, &ccArr)
					for ck, cv := range childData.(InsideData).Data {
						if util.Interface2String(cv) == "__PID__" {
							childData.(InsideData).Data[ck] = DataId
						}
					}
					newInsideData := InsideData{
						TableName: childData.(InsideData).TableName,
						Condition: ccArr,
						Data:      childData.(InsideData).Data,
						Children:  childData.(InsideData).Children,
					}
					newInsideDataList = append(newInsideDataList, newInsideData)
					d.Children[childIndex] = newInsideData
				}
				that.AutoInsideData(newInsideDataList)
			}
		}
	}
}

type InsideData struct {
	TableName string                 //数据表
	Condition [][]interface{}        //查询条件
	Data      map[string]interface{} //存储的数据
	Children  []interface{}          //下级数据集合
}
