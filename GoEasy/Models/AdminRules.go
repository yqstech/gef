/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-10-18 11:56:39
 * @LastEditTime: 2021-10-21 10:50:41
 * @Description:
 */
package Models

import (
	"errors"
	"github.com/gef/GoEasy/Utils/db"

	"github.com/wonderivan/logger"
)

type AdminRules struct {
}

//获取上级权限 0全部
func (ar AdminRules) GetParentMenus(Ruletype int64) ([]map[string]interface{}, error) {
	//搜索上级权限
	Ruletypes := []int64{1, 2}
	if Ruletype != 0 {
		Ruletypes = []int64{Ruletype}
	}
	menus, err := db.New().Table("tb_admin_rules").
		Where("is_delete", 0).
		Where("type", "in", Ruletypes).
		Order("index_num asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("系统运行出错")
	}
	menus = Model{}.GoroseDataLevelOrder(menus, "id", "pid", 0, 0)
	parentsMenus := []map[string]interface{}{
		{"value": "0", "name": "顶级模块"},
	}
	for _, v := range menus {
		if v["is_compel"].(int64) == 1 {
			v["name"] = v["name"].(string) + "[免检]"
		}
		if v["type"].(int64) == 2 {
			v["name"] = v["name"].(string) + "[权限]"
		}
		if v["level"] == int64(0) {
			parentsMenus = append(parentsMenus, map[string]interface{}{
				"value": v["id"].(int64),
				"name":  v["name"].(string),
			})
		} else if v["level"] == int64(1) {
			parentsMenus = append(parentsMenus, map[string]interface{}{
				"value": v["id"].(int64),
				"name":  "&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string),
			})
		} else if v["level"] == int64(2) {
			parentsMenus = append(parentsMenus, map[string]interface{}{
				"value": v["id"].(int64),
				"name":  "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string),
			})
		} else if v["level"] == int64(3) {
			parentsMenus = append(parentsMenus, map[string]interface{}{
				"value": v["id"].(int64),
				"name":  "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string),
			})
		}
	}
	return parentsMenus, nil
}

func (ar AdminRules) AllTypes() []map[string]interface{} {
	mt := []map[string]interface{}{
		{"value": 1, "name": "左侧导航菜单"},
		{"value": 2, "name": "按钮操作或页面"},
	}
	return mt
}
func (ar AdminRules) AllIsCompels() []map[string]interface{} {
	mt := []map[string]interface{}{
		{"value": 0, "name": "否(非必选需鉴权)"},
		{"value": 1, "name": "是(必选地址免检)"},
	}
	return mt
}
