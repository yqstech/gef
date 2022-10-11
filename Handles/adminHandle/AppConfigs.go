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
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"
	"net/http"
)

type AppConfigs struct {
	GroupId int
	Base
}

func (that *AppConfigs) NodeInit(pageBuilder *builder.PageBuilder) {
	that.NodePageActions["edit2"] = that.Edit2
	that.NodePageActions["dev_status"] = that.DevStatus
	that.NodePageActions["reset_server"] = that.ResetServer
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
func (that AppConfigs) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {

	pageBuilder.SetTitle(that.GroupName())
	pageBuilder.SetPageName("设置")
	pageBuilder.SetTbName("tb_app_configs")

	//自动清理重复项
	that.ClearAppConfigs()

	return nil, 0
}

func (that AppConfigs) ClearAppConfigs() {
	//!清理无效项，
	//查找所有有效的配置项
	ConfigItems, err := db.New().Table("tb_configs").
		Where("is_delete", 0).
		Order("id asc").Get()
	var configNames []interface{}
	for _, config := range ConfigItems {
		configNames = append(configNames, config["name"])
	}
	//删除所有无效的应用配置项
	db.New().Table("tb_app_configs").WhereNotIn("name", configNames).Delete()

	//!清理重复项
	appConfigs, err := db.New().Table("tb_app_configs").
		Where("is_delete", 0).Order("id asc").Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	appConfigNames := map[string]bool{}
	for _, appConfig := range appConfigs {
		//发现重复即删除
		if _, ok := appConfigNames[appConfig["name"].(string)]; ok {
			db.New().Table("tb_app_configs").Where("id", appConfig["id"]).Delete()
			continue
		}
		appConfigNames[appConfig["name"].(string)] = true
	}
}

// NodeList 初始化列表
func (that AppConfigs) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.ListRightBtnsClear()
	pageBuilder.ListTopBtnsClear()
	//列表清除
	pageBuilder.ListColumnClear()
	//添加两列信息
	pageBuilder.ListColumnAdd("name", "名称", "text", nil)
	pageBuilder.ListColumnAdd("key", "关键字", "text", nil)
	pageBuilder.ListColumnAdd("value", "内容", "html", nil)
	pageBuilder.SetListColumnStyle("name", "width:150px")
	pageBuilder.SetListColumnStyle("key", "width:150px")
	//隐藏分页
	pageBuilder.SetListPageHide()

	//新增右侧日志开启关闭按钮
	pageBuilder.SetButton("edit2", builder.Button{
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
	pageBuilder.SetListTopBtns("edit2")

	//新增开发模式开关,仅超级管理员有权限
	if that.GroupId == 1 && pageBuilder.HttpParams.ByName("account_id") == "1" {
		first, err := db.New().Table("tb_admin_rules").
			Where("is_delete", 0).
			Where("route", "#dev").
			Order("index_num desc").First()
		if err != nil {
			logger.Error(err.Error())
			return err, 0
		}
		if first == nil {
			return nil, 0
		}
		if first["status"].(int64) == 1 {
			pageBuilder.SetButton("dev_status", builder.Button{
				ButtonName: "关闭开发模式",
				Action:     "dev_status",
				ActionType: 2,
				LayerTitle: "正在关闭开发模式",
				ActionUrl:  "dev_status?status=0",
				Class:      "black",
				Icon:       "ri-bug-line",
				Display:    "",
				Expand: map[string]string{
					"w": "50%",
					"h": "50%",
				},
			})
		} else {
			pageBuilder.SetButton("dev_status", builder.Button{
				ButtonName: "开启开发模式",
				Action:     "dev_status",
				ActionType: 2,
				LayerTitle: "正在打开开发模式",
				ActionUrl:  "dev_status?status=1",
				Class:      "black",
				Icon:       "ri-bug-line",
				Display:    "",
				Expand: map[string]string{
					"w": "50%",
					"h": "50%",
				},
			})
		}
		pageBuilder.SetButton("reset_server", builder.Button{
			ButtonName: "重启",
			Action:     "reset_server",
			ActionType: 2,
			LayerTitle: "重启应用实例",
			ActionUrl:  "reset_server",
			Class:      "layui-btn-danger",
			Icon:       "ri-restart-fill",
			Display:    "",
			Expand: map[string]string{
				"w": "50%",
				"h": "50%",
			},
		})
		pageBuilder.SetListTopBtns("edit2", "dev_status", "reset_server")
	}

	return nil, 0
}

// NodeListCondition 修改查询条件
func (that AppConfigs) NodeListCondition(pageBuilder *builder.PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {
	//追加查询条件
	condition = append(condition, []interface{}{
		"group_id", that.GroupId,
	})
	return condition, nil, 0
}

// NodeListData 重写列表数据
func (that AppConfigs) NodeListData(pageBuilder *builder.PageBuilder, data []gorose.Data) ([]gorose.Data, error, int) {
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
				"key":   config["name"].(string),
				"value": value,
			})
		}
	}
	return result, nil, 0
}

func (that *AppConfigs) Edit2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.NodeBegin(that.PageBuilder)
	that.PageBuilder.SetActionName("修改")
	appConfigs := Models.Configs{}.GroupConfigs(that.GroupId)
	if len(appConfigs) > 0 {
		for _, config := range appConfigs {
			expand := map[string]interface{}{}
			if config["if"].(string) != "" {
				expand["if"] = config["if"]
			}
			that.PageBuilder.FormFieldsAdd(config["name"].(string), config["field_type"].(string),
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
					"group_id":  that.GroupId,
					"name":      keyName,
					"value":     value,
					"is_inside": 0,
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
						"is_inside":   0,
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
	that.PageBuilder.SetFormData(odata)

	that.ActShow(w, builder.Displayer{
		TplName: "edit.html",
	}, that.PageBuilder)
}

func (that *AppConfigs) DevStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	status := util.GetValue(r, "status")
	first, err := db.New().Table("tb_admin_rules").
		Where("is_delete", 0).
		Where("route", "#dev").
		Order("index_num desc").First()
	if err != nil {
		logger.Error(err.Error())
		that.ErrResult(w, r, 500, "查询数据库出错！", nil)
		return
	}
	if first == nil {
		that.ErrResult(w, r, 404, "查询不到开发信息！", nil)
		return
	}
	if ps.ByName("account_id") != "1" {
		that.ErrResult(w, r, 110, "操作受限！", nil)
		return
	}
	if status == "1" {
		//开启开发模式
		_, err = db.New().Table("tb_admin_rules").
			Where("is_delete", 0).
			Where("route", "#dev").
			Update(map[string]interface{}{
				"status":      1,
				"update_time": util.TimeNow(),
			})
		if err != nil {
			that.ErrResult(w, r, 500, "系统出错了！", nil)
			return
		}
		fmt.Fprint(w, "打开开发菜单成功！\n")

	} else {
		//关闭开发模式
		_, err = db.New().Table("tb_admin_rules").
			Where("is_delete", 0).
			Where("route", "#dev").
			Update(map[string]interface{}{
				"status":      0,
				"update_time": util.TimeNow(),
			})
		if err != nil {
			that.ErrResult(w, r, 500, "系统出错了！", nil)
			return
		}
		fmt.Fprint(w, "关闭开发菜单成功！\n")
	}
}

func (that *AppConfigs) ResetServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tpl := builder.Displayer{
		TplName: "reset_server.html",
	}
	that.ActShow(w, tpl, that.PageBuilder)
}
