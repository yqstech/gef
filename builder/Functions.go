/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Functions
 * @Version: 1.0.0
 * @Date: 2022/8/21 21:49
 */

package builder

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/util"
	"reflect"
	"strconv"
	"time"
)

// IsNum 任意值转字符串后判断是否是数字
func IsNum(value interface{}) bool {
	s := Interface2String(value)
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func Interface2String(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func inArray(v interface{}, Arr []interface{}) bool {
	if len(Arr) == 0 {
		return false
	}
	for _, v1 := range Arr {
		if Interface2String(v) == Interface2String(v1) {
			return true
		}
	}
	return false
}

func JsonEncode(data interface{}) string {
	if data == nil {
		return "[]"
	}
	jsondata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%d", err.Error())
		return ""
	} else {
		return string(jsondata)
	}
}

func IsInArray(v interface{}, Arr []interface{}) bool {
	if len(Arr) == 0 {
		return false
	}
	for _, v1 := range Arr {
		if Interface2String(v) == Interface2String(v1) {
			return true
		}
	}
	return false
}
func inMap(k string, data map[string]interface{}) bool {
	if _, ok := data[k]; ok {
		return true
	}
	return false
}

// NodePageCopy 拷贝节点页对象
func NodePageCopy(m NodePager) NodePager {
	vt := reflect.TypeOf(m).Elem()
	newObj := reflect.New(vt)
	newObj.Elem().Set(reflect.ValueOf(m).Elem())
	return newObj.Interface().(NodePager)
}

// 查询列表 tableName pageSize page where orderBy
func dbSelect(options ...interface{}) []gorose.Data {
	optionsLen := len(options)
	tableName := ""
	if optionsLen > 0 {
		tableName = util.Interface2String(options[0])
	} else {
		return nil
	}
	pageSize := 10
	if optionsLen > 1 {
		pageSize = util.String2Int(util.Interface2String(options[1]))
	}
	page := 1
	if optionsLen > 2 {
		page = util.String2Int(util.Interface2String(options[2]))
	}
	where := "is_delete = 0"
	if optionsLen > 3 {
		getWhere := util.Interface2String(options[3])
		if getWhere != "" {
			where = where + " and " + getWhere
		}
	}
	orderBy := "id desc"
	if optionsLen > 4 {
		orderBy = util.Interface2String(options[4])
	}

	conn := db.New()
	data, err := conn.Table(tableName).
		Where(where).
		Limit(pageSize).
		Page(page).
		OrderBy(orderBy).Get()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	for i, d := range data {
		//!针对于time类型（sqlite数据库）转为字符串
		for k, v := range d {
			switch v.(type) {
			case time.Time:
				d[k] = v.(time.Time).Format("2006-01-02 15:04:05")
			default:
			}
		}
		data[i] = d
	}
	return data
}

// 查询列表 tableName where
func dbCount(options ...interface{}) int64 {
	optionsLen := len(options)
	tableName := ""
	if optionsLen > 0 {
		tableName = util.Interface2String(options[0])
	} else {
		return 0
	}
	where := "is_delete = 0"
	if optionsLen > 1 {
		getWhere := util.Interface2String(options[1])
		if getWhere != "" {
			where = where + " and " + getWhere
		}
	}
	conn := db.New()
	count, err := conn.Table(tableName).Where(where).Count()
	if err != nil {
		logger.Error(err.Error())
		return 0
	}
	return count
}

// 模板查询单条数据,tableName where orderBy
func dbFind(options ...interface{}) gorose.Data {
	optionsLen := len(options)
	tableName := ""
	if optionsLen > 0 {
		tableName = util.Interface2String(options[0])
	}
	where := "is_delete = 0"
	if optionsLen > 1 {
		getWhere := util.Interface2String(options[1])
		if getWhere != "" {
			where = where + " and " + getWhere
		}
	}
	orderBy := "id desc"
	if optionsLen > 2 {
		orderBy = util.Interface2String(options[2])
	}

	conn := db.New()
	data, err := conn.Table(tableName).Where(where).OrderBy(orderBy).First()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	for k, v := range data {
		//时间格式转成字符串格式
		switch v.(type) {
		case time.Time:
			data[k] = v.(time.Time).Format("2006-01-02 15:04:05")
		default:
		}
	}
	return data
}

// 传入选项集key和选项key，得到选项名称
func optionName(optionModelsKey string, optionKey string) string {
	optionList := Models.OptionModels{}.ByKey(optionModelsKey, false)
	optionMap := map[string]interface{}{}
	for _, opt := range optionList {
		optionMap[util.Interface2String(opt["value"])] = opt["name"]
	}
	if v2, ok2 := optionMap[optionKey]; ok2 {
		return util.Interface2String(v2)
	}
	return ""
}

func isIn(options ...interface{}) bool {
	if len(options) < 2 {
		return false
	}
	first := util.Interface2String(options[0])
	for k, v := range options {
		if k > 0 {
			if first == util.Interface2String(v) {
				return true
			}
		}
	}
	return false
}
