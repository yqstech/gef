/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 后台权限和菜单管理
 * @File: AdminRules
 * @Version: 1.0.0
 * @Date: 2021/10/15 5:54 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/util"

	"github.com/gohouse/gorose/v2"
)

type AdminRules struct {
	Base
}

// NodeBegin 定义页面名称，数据库信息等
func (ar AdminRules) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("权限管理")
	pageData.SetPageName("权限")
	pageData.SetTbName("tb_admin_rules")
	return nil, 0
}

// NodeList 列表开始
func (ar AdminRules) NodeList(pageData *EasyApp.PageData) (error, int) {
	//列表查询条件
	pageData.SetListPageSize(200)
	pageData.SetListOrder("index_num,id asc")

	//设置列表项
	pageData.ListColumnAdd("name", "权限名称", "html", nil)
	pageData.ListColumnAdd("route", "权限地址", "", nil)
	pageData.ListColumnAdd("icon", "图标", "icon", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	pageData.ListColumnAdd("is_compel", "必选", "array", Models.OptionModels{}.ById(1, true))
	pageData.ListColumnAdd("type", "权限类型", "array", Models.OptionModels{}.ById(5, true))
	pageData.ListColumnAdd("open_log", "日志", "switch::text=开启|关闭", nil)
	pageData.ListColumnAdd("index_num", "排序", "input::width=50px&type=number", nil)

	//设置搜索表单
	//pageData.ListSearchFieldAdd("time_type", "select", "按订单时间", "1", nil, "width:auto;min-width:0px", nil)
	//pageData.ListSearchFieldAdd("start_time", "datetime", "", util.TimeNowFormat("2006-01-02 00:00:00", 0, 0, -2), nil, "", nil)
	//pageData.ListSearchFieldAdd("end_time", "datetime", "-", util.TimeNowFormat("2006-01-02 00:00:00", 0, 0, +1), nil, "", nil)
	pageData.ListSearchFieldAdd("id", "text", "ID", "", "", nil, "", nil)
	pageData.ListSearchFieldAdd("status", "select", "状态", "-1", "-1", Models.OptionModels{}.ById(2, false), "", nil)
	pageData.ListSearchFieldAdd("type", "select", "类型", "0", "0", Models.OptionModels{}.ById(5, false), "", nil)
	return nil, 0
}

func (ar AdminRules) NodeListCondition(pageData *EasyApp.PageData, data [][]interface{}) ([][]interface{}, error, int) {
	serachData := [][]interface{}{}
	for _, v := range data {
		if v[0].(string) == "status" {
			if util.Interface2String(v[1]) != "-1" {
				serachData = append(serachData, []interface{}{v})
			}
		}
		if v[0].(string) == "id" && v[1].(string) != "" {
			serachData = append(serachData, []interface{}{v})
		}
		if v[0].(string) == "type" && v[1].(string) != "0" {
			serachData = append(serachData, []interface{}{v})
		}
		//if v[0].(string)=="start_time"{
		//	serachData = append(serachData,[]interface{}{"update_time",">=",v[1]})
		//}
		//if v[0].(string)=="end_time"{
		//	serachData = append(serachData,[]interface{}{"update_time","<",v[1]})
		//}
	}
	return serachData, nil, 0
}

// NodeListData 重写数据
func (ar AdminRules) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	//将数据按上下级顺序重新排列
	data = Models.Model{}.GoroseDataLevelOrder(data, "id", "pid", 0, 0)
	for k, v := range data {
		if v["level"] == int64(0) {
			data[k]["name"] = "<i class=\"layui-icon " + v["icon"].(string) + "\"></i>&nbsp;&nbsp;" + v["name"].(string)
		} else if v["level"] == int64(1) {
			data[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		} else if v["level"] == int64(2) {
			data[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		} else if v["level"] == int64(3) {
			data[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		} else if v["level"] == int64(4) {
			data[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
		}
	}
	return data, nil, 0
}

// NodeForm 表单初始化数据
func (ar AdminRules) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	//获取上级菜单
	parentsMenus, err := Models.AdminRules{}.GetParentMenus(0)
	if err != nil {
		return err, 120
	}
	pageData.FormFieldsAdd("", "block", "权限信息", "", "", false, nil, "", nil)
	pageData.FormFieldsAdd("name", "text", "权限名称", "权限名称", "", true, nil, "", nil)
	pageData.FormFieldsAdd("pid", "select", "上级权限", "", "0", true, parentsMenus, "", nil)
	pageData.FormFieldsAdd("route", "text", "链接地址", "例如:/index/index", "", false, nil, "", nil)
	pageData.FormFieldsAdd("icon", "icon", "Icon图标", "请选择图标（remixicon）", "", false, nil, "", nil)
	pageData.FormFieldsAdd("type", "radio", "权限类型", "", "1", false, Models.OptionModels{}.ById(5, false), "", nil)
	pageData.FormFieldsAdd("is_compel", "radio", "必选权限", "", "0", true, Models.AdminRules{}.AllIsCompels(), "", nil)
	pageData.FormFieldsAdd("index_num", "number", "菜单排序", "", "200", true, nil, "", nil)

	return nil, 0
}
