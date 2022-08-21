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
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/util"
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
		//OptionModelsListLock.Lock()
		//defer OptionModelsListLock.Unlock()
		//通过数据库查询
		conn := db.New().Table("tb_option_models")
		data, err := conn.
			Where("is_delete", 0).
			Where("status", 1).
			Where(func() {
				if id == 0 {
					conn.Where("1=1")
				} else {
					conn.Where("id", id)
				}
			}).
			Where(func() {
				if where == "" {
					conn.Where("1=1")
				} else {
					conn.Where(where)
				}
			}).First()
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
			if data["default_data"].(string) != "" {
				util.JsonDecode(data["default_data"].(string), &selectData)
			}
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
			//设置了上级字段，获取上级字段,上级字段统一格式化成pid
			if data["parent_field"].(string) != "" {
				keyTrans[data["parent_field"].(string)] = "pid"
			}
			//如果value不是id，那么也需要查询一下id
			if data["value_field"].(string) != "id" {
				keyTrans["id"] = "id"
			}

			arrOptions, err, _ := Model{}.SelectOptionsData(data["table_name"].(string), keyTrans, "", "", where, data["select_order"].(string))
			if err != nil {
				logger.Error(err.Error())
				return nil
			}
			//假如value就是ID，需要将value字段值 恢复 到id字段上
			if data["value_field"].(string) == "id" {
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
			selectData = append(selectData, arrOptions...)
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
				if icon, ok := option["icon"].(string); ok && icon != "" {
					optionIcon = "<i class=\"" + icon + "\"></i>"
				}
				optionName := "<div class=\"option-tag\" style=\"" + optionColor + "\">" + optionIcon + util.Interface2String(
					option["name"]) + "</div>"
				selectData[index]["name"] = optionName
			}
		}
		//logger.Alert("美化后的数据集", selectData)
		//!设置了禁用
		if data["options_disable"].(int64) == 1 {
			for thisIndex, _ := range selectData {
				selectData[thisIndex]["disabled"] = "disabled"
			}
		}
		//! 合并下级选项集
		if data["children_option_model_key"].(string) != "" {
			nextOptionModels := that.ByKey(data["children_option_model_key"], false)
			for _, nextOption := range nextOptionModels {
				//下级必须存在pid字段
				if pid, ok := nextOption["pid"]; ok {
					//循环当前选项
					for thisIndex, thisOption := range selectData {
						//获取值
						if thisId, ok2 := thisOption[data["value_field"].(string)]; ok2 {
							//当前选项集选项ID 和 下级选项集选项的PID相同
							if util.Interface2String(thisId) == util.Interface2String(pid) {
								if _, ok3 := selectData[thisIndex]["_child"]; !ok3 {
									selectData[thisIndex]["_child"] = []map[string]interface{}{}
								}
								selectData[thisIndex]["_child"] = append(selectData[thisIndex]["_child"].([]map[string]interface{}), nextOption)
							}
						}
					}
				}
			}
			//logger.Alert("合并下级以后", selectData)
			//! 跨表会造成value值重复，需要修改一下value值，且需要在补充完下级以后修改
			//!取数据表名最后一个字段作为新ID的前缀
			tbName := strings.Split(data["table_name"].(string), "_")
			Prefix := tbName[len(tbName)-1]
			for thisIndex, thisOption := range selectData {
				selectData[thisIndex]["value"] = Prefix + "_" + util.Interface2String(thisOption["value"])
			}
			//logger.Alert("修改value值以后", selectData)
		}
		//!转多维数组
		if data["to_tree_array"].(int64) == 1 {
			selectData, _, _ = util.ArrayMap2Tree(selectData, 0, data["value_field"].(string), "pid", "_child")
		}
		//!迭代更新一下标记值
		selectData, _, _ = TreeArrayExtendField(selectData)
		//
		//logger.Info(id, util.JsonEncode(selectData))

		//!定时删除
		go func() {
			t := time.After(time.Second * 10) //十秒钟后删除
			_, _ = <-t
			OptionModelsListLock.Lock()
			//defer OptionModelsListLock.Unlock()
			//delete(OptionModelsList, cacheKey)
		}()
		OptionModelsList[cacheKey] = selectData
		return selectData
	}
}

// ById 获取选项集数据（可选择美化数据）
func (that OptionModels) ById(id int, beautify bool) []map[string]interface{} {
	return that.Select(id, "", beautify)
}

// ByKey 获取选项集数据（可选择美化数据）
func (that OptionModels) ByKey(uniqueKey interface{}, beautify bool) []map[string]interface{} {
	var data []map[string]interface{}
	switch uniqueKey.(type) {
	case int:
		data = that.ById(uniqueKey.(int), beautify)
	case int64:
		data = that.ById(util.Int642Int(uniqueKey.(int64)), beautify)
	case string:
		data = that.Select(0, "unique_key = '"+uniqueKey.(string)+"'", beautify)
	default:
		logger.Error("选项集key值类型异常！")
	}
	return data
}

// OptionDynamicParam 选项动态参数
type OptionDynamicParam struct {
	ParamKey string
	FieldKey string
	DefValue string
}

