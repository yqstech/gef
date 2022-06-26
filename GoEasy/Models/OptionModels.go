/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: OptionModels
 * @Version: 1.0.0
 * @Date: 2022/5/14 15:06
 */

package Models

import (
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
	"strings"
	"sync"
	"time"
)

type OptionModels struct {
}

var OptionModelsList = map[string][]map[string]interface{}{}
var OptionModelsListLock sync.Mutex

// Select 新增选项集底层查询方法，支持where条件查询
func (that OptionModels) Select(id int, where string, beautify bool) []map[string]interface{} {
	cacheKey := util.Int2String(id) + "_" + util.Is(beautify, "1", "0").(string)
	if where != "" {
		cacheKey = cacheKey + "_" + util.MD5(where)
	}
	if selectData, ok := OptionModelsList[cacheKey]; ok {
		return selectData
	} else {
		OptionModelsListLock.Lock()
		defer OptionModelsListLock.Unlock()
		//通过数据库查询
		data, err := db.New().Table("tb_option_models").
			Where("is_delete", 0).
			Where("status", 1).
			Where("id", id).First()
		if err != nil {
			logger.Error(err.Error())
			return nil
		}
		if data == nil {
			return nil
		}
		if data["data_type"].(int64) == 0 {
			//解析json格式的静态数据
			if data["static_data"].(string) != "" {
				util.JsonDecode(data["static_data"].(string), &selectData)
			}
		} else {
			//数据表查询支持补充颜色和图标
			colorArray := strings.Split(data["color_array"].(string), ",")
			iconArray := strings.Split(data["icon_array"].(string), ",")
			//!合并传入的查询条件和模型定义的查询条件
			if where != "" && data["select_where"].(string) != "" {
				where = where + " and " + data["select_where"].(string)
			} else if data["select_where"].(string) != "" {
				where = data["select_where"].(string)
			}
			//通过数据表获取列表
			keyTrans := map[string]string{
				data["value_field"].(string): "value",
				data["name_field"].(string):  "name",
			}
			if data["icon_field"].(string) != "" {
				keyTrans[data["icon_field"].(string)] = "icon"
			}
			if data["color_field"].(string) != "" {
				keyTrans[data["color_field"].(string)] = "color"
			}
			//设置了上级字段，获取上级字段
			if data["parent_field"].(string) != "" {
				keyTrans[data["parent_field"].(string)] = "pid"
				//如果value不是id，那么直接查询一下
				if data["value_field"].(string) != "id" {
					keyTrans["id"] = "id"
				}
			}
			arrOptions, err, _ := Model{}.SelectOptionsData(data["table_name"].(string), keyTrans, "", "", where, data["select_order"].(string))
			if err != nil {
				logger.Error(err.Error())
				return nil
			}
			if data["value_field"].(string) == "id" && data["parent_field"].(string) != "" {
				//假如value就是ID，需要拷贝一下
				for index, options := range arrOptions {
					arrOptions[index]["id"] = options["value"]
				}
			}
			colorArrayLen := len(colorArray)
			iconArrayLen := len(iconArray)
			for index, _ := range arrOptions {
				if colorArrayLen > index {
					if c, ok := arrOptions[index]["color"]; !ok || c == "" {
						arrOptions[index]["color"] = colorArray[index]
					}
				}
				if iconArrayLen > index {
					if ic, ok := arrOptions[index]["icon"]; !ok || ic == "" {
						arrOptions[index]["icon"] = iconArray[index]
					}
				}
			}
			selectData = arrOptions
		}
		if beautify {
			//静态数据渲染文本
			for index, option := range selectData {
				optionColor := ""
				if bgColor, ok := option["bgColor"].(string); ok {
					optionColor = ";background:" + bgColor + ";border-color:" + bgColor
					if color, ok2 := option["color"].(string); ok2 {
						optionColor = optionColor + ";color:" + color
					} else {
						//有背景色，默认字体就是白色
						optionColor = optionColor + ";color:#FFFFFF"
					}
				} else {
					//无背景色，有颜色，设置字体颜色
					if color, ok := option["color"].(string); ok {
						optionColor = ";color:" + color + ";border-color:" + color
					}
				}
				optionIcon := ""
				if icon, ok := option["icon"].(string); ok {
					optionIcon = "<i class=\"" + icon + "\"></i>"
				}
				optionName := "<div class=\"option-tag\" style=\"" + optionColor + "\">" + optionIcon + util.Interface2String(
					option["name"]) + "</div>"
				selectData[index]["name"] = optionName
			}
		}
		//!定时删除
		go func() {
			t := time.After(time.Second * 10) //十秒钟后删除
			_, _ = <-t
			OptionModelsListLock.Lock()
			defer OptionModelsListLock.Unlock()
			delete(OptionModelsList, cacheKey)
		}()
		OptionModelsList[cacheKey] = selectData
		return selectData
	}
}

// ById 获取选项集数据（可选择美化数据）
func (that OptionModels) ById(id int, beautify bool) []map[string]interface{} {
	return that.Select(id, "", beautify)
}

//定时更新标记
var matchRuleValidity = false

//原子操作锁
var fieldMatchSelectData sync.Mutex

//字段完全匹配的id索引
var fieldKeyDataId = map[string]int64{}

//字段匹配规则列表
var fieldRules []fieldRule

type fieldRule struct {
	Key       string //匹配关键字
	KeyLength int    //匹配字段位数
	MatchType int    //匹配规则（0开头）
	DataId    int64  //数据
}

// FieldMatchOptionModelsId 字段自动匹配选项集ID
func FieldMatchOptionModelsId(fieldKey string) int64 {
	//原子操作
	fieldMatchSelectData.Lock()
	defer fieldMatchSelectData.Unlock()
	//默认返回0
	selectId := int64(0)
	//需要更新规则
	if !matchRuleValidity {
		//设置当前有效期有效
		matchRuleValidity = true
		//!定时十秒钟以后需要再次更新
		go func() {
			t := time.After(time.Second * 10) //十秒钟后删除
			_, _ = <-t
			fieldMatchSelectData.Lock()
			defer fieldMatchSelectData.Unlock()
			matchRuleValidity = false
		}()
		//清空原始的数据
		fieldKeyDataId = make(map[string]int64)
		fieldRules = []fieldRule{}
		//查询关联数据列表和匹配规则
		data, err := db.New().Table("tb_option_models").
			Where("is_delete", 0).
			Where("status", 1).
			OrderBy("id asc").
			Fields("id,match_fields").
			Get()
		if err != nil {
			logger.Error(err.Error())
			return selectId
		}
		for _, item := range data {
			if item["match_fields"].(string) != "" {
				fields := strings.Split(item["match_fields"].(string), "\n")
				for _, field := range fields {
					idx := strings.IndexAny(field, "*")
					if idx == -1 {
						fieldKeyDataId[field] = item["id"].(int64)
					} else {
						fieldRules = append(fieldRules, fieldRule{
							Key:       field[0:idx],
							KeyLength: idx,
							MatchType: 0,
							DataId:    item["id"].(int64),
						})
					}
				}
			}
		}
	}
	if id, ok := fieldKeyDataId[fieldKey]; ok {
		return id
	}
	fieldKeyLength := len(fieldKey)
	for _, rule := range fieldRules {
		if fieldKeyLength >= rule.KeyLength {
			if fieldKey[0:rule.KeyLength] == rule.Key {
				return rule.DataId
			}
		}
	}
	return selectId
}
