/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-04-01 16:19:26
 * @LastEditTime: 2021-06-21 11:11:36
 * @Description: 基础控制器
 */

package commHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"net/http"
	"strings"

	"github.com/gohouse/gorose/v2"
)

type Base struct {
	EasyApp.Page
}

func (base Base) GetTabIndex(pageData *EasyApp.PageData, keyName string) int {
	//获取tab序列号并设置选中序列
	tab := util.GetValue(pageData.GetHttpRequest(), keyName)
	tabIndex := 0
	if tab != "" {
		tabIndex = util.String2Int(tab)
	}
	return tabIndex
}

//
////定义空白方法
//func (base Base) Empty(w http.ResponseWriter, r *http.Request, _ httprouter.UrlParams) {
//
//}

//快速获取数据模型
func (base Base) DbModel(tbname string) gorose.IOrm {
	db := db.New()
	obj := db.Table(tbname)
	return obj
}

//快速获取查询条件map结构
func (base Base) WhereObj() map[string]interface{} {
	return map[string]interface{}{"is_delete": 0}
}

func (base Base) DefPageInfo(r *http.Request) (int, int) {
	page := util.PostValue(r, "page")
	pageNum := util.String2Int(page)
	if pageNum < 1 {
		pageNum = 1
	}
	page_size := util.PostValue(r, "page_size")
	pageSize := util.String2Int(page_size)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return pageNum, pageSize
}

/**
 * @description: 获取IP地址
 * @param {*http.Request} r
 * @return {*}
 */
func GetIp(r *http.Request) string {
	IP := r.RemoteAddr
	IPinfo := strings.Split(IP, ":")
	IP_ADDR := IPinfo[0]
	RealIP := r.Header.Get("X-Real-IP")
	if RealIP != "" && IP_ADDR == "127.0.0.1" {
		IP_ADDR = RealIP
	}
	util.Echo("IP_ADDR", IP_ADDR)
	return IP_ADDR
}

/**
 * @description: 获取string参数
 * @param {*http.Request} r
 * @param {string} name
 * @return {*}
 */
func StringValue(r *http.Request, name string) string {
	value := util.PostValue(r, name)
	if value == "" {
		value = util.GetValue(r, name)
	}
	return value
}

/**
 * @description: 获取Int参数
 * @param {*http.Request} r
 * @param {string} name
 * @return {*}
 */
func IntValue(r *http.Request, name string) int {
	value := StringValue(r, name)
	if value == "" {
		return 0
	}
	return util.String2Int(value)
}

/**
 * @description: 数据转为上下级结构数据
 * @param {[]map[string]interface{}} rules
 * @param {int64} pid
 * @return {*}
 */
func RankData(rules []map[string]interface{}, pid int64) []map[string]interface{} {
	rules_arr := []map[string]interface{}{}
	//得到一级结构
	for _, rule := range rules {
		if rule["pid"].(int64) == pid {
			rules_arr = append(rules_arr, map[string]interface{}{
				"id":          rule["id"].(int64),
				"name":        rule["name"],
				"selected":    false,
				"subs":        RankData(rules, rule["id"].(int64)),
				"rule_values": "",
			})
		}
	}
	return rules_arr
}
