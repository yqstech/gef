/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: GroupConfigs
 * @Version: 1.0.0
 * @Date: 2022/1/24 5:55 下午
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"net/http"
)

type AppConfigs struct {
	GroupId int
	Base
}

func (that AppConfigs) PageInit(pageData *EasyApp.PageData) {
	pageData.ActionAdd("edit2", that.Edit2)
}

// 查询分组名
func (that AppConfigs) GroupName() string {
	group, err := db.New().Table("tb_configs_group").Where("id", that.GroupId).First()
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	if group != nil {
		return group["group_name"].(string)
	}
	return ""
}

// NodeBegin 开始
func (that AppConfigs) NodeBegin(pageData *EasyApp.PageData) (error, int) {

	pageData.SetTitle(that.GroupName())
	pageData.SetPageName("设置")
	pageData.SetTbName("tb_app_configs")
	return nil, 0
}

// NodeList 初始化列表
func (that AppConfigs) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.ListRightBtnsClear()
	pageData.ListTopBtnsClear()
	//列表清除
	pageData.ListColumnClear()
	//添加两列信息
	pageData.ListColumnAdd("name", "名称", "text", nil)
	pageData.ListColumnAdd("value", "内容", "html", nil)
	pageData.SetListColumnStyle("name", "width:150px")
	//隐藏分页
	pageData.SetListPageHide()

	//新增右侧日志开启关闭按钮
	pageData.SetButton("edit2", EasyApp.Button{
		ButtonName: "编辑" + that.GroupName(),
		Action:     "edit2",
		ActionType: 2,
		ActionUrl:  "edit2",
		Class:      "def",
		Icon:       "ri-equalizer-line",
		Display:    "",
		Expand: map[string]string{
			"w": "98%",
			"h": "98%",
		},
	})
	pageData.SetListTopBtns("edit2")

	return nil, 0
}

// NodeListCondition 修改查询条件
func (that AppConfigs) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	//追加查询条件
	condition = append(condition, []interface{}{
		"group_id", that.GroupId,
	})
	return condition, nil, 0
}

// NodeListData 重写列表数据
func (that AppConfigs) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	//查询出的数据转化一下
	configValue := map[string]interface{}{}
	for _, v := range data {
		configValue[v["name"].(string)] = v["value"]
	}

	//!按配置顺序显示出来
	result := []gorose.Data{}
	//获取所有应用内配置项
	appConfigs := Models.Configs{}.GroupConfigs(that.GroupId)
	if len(appConfigs) > 0 {
		//遍历配置项
		for _, config := range appConfigs {
			//查找应用内对应配置项值
			value := config["value"].(string)
			if v, ok := configValue[config["name"].(string)]; ok {
				value = v.(string)
				if config["field_type"].(string) == "image" {
					value = "<img style=\"width:80px;max-hight:80px\" src=\"" + value + "\"/>"
				}
				if config["field_type"].(string) == "select" {
					for _, op := range config["options"].([]map[string]interface{}) {
						if op["value"] == value {
							value = op["name"].(string)
							break
						}
					}
				}
			}
			//数据添加
			result = append(result, gorose.Data{
				"name":  config["title"].(string),
				"value": value,
			})
		}
	}
	return result, nil, 0
}

func (that AppConfigs) Edit2(pageData *EasyApp.PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.NodeBegin(pageData)
	pageData.SetActionName("修改")
	appConfigs := Models.Configs{}.GroupConfigs(that.GroupId)
	if len(appConfigs) > 0 {
		for _, config := range appConfigs {
			expand := map[string]interface{}{}
			if config["if"].(string) != "" {
				expand["if"] = config["if"]
			}
			pageData.FormFieldsAdd(config["name"].(string), config["field_type"].(string),
				config["title"].(string),
				config["notice"].(string),
				config["value"].(string), false, config["options"].([]map[string]interface{}), "", expand)
		}
	}

	if r.Method == "POST" {
		PostData := util.PostJson(r, "formFields")

		for _, config := range appConfigs {
			keyName := config["name"].(string)
			value := PostData[keyName]
			//检索是否有数据，有就保存没有就插入
			cfg, err := db.New().Table("tb_app_configs").
				Where("is_delete", 0).
				Where("group_id", that.GroupId).
				Where("name", keyName).First()
			if err != nil {
				logger.Error(err.Error())
				that.ApiResult(w, 500, "系统运行出错!", "")
				return
			}
			if cfg == nil {
				_, err = db.New().Table("tb_app_configs").Insert(map[string]interface{}{
					"group_id": that.GroupId,
					"name":     keyName,
					"value":    value,
				})
				if err != nil {
					logger.Error(err.Error())
					that.ApiResult(w, 500, "系统运行出错!", "")
					return
				}
			} else {
				_, err := db.New().Table("tb_app_configs").
					Where("id", cfg["id"]).
					Update(map[string]interface{}{
						"value":       value,
						"update_time": util.TimeNow(),
					})
				if err != nil {
					logger.Error(err.Error())
					that.ApiResult(w, 500, "系统运行出错!", "")
					return
				}
			}
		}

		that.ApiResult(w, 200, "修改成功!", "success")
		return
	}

	cfgs, err := db.New().Table("tb_app_configs").
		Where("is_delete", 0).
		Where("group_id", that.GroupId).Get()
	if err != nil {
		logger.Error(err.Error())
		that.ApiResult(w, 500, "系统运行出错!", "")
		return
	}
	odata := gorose.Data{}
	for _, cfg := range cfgs {
		odata[cfg["name"].(string)] = cfg["value"]
	}
	pageData.SetFormData(odata)

	that.ActShow(w, EasyApp.Template{
		TplName: "edit.html",
	}, pageData)
}
