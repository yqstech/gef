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
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"

	"github.com/gohouse/gorose/v2"
)

type AdminRules struct {
	Base
}

// NodeBegin 定义页面名称，数据库信息等
func (ar AdminRules) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("权限管理")
	pageBuilder.SetPageName("权限")
	pageBuilder.SetTbName("tb_admin_rules")
	return nil, 0
}

// NodeList 列表开始
func (ar AdminRules) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	//列表查询条件
	pageBuilder.SetListPageSize(200)
	pageBuilder.SetListOrder("index_num,id asc")

	//设置列表项
	pageBuilder.ListColumnAdd("name", "权限名称", "html", nil)
	pageBuilder.ListColumnAdd("route", "权限地址", "", nil)
	pageBuilder.ListColumnAdd("icon", "图标", "icon", nil)
	pageBuilder.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ByKey("status", true))
	pageBuilder.ListColumnAdd("is_compel", "必选", "array", Models.OptionModels{}.ByKey("is", true))
	pageBuilder.ListColumnAdd("type", "权限类型", "array", Models.OptionModels{}.ByKey("rule_type", true))
	pageBuilder.ListColumnAdd("open_log", "日志", "switch::text=开启|关闭", nil)
	pageBuilder.ListColumnAdd("index_num", "排序", "input::width=50px&type=number", nil)

	//设置搜索表单
	//pageBuilder.ListSearchFieldAdd("time_type", "select", "按订单时间", "1", nil, "width:auto;min-width:0px", nil)
	//pageBuilder.ListSearchFieldAdd("start_time", "datetime", "", util.TimeNowFormat("2006-01-02 00:00:00", 0, 0, -2), nil, "", nil)
	//pageBuilder.ListSearchFieldAdd("end_time", "datetime", "-", util.TimeNowFormat("2006-01-02 00:00:00", 0, 0, +1), nil, "", nil)
	pageBuilder.ListSearchFieldAdd("id", "text-sm", "ID", "", "", nil, "", map[string]interface{}{
		"placeholder": "ID",
	})
	pageBuilder.ListSearchFieldAdd("status", "select", "状态", "-1", "-1", Models.OptionModels{}.ByKey("status", false), "", map[string]interface{}{
		"placeholder": "选择状态",
	})
	pageBuilder.ListSearchFieldAdd("type", "select", "类型", "0", "0", Models.OptionModels{}.ByKey("rule_type", false), "", map[string]interface{}{
		"placeholder": "选择类型",
	})
	return nil, 0
}

func (ar AdminRules) NodeListCondition(pageBuilder *builder.PageBuilder, data [][]interface{}) ([][]interface{}, error, int) {
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
func (ar AdminRules) NodeListData(pageBuilder *builder.PageBuilder, data []gorose.Data) ([]gorose.Data, error, int) {
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
func (ar AdminRules) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
	//获取上级菜单
	parentsMenus, err := Models.AdminRules{}.GetParentMenus(0)
	if err != nil {
		return err, 120
	}
	pageBuilder.FormFieldsAdd("", "block", "权限信息", "", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("name", "text", "权限名称", "权限名称", "", true, nil, "", nil)
	pageBuilder.FormFieldsAdd("pid", "select", "上级权限", "", "0", true, parentsMenus, "", nil)
	pageBuilder.FormFieldsAdd("route", "text", "链接地址", "例如:/index/index", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("icon", "icon", "Icon图标", "请选择图标（remixicon）", "", false, nil, "", nil)
	pageBuilder.FormFieldsAdd("type", "radio", "权限类型", "", "1", false, Models.OptionModels{}.ByKey("rule_type", false), "", nil)
	pageBuilder.FormFieldsAdd("is_compel", "radio", "必选权限", "", "0", true, Models.AdminRules{}.AllIsCompels(), "", nil)
	pageBuilder.FormFieldsAdd("index_num", "number-xxs", "菜单排序", "值越小越靠前", "200", true, nil, "", nil)

	return nil, 0
}

// NodeSaveData 表单保存数据前使用
func (ar AdminRules) NodeSaveData(pageBuilder *builder.PageBuilder, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	postData["is_inside"] = 0
	return postData, nil, 0
}