// DynamicParams 获取选项集的动态参数
func (that OptionModels) DynamicParams(uniqueKey string) []OptionDynamicParam {
	var DynamicParams []OptionDynamicParam
	data, err := db.New().Table("tb_option_models").
		Where("is_delete", 0).
		Where("status", 1).
		Where("unique_key", uniqueKey).First()
	if err != nil {
		logger.Error(err.Error())
		return DynamicParams
	}
	if data == nil {
		return DynamicParams
	}
	dp := data["dynamic_params"].(string)
	Params := strings.Split(dp, "\n")
	for _, Param := range Params {
		ps := strings.Split(Param, ":")
		if len(ps) == 3 {
			DynamicParams = append(DynamicParams, OptionDynamicParam{
				ParamKey: ps[0],
				FieldKey: ps[1],
				DefValue: ps[2],
			})
		} else if len(ps) == 2 {
			DynamicParams = append(DynamicParams, OptionDynamicParam{
				ParamKey: ps[0],
				FieldKey: ps[1],
				DefValue: "",
			})
		}
	}
	return DynamicParams
}

//定时更新标记
var matchRuleValidity = false

//原子操作锁
var fieldMatchSelectData sync.Mutex

//字段完全匹配的id索引
var fieldKey2OptionModelKey = map[string]string{}

//字段匹配规则列表
var fieldRules []fieldRule

type fieldRule struct {
	Key       string //匹配关键字
	KeyLength int    //匹配字段位数
	MatchType int    //匹配规则（0开头）
	DataId    int64  //模型ID
	UniqueKey string //模型关键字
}

// FieldMatchOptionModelsKey 字段自动匹配选项集key
func FieldMatchOptionModelsKey(fieldKey string) string {
	//原子操作
	fieldMatchSelectData.Lock()
	defer fieldMatchSelectData.Unlock()
	//默认返回0
	uniqueKey := ""
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
		fieldKey2OptionModelKey = make(map[string]string)
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
			return uniqueKey
		}
		for _, item := range data {
			if item["match_fields"].(string) != "" {
				fields := strings.Split(item["match_fields"].(string), "\n")
				for _, field := range fields {
					idx := strings.IndexAny(field, "*")
					if idx == -1 {
						fieldKey2OptionModelKey[field] = item["unique_key"].(string)
					} else {
						fieldRules = append(fieldRules, fieldRule{
							Key:       field[0:idx],
							KeyLength: idx,
							MatchType: 0,
							DataId:    item["id"].(int64),
							UniqueKey: item["unique_key"].(string),
						})
					}
				}
			}
		}
	}
	if key, ok := fieldKey2OptionModelKey[fieldKey]; ok {
		return key
	}
	fieldKeyLength := len(fieldKey)
	for _, rule := range fieldRules {
		if fieldKeyLength >= rule.KeyLength {
			if fieldKey[0:rule.KeyLength] == rule.Key {
				return rule.UniqueKey
			}
		}
	}
	return uniqueKey
}

// TreeArrayExtendField 多维数据补充数组的附加信息
//多维数组查询出子级的层级，遇到选项有_lastLevel值，就直接返回
func TreeArrayExtendField(data []map[string]interface{}) ([]map[string]interface{}, []interface{}, int64) {
	//迭代数组
	//最多有几级子项
	childMaxLevel := int64(0)
	var childrenIds []interface{}
	//循环数组
	for k, v := range data {
		//明确判断值是数字还是字符串
		valueStr := util.Interface2String(v["value"])
		if util.IsNum(valueStr) {
			childrenIds = append(childrenIds, int64(util.String2Int(valueStr)))
		} else {
			childrenIds = append(childrenIds, valueStr)
		}
		//逐项判断，是否含有下级
		if _child, ok := v["_child"]; ok {
			if len(_child.([]map[string]interface{})) > 0 {
				//如果当前项，已经计算出了下级信息，就不用继续去处理了
				if _children, ok2 := v["_children"]; ok2 {
					childrenIds = append(childrenIds, _children.([]interface{})...)
					if _lastLevel, ok2 := v["_lastLevel"]; ok2 {
						if _lastLevel.(int64)+1 > childMaxLevel {
							childMaxLevel = _lastLevel.(int64) + 1
						}
					}
					continue
				}
				//下级返回信息
				newChild, childrenIds2, childMaxLevel2 := TreeArrayExtendField(_child.([]map[string]interface{}))
				//判断子项有几层下级，如果就子集自己，则返回0
				if childMaxLevel2+1 > childMaxLevel {
					childMaxLevel = childMaxLevel2 + 1
				}
				data[k]["_child"] = newChild
				data[k]["_children"] = childrenIds2
				data[k]["_lastLevel"] = childMaxLevel2 + 1 //后面还有几级

				childrenIds = append(childrenIds, childrenIds2...)

				//处理完毕，跳转到下次循环
				continue
			} else {
				//下级列表为空
			}
		} else {
			data[k]["_child"] = []map[string]interface{}{}
		}

		//获取下级失败或者下级为空
		data[k]["_children"] = []interface{}{}
		data[k]["_lastLevel"] = int64(0)
	}
	return data, childrenIds, childMaxLevel
}
