/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-05-13 15:23:10
 * @LastEditTime: 2021-10-19 22:14:13
 * @Description:后台权限组（角色）管理
 */

package adminHandle

import (
	"errors"
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"

	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
)

type AdminGroup struct {
	Base
}

func (ad AdminGroup) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetPageName("角色")
	pageData.SetTitle("角色管理")
	pageData.SetTbName("tb_admin_group")
	return nil, 0
}

func (ad AdminGroup) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.ListColumnAdd("group_name", "角色名称", "text", nil)
	pageData.ListColumnAdd("status", "状态", "array", Models.OptionModels{}.ById(2, true))
	return nil, 0
}

func (ad AdminGroup) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	//获取所有权限列表
	conn := db.New().Table("tb_admin_rules")
	rules, err := conn.Where("is_delete", 0).
		Where("status", 1).
		Where("type", 1).
		Order("index_num asc").
		Get()
	if err != nil {
		logger.Error(err.Error())
		return errors.New("出错了"), 500
	}

	//goroseData 转成 []map
	var data []map[string]interface{}
	for _, v := range rules {
		if v["is_compel"].(int64) == 1 {
			v["name"] = v["name"].(string) + "[免检]"
		}
		if v["type"].(int64) == 2 {
			v["name"] = v["name"].(string) + "[权限]"
		}
		data = append(data, map[string]interface{}{
			"id":        v["id"],
			"pid":       v["pid"],
			"name":      v["name"],
			"type":      v["type"],
			"is_compel": v["is_compel"],
		})
	}
	//[]map转成上下级结构
	data, _, _ = util.ArrayMap2Tree(data, 0, "id", "pid", "_child")
	//表单信息
	pageData.FormFieldsAdd("group_name", "text", "角色名称", "", "", true, nil, "", nil)
	pageData.FormFieldsAdd("rules", "checkbox_level", "配置权限", "", "", false, data, "", nil)
	return nil, 0
}

/**
 * @description:
 * @param {gorose.Data} formData
 * @param {map[string]interface{}} postData
 * @return {*}
 */
func (ad AdminGroup) NodeSaveData(pageData *EasyApp.PageData, formData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	postData["rules"] = util.JsonEncode(postData["rules"])
	return postData, nil, 0
}
