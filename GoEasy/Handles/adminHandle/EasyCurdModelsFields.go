/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyCurdModelsFields
 * @Version: 1.0.0
 * @Date: 2022/5/3 21:26
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gef/config"
	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
	"strings"
)

type EasyCurdModelsFields struct {
	Base
}

// NodeBegin 开始
func (that EasyCurdModelsFields) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("easyCurd接口模型字段管理")
	pageData.SetPageName("模型字段")
	pageData.SetTbName("tb_easy_curd_models_fields")
	if pageData.GetHttpRequest().Method == "GET" {
		//同步数据库字段
		id := util.GetValue(pageData.GetHttpRequest(), "id")
		that.syncModelFields(util.String2Int(id))
	}
	return nil, 0
}

// NodeList 初始化列表
func (that EasyCurdModelsFields) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.SetListOrder("id asc")
	pageData.SetListTopBtns()
	pageData.ListColumnAdd("field_key", "字段Key", "text", nil)
	pageData.ListColumnAdd("field_name", "字段名称", "text", nil)
	pageData.ListColumnAdd("field_note", "字段备注", "text", nil)
	pageData.ListColumnAdd("option_models_id", "关联选项集", "array", that.OptionModelsList())
	pageData.ListColumnAdd("is_private", "私密数据", "switch::text=私密|公开", nil)
	pageData.ListColumnAdd("is_lock", "锁定数据", "switch::text=锁定|可改", nil)
	pageData.ListColumnAdd("update_time", "最后同步", "text", nil)
	pageData.SetListRightBtns("edit")
	return nil, 0
}

// NodeListCondition 修改查询条件
func (that EasyCurdModelsFields) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	modelID := 0
	modelId := util.GetValue(pageData.GetHttpRequest(), "id")
	if modelId != "" {
		modelID = util.String2Int(modelId)
		//追加查询条件
		condition = append(condition, []interface{}{
			"model_id", "=", modelID,
		})
	}

	return condition, nil, 0
}

// NodeForm 初始化表单
func (that EasyCurdModelsFields) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	pageData.FormFieldsAdd("option_models_id", "select", "关联选项集", "", "", true, that.OptionModelsList(), "", nil)
	return nil, 0
}

// NodeFormData 表单显示前修改数据
func (that EasyCurdModelsFields) NodeFormData(pageData *EasyApp.PageData, data gorose.Data, id int64) (gorose.Data, error, int) {
	if data["option_models_id"].(int64) == 0 {
		data["option_models_id"] = ""
	}
	return data, nil, 0
}

// NodeSaveData 表单保存数据前使用
func (that EasyCurdModelsFields) NodeSaveData(pageData *EasyApp.PageData, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	if postData["option_models_id"] == "" {
		postData["option_models_id"] = 0
	}
	return postData, nil, 0
}

//默认设置成私密的字段名
var defaultPrivateFields = []interface{}{"is_delete", "status", "create_time", "pid"}

//默认锁定(禁止直接修改)的字段
var defaultLockFields = []interface{}{"is_delete", "update_time", "create_time", "pid", "id"}

//同步模型字段
func (that EasyCurdModelsFields) syncModelFields(easyModelId int) {
	//查询模型信息
	easyModelInfo, err := db.New().Table("tb_easy_curd_models").
		Where("id", easyModelId).
		Where("is_delete", 0).First()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if easyModelInfo == nil {
		return
	}
	if easyModelInfo["table_name"].(string) == "" {
		return
	}
	//得到模型数据表名称
	tableName := easyModelInfo["table_name"].(string)

	//查询数据表实时字段信息
	query, err := db.New().Query("select COLUMN_NAME,COLUMN_COMMENT,COLUMN_DEFAULT,COLUMN_TYPE from information_schema.COLUMNS where table_name = ? and table_schema = ? order by ordinal_position",
		tableName,
		config.DbName)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//logger.Debug(util.JsonEncode(query))

	//查询模型字段列表
	fields, err := db.New().Table("tb_easy_curd_models_fields").
		Where("model_id", easyModelId).
		Where("is_delete", 0).
		Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//字段列表转成以fieldKey为键的map数据
	fieldsMap := map[string]gorose.Data{}
	for _, field := range fields {
		fieldsMap[field["field_key"].(string)] = field
	}
	timeNow := util.TimeNow()
	//对比更新模型字段信息
	for _, fieldInfo := range query {
		//字段key
		fieldKey := fieldInfo["COLUMN_NAME"].(string)
		//字段名称和备注，直接同步数据库的备注信息
		fieldName := ""
		fieldNote := ""
		commands := strings.Split(fieldInfo["COLUMN_COMMENT"].(string), "|")
		if len(commands) > 0 {
			fieldName = commands[0]
		}
		if len(commands) > 1 {
			fieldNote = commands[1]
		}
		//是否默认设置成私密
		isPrivate := 0
		if util.IsInArray(fieldKey, defaultPrivateFields) {
			isPrivate = 1
		}
		//是否默认设置成锁定
		isLock := 0
		if util.IsInArray(fieldKey, defaultLockFields) {
			isLock = 1
		}
		//是否存在字段
		if field, ok := fieldsMap[fieldKey]; ok {
			db.New().Table("tb_easy_curd_models_fields").
				Where("id", field["id"]).
				Update(map[string]interface{}{
					"field_name":  fieldName,
					"field_note":  fieldNote,
					"update_time": timeNow,
				})
			fieldsMap[fieldKey]["sync_tag"] = true
		} else {
			db.New().Table("tb_easy_curd_models_fields").Insert(map[string]interface{}{
				"model_id":         easyModelId,
				"field_key":        fieldKey,
				"field_name":       fieldName,
				"field_note":       fieldNote,
				"option_models_id": Models.FieldMatchOptionModelsId(fieldKey), //自动匹配选择数据源
				"is_private":       isPrivate,
				"is_lock":          isLock,
				"update_time":      timeNow,
			})
		}
	}

	for _, field := range fieldsMap {
		//未标记的都删除
		if _, ok2 := field["sync_tag"]; !ok2 {
			db.New().Table("tb_easy_curd_models_fields").
				Where("id", field["id"]).
				Update(map[string]interface{}{
					"update_time": timeNow,
					"is_delete":   1,
				})
		}
	}
}
