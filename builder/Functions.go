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
	"reflect"
	"strconv"
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
