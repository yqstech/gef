/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-04-11 11:18:52
 * @LastEditTime: 2021-05-15 22:26:44
 * @Description: 数组操作
 */
package util

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

/**
 * @description: 获取默认值（int类型）
 * @param {interface{}} v
 * @param {int} defval
 * @return {*}
 */
func GetDefInt(v interface{}, defval int) int {
	//值是nil 返回默认值
	if v == nil {
		return defval
	}
	//强制转换成字符串，字符串为空 则返回默认值
	s := Interface2String(v)
	if s == "" {
		return defval
	}
	//字符串转为 int 并返回
	return String2Int(s)
}

/**
 * @description: 获取默认值（string类型）
 * @param {interface{}} v
 * @param {string} defval
 * @return {*}
 */
func GetDefString(v interface{}, defval string) string {
	//值是nil 返回默认值
	if v == nil {
		return defval
	}
	//强制转换成字符串，字符串为空 则返回默认值
	s := Interface2String(v)
	if s == "" {
		return defval
	}
	return s
}

// ArrayMap2Tree
// []map[string]interface{} 转树状上下级结构
// 返回上下级数据，返回下级所有ID列表
func ArrayMap2Tree(data []map[string]interface{}, pid int64, idKey, pidKey, subsKey string) ([]map[string]interface{},
	[]int64, int64) {
	var result []map[string]interface{}
	children := []int64{}
	count := 0
	//子项的最大级别
	childmaxLevel := int64(0)
	for _, v := range data {
		//根据PID挑选出数据
		if v[pidKey] == pid {
			count++
			//迭代获取下一层级，下边的层级数
			next, _children, _lastLevel := ArrayMap2Tree(data, v["id"].(int64), idKey, pidKey, subsKey)
			if _lastLevel > childmaxLevel {
				childmaxLevel = _lastLevel
			}
			v[subsKey] = next
			v["_children"] = _children
			v["_lastLevel"] = _lastLevel
			//当前ID和下级ID都加入到 _children
			children = append(children, v[idKey].(int64))
			children = append(children, _children...)
			result = append(result, v)
		}
	}
	if count > 0 {
		childmaxLevel++
	}
	return result, children, childmaxLevel
}
