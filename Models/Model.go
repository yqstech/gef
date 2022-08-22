/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Comm
 * @Version: 1.0.0
 * @Date: 2021/10/19 12:10 下午
 */

package Models

import (
	"errors"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/util"

	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
)

type Model struct {
}

// DefaultStatus 通用状态列表
var DefaultStatus = []map[string]interface{}{
	{"value": 1, "name": "启用"},
	{"value": 0, "name": "禁用"},
}
var DefaultIsOrNot = []map[string]interface{}{
	{"value": 1, "name": "是"},
	{"value": 0, "name": "否"},
}

//
// SelectOptionsData
//  @Description: 通用的获取select组件数据或者列表array类型数据的方法，支持键名转换
//  @receiver mod
//  @param tbName string 表名称
//  @param keyTrans 数据键名转换，map[string(旧键名)]string(新键名)
//  @param defValue 默认值
//  @param defName 默认值名称
//  @return []map[string]interface{}
//  @return error
//  @return int
//
func (mod Model) SelectOptionsData(tbName string, keyTrans map[string]string, defValue, defName, where string, order string) ([]map[string]interface{}, error, int) {

	//设置默认值
	var result []map[string]interface{}
	if defValue != "" && defName != "" {
		result = append(result, map[string]interface{}{
			"value": defValue,
			"name":  defName,
		})
	}
	//查询数据库
	conn := db.New().Table(tbName)
	if where != "" {
		conn = conn.Where(where)
	}
	if order == "" {
		order = "id asc"
	}
	data, err := conn.Where("is_delete", 0).Order(order).Get()
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("系统错误！"), 500
	}
	//数据转成map数组并合并
	gList := mod.GoroseArrayToMap(data, keyTrans)
	if gList != nil {
		result = append(result, gList...)
	}

	//返回数组
	return result, nil, 0
}

//
// GoroseArrayToMap
//  @Description: goroseData数据转为map数组，支持键名转换，忽略未设置键名的数据
//  @receiver mod
//  @param data
//  @param keyTrans
//  @return []map[string]interface{}
//
func (mod Model) GoroseArrayToMap(data []gorose.Data, keyTrans map[string]string) []map[string]interface{} {
	var result []map[string]interface{}
	for _, d := range data {
		row := map[string]interface{}{}
		for oldKey, newKey := range keyTrans {
			if value, ok := d[oldKey]; ok {
				row[newKey] = util.Interface2String(value)
			}
		}
		result = append(result, row)
	}
	return result
}

// GoroseDataLevelOrder 将数据按上下级顺序排序并标记是第几级
func (mod Model) GoroseDataLevelOrder(data []gorose.Data, idKey string, pidKey string, pid, level int64) []gorose.Data {
	result := []gorose.Data{}
	for _, v := range data {
		if v[pidKey] == pid && level < 10 {
			v["level"] = level
			result = append(result, v)
			if pid != v["id"].(int64) {
				next := mod.GoroseDataLevelOrder(data, idKey, pidKey, v["id"].(int64), level+1)
				for _, nextItem := range next {
					result = append(result, nextItem)
				}
			}
		}
	}
	return result
}

// GoroseDataLevelTree 将数据转换成树形上下级
func (mod Model) GoroseDataLevelTree(data []gorose.Data, idKey string, pidKey string, pid int64, subsKey string) []gorose.Data {
	//创建当前一级数组
	result := []gorose.Data{}
	for _, v := range data {
		//根据PID挑选出数据
		if v[pidKey] == pid {
			//迭代获取下一层级
			next := mod.GoroseDataLevelTree(data, idKey, pidKey, v["id"].(int64), subsKey)
			v[subsKey] = next
			result = append(result, v)
		}
	}
	return result
}

// Select 万能查询方法
// 数据表 查询条件 排序条件 查询数量
func (mod Model) Select(dbName string, Where string, orderBy string, limit int, offset int) []map[string]string {
	var results []map[string]string
	conn := db.New().Table(dbName).
		Where("is_delete", 0).
		Where("status", 1)
	if Where != "" {
		conn = conn.Where(Where)
	}
	if orderBy == "" {
		orderBy = "id desc"
	}
	data, err := conn.Offset(offset).Limit(limit).Order(orderBy).Get()
	if err != nil {
		logger.Error(err.Error())
		return results
	}
	for _, row := range data {
		item := map[string]string{}
		for key, value := range row {
			item[key] = util.Interface2String(value)
		}
		results = append(results, item)
	}
	return results
}
