/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: FmtRuleFuncs
 * @Version: 1.0.0
 * @Date: 2022/3/31 3:13 下午
 */

package EasyCurd

import (
	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/util"
)
// JsonDecode 字符串转json
func JsonDecode(json interface{}) interface{} {
	var data interface{}
	util.JsonDecode(json.(string), &data)
	return data
}

// Interface2String 任意值转字符串
func Interface2String(data interface{}) interface{} {
	return util.Interface2String(data)
}


func FindSuccess(data gorose.Data, easyCurd EasyCurd) gorose.Data {
	//查询字段信息
	easyModelFields, err := db.New().Table("tb_easy_curd_models_fields").
		Where("model_id", easyCurd.DbModel.ID).
		Where("is_delete", 0).
		Get()
	if err != nil {
		logger.Error(err.Error())
		return data
	}
	//遍历所有字段信息
	for _, emf := range easyModelFields {
		//字段关联了选项集
		if emf["option_models_key"].(string) != "" && emf["is_private"].(int64) == 0 {
			//查询选项集
			optionModels := Models.OptionModels{}.ByKey(emf["option_models_key"], false)
			//遍历数据列表，每一条数据都加上一个新的字段，字段值
			//添加一个需要转换的键（原键 + _name）
			data[emf["field_key"].(string)+"_name"] = "-"
			//字段的值
			fieldValue := data[emf["field_key"].(string)]
			//根据值对比出来选项集的名称
			for _, optionModel := range optionModels {
				if util.Interface2String(optionModel["value"]) == util.Interface2String(fieldValue) {
					data[emf["field_key"].(string)+"_name"] = optionModel["name"]
					break
				}
			}
		}
	}
	return data
}

func SelectSuccess(data []gorose.Data, easyCurd EasyCurd) []gorose.Data {
	//查询字段信息
	easyModelFields, err := db.New().Table("tb_easy_curd_models_fields").
		Where("model_id", easyCurd.DbModel.ID).
		Where("is_delete", 0).
		Get()
	if err != nil {
		logger.Error(err.Error())
		return data
	}
	//遍历所有字段信息
	for _, emf := range easyModelFields {
		//字段关联了选项集
		if emf["option_models_key"].(string) != "" && emf["is_private"].(int64) == 0 {
			//查询选项集
			optionModels := Models.OptionModels{}.ByKey(emf["option_models_key"], false)
			//遍历数据列表，每一条数据都加上一个新的字段，字段值
			for index, item := range data {
				//添加一个需要转换的键（原键 + _name）
				data[index][emf["field_key"].(string)+"_name"] = "-"
				//字段的值
				fieldValue := item[emf["field_key"].(string)]
				//根据值对比出来选项集的名称
				for _, optionModel := range optionModels {
					if util.Interface2String(optionModel["value"]) == util.Interface2String(fieldValue) {
						data[index][emf["field_key"].(string)+"_name"] = optionModel["name"]
						break
					}
				}
			}
		}
	}
	return data
}
