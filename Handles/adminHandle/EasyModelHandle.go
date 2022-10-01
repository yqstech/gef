/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyModelHandle
 * @Version: 1.0.0
 * @Date: 2022/5/2 10:54
 */

package adminHandle

import (
	"errors"
	"github.com/gohouse/gorose/v2"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"
	"strings"
	"sync"
	"time"
)

type EasyModelHandle struct {
	Base
	ModelKey string //模型关键字
}

// NodeBegin 开始
func (that EasyModelHandle) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, "begin")
	if err != nil {
		return err, 500
	} else {
		pageBuilder.SetTitle(easyModel.ModelName + "管理")
		pageBuilder.SetPageName(easyModel.ModelName)
		pageBuilder.SetTbName(easyModel.TableName)
		return nil, 0
	}
}

// NodeList 初始化列表
func (that EasyModelHandle) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, "list")
	if err != nil {
		return err, 500
	} else {
		//!排序、分页、列表页面提示
		pageBuilder.SetListOrder(easyModel.OrderType)
		if easyModel.PageSize > 0 {
			pageBuilder.SetListPageSize(easyModel.PageSize)
		}
		if easyModel.PageNotice != "" {
			pageBuilder.SetPageNotice(easyModel.PageNotice)
		}
		//!模型url参数透传到tabs
		var validParams []string
		for _, urlParam := range easyModel.UrlParams {
			validParams = append(validParams, urlParam.ParamKey)
		}
		//!设置tabs列表和选中项
		validUrl := util.UrlScreenParam(pageBuilder.GetHttpRequest(), validParams, false, true)
		for index, tab := range easyModel.ListTabs {
			pageBuilder.PageTabAdd(tab.TabName, validUrl+"tab="+util.Int2String(index))
		}
		//获取第几页
		tabIndex := that.GetTabIndex(pageBuilder, "tab")
		pageBuilder.SetPageTabSelect(tabIndex)

		//!获取自定义按钮列表
		for btnName, Btn := range easyModel.Buttons {
			pageBuilder.SetButton(btnName, Btn)
		}

		//!顶部按钮
		topBtns := easyModel.TopButtons
		//新增
		if !easyModel.AllowCreate {
			for index, text := range topBtns {
				if text == "add" {
					topBtns = append(topBtns[0:index], topBtns[index+1:]...)
				}
			}
		}
		pageBuilder.SetListTopBtns(topBtns...)
		//!URL参数透传到顶部按钮上
		for _, urlParam := range easyModel.UrlParams {
			for _, btn := range topBtns {
				urlAppend := urlParam.FieldKey + "=" + util.GetValue(pageBuilder.GetHttpRequest(), urlParam.ParamKey)
				pageBuilder.SetButtonActionUrl(btn, urlAppend, true)
			}
		}
		//!右侧按钮
		rightBtns := easyModel.RightButtons
		//修改
		if !easyModel.AllowUpdate {
			for index, text := range rightBtns {
				if text == "edit" {
					rightBtns = append(rightBtns[0:index], rightBtns[index+1:]...)
				}
			}
		}
		//状态
		if !easyModel.AllowStatus {
			for index, text := range rightBtns {
				if text == "disable" || text == "enable" {
					rightBtns = append(rightBtns[0:index], rightBtns[index+1:]...)
				}
			}
		}
		//删除
		if !easyModel.AllowDelete {
			for index, text := range rightBtns {
				if text == "delete" {
					rightBtns = append(rightBtns[0:index], rightBtns[index+1:]...)
				}
			}
		}
		pageBuilder.SetListRightBtns(rightBtns...)
		//ID列同步到字段管理了，统一不显示默认id列了
		pageBuilder.ListColumnClear()
		if easyModel.BatchAction {
			pageBuilder.SetListBatchAction(true)
		}
		//增加列表列
		for _, field := range easyModel.Fields {
			if field.IsShowOnList {
				//自定义列名称
				if field.FieldNameReset != "" {
					field.FieldName = field.FieldNameReset
				}
				pageBuilder.ListColumnAdd(field.FieldKey, field.FieldName, field.DataTypeOnList+"::"+field.DataTypeCommandOnList, field.FieldOptions)
				//自定义列样式
				if field.FieldStyleReset != "" {
					pageBuilder.SetListColumnStyle(field.FieldKey, field.FieldStyleReset)
				}
			}
		}

		//当前选中的tab页面
		TabSelectParams := map[string]interface{}{}
		if tabIndex >= 0 && tabIndex+1 <= len(easyModel.ListTabs) {
			TabSelectParams = easyModel.ListTabs[tabIndex].SearchFormParams
		}

		//!搜索表单
		for _, searchItem := range easyModel.SearchForm {
			//选项集
			var ItemOptionModels []map[string]interface{}
			if searchItem.OptionModelsKey != "" {
				//查询当前选项集是否支持联动
				//# 配置了联动的监听参数
				//# 如果参数1和tab参数一致
				//# 则查询条件为 参数2=tab参数值
				var Sqls []string
				DynamicParams := Models.OptionModels{}.DynamicParams(searchItem.OptionModelsKey)
				for _, DynamicParam := range DynamicParams {
					if SelectParamValue, ok := TabSelectParams[DynamicParam.ParamKey]; ok {
						Sqls = append(Sqls, DynamicParam.FieldKey+"="+SelectParamValue.(string))
					}
				}
				ItemOptionModels = Models.OptionModels{}.ByKeySelect(searchItem.OptionModelsKey, strings.Join(Sqls, " and "), false)
			}
			//自定义追加选项集
			for _, oma := range searchItem.OptionModelsAdd {
				ItemOptionModels = append(ItemOptionModels, oma)
			}
			pageBuilder.ListSearchFieldAdd(
				searchItem.SearchKey,
				searchItem.DataType,
				searchItem.SearchName,
				"",
				searchItem.DefaultValue,
				ItemOptionModels,
				searchItem.Style,
				map[string]interface{}{
					"placeholder": searchItem.Placeholder,
				})
		}
		return nil, 0
	}
}

// NodeListCondition 修改查询条件
func (that EasyModelHandle) NodeListCondition(pageBuilder *builder.PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {
	//!删除自动添加的搜索项，由EasyModel搜索表单配置信息统一处理
	pageBuilder.SetListCondition([][]interface{}{})
	condition = [][]interface{}{}

	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, "list")
	if err != nil {
		return condition, nil, 0
	} else {
		//设置url参数查询条件
		for _, urlParam := range easyModel.UrlParams {
			if urlParam.ParamKey != "" && urlParam.FieldKey != "" {
				urlValue := util.GetValue(pageBuilder.GetHttpRequest(), urlParam.ParamKey)
				if urlValue == "" && urlParam.DefaultValue != "" {
					urlValue = urlParam.DefaultValue
				}
				if urlValue != "" {
					//追加查询条件
					condition = append(condition, []interface{}{
						urlParam.FieldKey, "=", urlValue,
					})
				}
			}
		}
		//设置tab查询条件
		tabIndex := that.GetTabIndex(pageBuilder, "tab")
		if len(easyModel.ListTabs) > tabIndex {
			if easyModel.ListTabs[tabIndex].SelectCondition != "" {
				condition = append(condition, []interface{}{
					easyModel.ListTabs[tabIndex].SelectCondition,
				})
			}
		}
		//全部查询条件
		PostSearch := map[string]interface{}{}
		postSearch := pageBuilder.HttpRequest.PostFormValue("search")
		if postSearch != "" {
			util.JsonDecode(postSearch, &PostSearch)
		}

		//? 设置表单查询条件
		for _, searchFormItem := range easyModel.SearchForm {
			//!仅支持查询单个字段
			if len(searchFormItem.SearchFields) == 1 {
				if postValue, ok := PostSearch[searchFormItem.SearchKey]; ok {
					if postValue != "" {
						//是否使用自定义sql
						useSelfSql := false
						for _, v := range searchFormItem.OptionModelsAdd {
							if postValue == util.Interface2String(v["value"]) {
								if sql, ok := v["sql"]; ok {
									useSelfSql = true
									//设置自定义查询条件
									condition = append(condition, []interface{}{sql})
									break
								}
							}
						}
						if useSelfSql {
							continue
						}

						postValueTs := util.Interface2String(postValue)
						if searchFormItem.SubQuery != "" && searchFormItem.MatchType == "like" {
							//子查询+like
							postValueTs = "CONCAT('%', (" + strings.Replace(searchFormItem.SubQuery, "$1", util.Interface2String(postValue), -1) + "), '%')"
						} else if searchFormItem.SubQuery != "" {
							//子查询+普通运算符
							postValueTs = "(" + strings.Replace(searchFormItem.SubQuery, "$1", util.Interface2String(postValue), -1) + ")"
						} else if searchFormItem.MatchType == "like" {
							//like模糊查询
							postValueTs = "%" + util.Interface2String(postValue) + "%"
						} else {
							//其他不变
						}
						if searchFormItem.SubQuery != "" {
							condition = append(condition, []interface{}{
								searchFormItem.SearchFields[0] + " " + searchFormItem.MatchType + " " + postValueTs,
							})
						} else {
							condition = append(condition, []interface{}{
								searchFormItem.SearchFields[0], searchFormItem.MatchType, postValueTs,
							})
						}

					}
				}
			}
		}
	}
	//追加查询条件
	return condition, nil, 0
}

// NodeListOrm 追加高级查询条件
func (that EasyModelHandle) NodeListOrm(pageBuilder *builder.PageBuilder, conn *gorose.IOrm) (error, int) {
	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, "list")
	if err != nil {
		return err, 0
	} else {
		//全部查询条件
		PostSearch := map[string]interface{}{}
		postSearch := pageBuilder.HttpRequest.PostFormValue("search")
		if postSearch != "" {
			util.JsonDecode(postSearch, &PostSearch)
		}
		//!多个字段同时查询
		for _, searchFormItem := range easyModel.SearchForm {
			//当前查询关联了多个字段 使用or拼接
			if len(searchFormItem.SearchFields) > 1 {
				if postValue, ok := PostSearch[searchFormItem.SearchKey]; ok {
					if postValue != "" {
						//判断值是否是自定义的追加选项，且设置了追加选项自定义的sql查询语句
						SelfSql := ""
						for _, v := range searchFormItem.OptionModelsAdd {
							if postValue == util.Interface2String(v["value"]) {
								if sql, ok := v["sql"]; ok {
									SelfSql = sql.(string)
								}
							}
						}

						//提交的值转为sql语句
						postValueTs := ""
						if searchFormItem.SubQuery != "" && searchFormItem.MatchType == "like" {
							//子查询+like模糊查询的方式
							postValueTs = "CONCAT('%', (" + strings.Replace(searchFormItem.SubQuery, "$1", util.Interface2String(postValue), -1) + "), '%')"
						} else if searchFormItem.SubQuery != "" {
							//子查询+普通运算符
							postValueTs = "(" + strings.Replace(searchFormItem.SubQuery, "$1", util.Interface2String(postValue), -1) + ")"
						} else if searchFormItem.MatchType == "like" {
							//like模糊查询
							postValueTs = "%" + util.Interface2String(postValue) + "%"
						} else {
							//其他不变
							postValueTs = util.Interface2String(postValue)
						}
						logger.Error(postValueTs + "|||")

						//! 不可直接引用，gorose可能是异步处理，searchFormItem会改变
						SearchFields := searchFormItem.SearchFields
						SubQuery := searchFormItem.SubQuery
						MatchType := searchFormItem.MatchType
						(*conn).Where(func() {
							for index, field := range SearchFields {
								if index == 0 {
									if SelfSql != "" {
										(*conn).Where(SelfSql)
									} else if SubQuery != "" {
										(*conn).Where(field + " " + MatchType + " " + postValueTs)
									} else {
										(*conn).Where(field, MatchType, postValueTs)
									}
								} else {
									if SelfSql != "" {
										(*conn).OrWhere(SelfSql)
									} else if SubQuery != "" {
										(*conn).OrWhere(field + " " + MatchType + " " + postValueTs)
									} else {
										(*conn).OrWhere(field, MatchType, postValueTs)
									}
								}
							}
						})
					}
				}
			}
		}
	}
	return nil, 200
}

// NodeListData 重写列表数据
func (that EasyModelHandle) NodeListData(pageBuilder *builder.PageBuilder, data []gorose.Data) ([]gorose.Data, error, int) {
	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, "list")
	if err != nil {
		return data, nil, 0
	} else {
		//按级缩进
		if easyModel.LevelIndent != "" {
			LevelIndent := strings.Split(easyModel.LevelIndent, ":")
			if len(LevelIndent) == 2 {
				data = Models.Model{}.GoroseDataLevelOrder(data, "id", LevelIndent[0], 0, 0)
				for k, v := range data {
					if v["level"] == int64(0) {
						data[k][LevelIndent[1]] = "&nbsp;&nbsp;" + v[LevelIndent[1]].(string)
					} else if v["level"] == int64(1) {
						data[k][LevelIndent[1]] = "&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v[LevelIndent[1]].(string)
					} else if v["level"] == int64(2) {
						data[k][LevelIndent[1]] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v[LevelIndent[1]].(string)
					} else if v["level"] == int64(3) {
						data[k][LevelIndent[1]] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v[LevelIndent[1]].(string)
					} else if v["level"] == int64(4) {
						data[k][LevelIndent[1]] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v[LevelIndent[1]].(string)
					}
				}
			}
		}
		//存储数据还原
		//遍历所有键
		for _, field := range easyModel.Fields {
			//当前键需要元和分互转
			if field.SaveTransRule == "yuan2fen" {
				//循环所有数据
				for k, v := range data {
					//转换同键的数据
					data[k][field.FieldKey] = util.Money(v[field.FieldKey].(int64))
				}
			}
			//秒和小时互转
			if field.SaveTransRule == "hour2second" {
				//循环所有数据
				for k, v := range data {
					//转换同键的数据
					data[k][field.FieldKey] = util.Float642String(float64(v[field.FieldKey].(int64)) / 3600.0)
				}
			}
			//分钟和秒互转
			if field.SaveTransRule == "minute2second" {
				//循环所有数据
				for k, v := range data {
					//转换同键的数据
					data[k][field.FieldKey] = util.Float642String(float64(v[field.FieldKey].(int64)) / 60.0)
				}
			}
			//秒转天时分秒
			if field.SaveTransRule == "dhms2second" {
				//循环所有数据
				for k, v := range data {
					//转换同键的数据
					data[k][field.FieldKey] = util.Second2Dhms(v[field.FieldKey].(int64), "天", "时", "分", "秒")
				}
			}
		}
		//数据增强
		for _, field := range easyModel.Fields {
			if field.FieldAugment != "" {
				for k, v := range data {
					data[k][field.FieldKey] = strings.Replace(field.FieldAugment, "{{this}}", util.Interface2String(v[field.FieldKey]), -1)
				}
			}
		}
		//数据合并
		for _, field := range easyModel.Fields {
			if field.Attach2Field != "" {
				for k, v := range data {
					if _, ok := v[field.Attach2Field]; ok {
						data[k][field.Attach2Field] = util.Interface2String(v[field.Attach2Field]) + "<br>" + util.Interface2String(v[field.FieldKey])
					}
				}
			}
		}
		return data, nil, 0
	}

}

// NodeForm 初始化表单
func (that EasyModelHandle) NodeForm(pageBuilder *builder.PageBuilder, id int64) (error, int) {
	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, util.Is(id == 0, "add", "edit").(string))
	if err != nil {
		return err, 500
	} else {
		if id == 0 {
			//新增页
			for _, field := range easyModel.Fields {
				if field.AllowCreate {
					//格式化拓展数据
					expand := map[string]interface{}{}
					if field.ExpandIf != "" {
						expand["if"] = field.ExpandIf
					}
					if field.ExpandWatchFields != "" && field.ExpandDynamicOptionModelsKey != "" {
						expand["watch_fields"] = field.ExpandWatchFields
						expand["dynamic_option_model_key"] = field.ExpandDynamicOptionModelsKey
					}
					if field.GroupTitle != "" {
						pageBuilder.FormFieldsAdd("", "block", field.GroupTitle, "", "", false, nil, "", nil)
					}
					//选项缩进
					var FieldOptions []map[string]interface{}
					//深度拷贝，直接复制的话，多次刷新会造成多次缩进
					for _, v := range field.FieldOptions {
						x := map[string]interface{}{}
						for k1, v1 := range v {
							x[k1] = v1
						}
						FieldOptions = append(FieldOptions, x)
					}
					if field.OptionIndent {
						//转一下类型
						var FieldOptionsCopy []gorose.Data
						for _, v := range FieldOptions {
							v["id"] = int64(util.String2Int(util.Interface2String(v["id"])))
							v["pid"] = int64(util.String2Int(util.Interface2String(v["pid"])))
							FieldOptionsCopy = append(FieldOptionsCopy, v)
						}
						//执行缩进
						FieldOptionsCopy = Models.Model{}.GoroseDataLevelOrder(FieldOptionsCopy, "id", "pid", 0, 0)
						for k, v := range FieldOptionsCopy {
							if v["level"] == int64(0) {
							} else if v["level"] == int64(1) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							} else if v["level"] == int64(2) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							} else if v["level"] == int64(3) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							} else if v["level"] == int64(4) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							}
						}
						//类型再回来
						FieldOptions = []map[string]interface{}{}
						for _, v := range FieldOptionsCopy {
							FieldOptions = append(FieldOptions, v)
						}
					}
					//增加表单项
					pageBuilder.FormFieldsAdd(field.FieldKey, field.DataTypeOnCreate, field.FieldName, field.FieldNotice, field.DefaultValue, field.IsMust, FieldOptions, "", expand)
				}
			}
		} else {
			for _, field := range easyModel.Fields {
				if field.AllowUpdate {
					//格式化拓展数据
					expand := map[string]interface{}{}
					if field.ExpandIf != "" {
						expand["if"] = field.ExpandIf
					}
					if field.ExpandWatchFields != "" && field.ExpandDynamicOptionModelsKey != "" {
						expand["watch_fields"] = field.ExpandWatchFields
						expand["dynamic_option_model_key"] = field.ExpandDynamicOptionModelsKey
					}
					if field.GroupTitle != "" {
						pageBuilder.FormFieldsAdd("", "block", field.GroupTitle, "", "", false, nil, "", nil)
					}
					//选项缩进
					var FieldOptions []map[string]interface{}
					//深度拷贝，直接复制的话，多次刷新会造成多次缩进
					for _, v := range field.FieldOptions {
						x := map[string]interface{}{}
						for k1, v1 := range v {
							x[k1] = v1
						}
						FieldOptions = append(FieldOptions, x)
					}
					if field.OptionIndent {
						//转一下类型
						var FieldOptionsCopy []gorose.Data
						for _, v := range FieldOptions {
							v["id"] = int64(util.String2Int(util.Interface2String(v["id"])))
							v["pid"] = int64(util.String2Int(util.Interface2String(v["pid"])))
							FieldOptionsCopy = append(FieldOptionsCopy, v)
						}
						//执行缩进
						FieldOptionsCopy = Models.Model{}.GoroseDataLevelOrder(FieldOptionsCopy, "id", "pid", 0, 0)
						for k, v := range FieldOptionsCopy {
							if v["level"] == int64(0) {
							} else if v["level"] == int64(1) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							} else if v["level"] == int64(2) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							} else if v["level"] == int64(3) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							} else if v["level"] == int64(4) {
								FieldOptionsCopy[k]["name"] = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├&nbsp;&nbsp;&nbsp;&nbsp;├─&nbsp;" + v["name"].(string)
							}
						}
						//类型再回来
						FieldOptions = []map[string]interface{}{}
						for _, v := range FieldOptionsCopy {
							FieldOptions = append(FieldOptions, v)
						}
					}
					//增加表单项
					pageBuilder.FormFieldsAdd(field.FieldKey, field.DataTypeOnUpdate, field.FieldName, field.FieldNotice, field.DefaultValue, field.IsMust, FieldOptions, "", expand)
				}
			}
		}

		return nil, 0
	}
}

// NodeFormData 表单显示前修改数据
func (that EasyModelHandle) NodeFormData(pageBuilder *builder.PageBuilder, data gorose.Data, id int64) (gorose.Data, error, int) {
	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, util.Is(id == 0, "add", "edit").(string))
	if err != nil {
		return data, nil, 0
	} else {
		if id > 0 {
			//修改页
			for _, field := range easyModel.Fields {
				//当前键需要元和分互转
				if field.SaveTransRule == "yuan2fen" {
					data[field.FieldKey] = util.Money(data[field.FieldKey].(int64))
				}
				if field.SaveTransRule == "hour2second" {
					data[field.FieldKey] = util.Float642String(float64(data[field.FieldKey].(int64)) / 3600.0)
				}
				if field.SaveTransRule == "minute2second" {
					data[field.FieldKey] = util.Float642String(float64(data[field.FieldKey].(int64)) / 60.0)
				}
				if field.SaveTransRule == "dhms2second" {
					data[field.FieldKey] = util.Second2Dhms(data[field.FieldKey].(int64), "天", "时", "分", "秒")
				}

			}
		} else {
			//新增页链接参数透传
			for _, urlParam := range easyModel.UrlParams {
				field := urlParam.FieldKey
				value := util.GetValue(pageBuilder.GetHttpRequest(), field)
				data[field] = value
			}
		}
	}
	return data, nil, 0
}

// NodeSaveData 表单保存数据前使用
func (that EasyModelHandle) NodeSaveData(pageBuilder *builder.PageBuilder, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	easyModel, err := GetEasyModelInfo(pageBuilder, that.ModelKey, util.Is(oldData["id"].(int64) == 0, "add", "edit").(string))
	if err != nil {
		return postData, nil, 0
	} else {
		for _, field := range easyModel.Fields {
			//当前键需要元和分互转
			if field.SaveTransRule == "yuan2fen" {
				v, err := util.Money2Cent(postData[field.FieldKey].(string))
				if err != nil {
					return nil, errors.New("金额格式错误"), 500
				}
				postData[field.FieldKey] = v
			}
			if field.SaveTransRule == "hour2second" {
				v, err := util.FloatString2Int64(postData[field.FieldKey].(string), 3600.0)
				if err != nil {
					return nil, errors.New("金额格式错误"), 500
				}
				postData[field.FieldKey] = v
			}
			if field.SaveTransRule == "minute2second" {
				v, err := util.FloatString2Int64(postData[field.FieldKey].(string), 60.0)
				if err != nil {
					return nil, errors.New("金额格式错误"), 500
				}
				postData[field.FieldKey] = v
			}
			if field.SaveTransRule == "dhms2second" {
				postData[field.FieldKey] = util.Dhms2Second(postData[field.FieldKey].(string))
			}
		}
	}
	return postData, nil, 0
}

// EasyModel 定义模型信息结构体
type EasyModel struct {
	ModelKey     string                    //模型关键字
	ModelName    string                    //模型名称
	TableName    string                    //数据表名称
	AllowCreate  bool                      //允许新增
	AllowUpdate  bool                      //允许新增
	AllowDelete  bool                      //允许删除
	AllowStatus  bool                      //允许改状态
	Fields       []EasyModelField          //字段列表
	ListTabs     []EasyModelListTab        //列表多Tab页
	SearchForm   []SearchFromItem          //搜索表单项
	OrderType    string                    //排序方式
	PageSize     int                       //分页大小
	PageNotice   string                    //页面备注
	Buttons      map[string]builder.Button //页面
	TopButtons   []string                  //顶部按钮
	RightButtons []string                  //右侧操作按钮
	UrlParams    []EasyModelUrlParam       //链接参数
	LevelIndent  string                    //按级缩进
	BatchAction  bool                      //是否支持批量操作
}
type EasyModelListTab struct {
	TabName          string                 //tab页名称
	SelectCondition  string                 //tab页查询条件，根据tab值，解析查询条件
	SearchFormParams map[string]interface{} //参数:值
}
type EasyModelUrlParam struct {
	ParamKey     string //参数key
	FieldKey     string //数据库字段
	DefaultValue string //默认值
}

// EasyModelField 模型字段信息
type EasyModelField struct {
	FieldKey                     string                   //字段关键字
	FieldName                    string                   //字段名称
	FieldNameReset               string                   //字段重置名称（列表）
	FieldStyleReset              string                   //字段重设样式（列表）
	FieldNotice                  string                   //字段提示信息
	IsShowOnList                 bool                     //是否在列表页展示该字段
	DataTypeOnList               string                   //列表页展示时的数据类型
	DataTypeCommandOnList        string                   //列表页展示时数据类型的指令
	DataTypeOnCreate             string                   //表单页展示时的数据类型
	DataTypeOnUpdate             string                   //表单页展示时的数据类型
	AllowCreate                  bool                     //新增页显示
	AllowUpdate                  bool                     //修改页显示
	IsMust                       bool                     //是否必填
	DefaultValue                 string                   //默认值
	FieldOptions                 []map[string]interface{} //待选数据
	OptionIndent                 bool                     //选项按照上下级缩进
	ExpandIf                     string                   //拓展数据-if条件
	ExpandWatchFields            string                   //拓展数据-联动监听字段
	ExpandDynamicOptionModelsKey string                   //拓展数据-联动绑定动态选项集
	GroupTitle                   string                   //创建一个新分组的标题
	FieldAugment                 string                   //数据增强转换规则
	Attach2Field                 string                   //附加到其他字段
	SaveTransRule                string                   //保存时数据变换规则
}

type SearchFromItem struct {
	DataType        string                   //组件类型名称
	SearchKey       string                   //组件关键字
	SearchName      string                   //组件名，提示信息
	Placeholder     string                   //提示信息
	Style           string                   //附加样式
	OptionModelsKey string                   //关联选项集
	OptionModelsAdd []map[string]interface{} //选项集追加
	SearchFields    []string                 //搜索字段集合
	MatchType       string                   //匹配规则
	DefaultValue    string                   //默认值
	SubQuery        string                   //子查询语句
	//DefaultOptionValue string                   //默认项的值
	//EmptySetValue      string                   //空值设置值
}

// easyModelList 存储的模型信息列表
var easyModelList = map[string]EasyModel{}
var easyModelListLock sync.Mutex

// GetEasyModelInfo 获取模型信息
func GetEasyModelInfo(pageBuilder *builder.PageBuilder, modelKey string, actionName string) (EasyModel, error) {
	saveKey := modelKey + "_" + actionName
	if easyModel, ok := easyModelList[saveKey]; ok {
		return easyModel, nil
	} else {
		easyModelListLock.Lock()
		defer easyModelListLock.Unlock()
		//!查询模型信息
		modelInfo, err := db.New().Table("tb_easy_models").
			Where("model_key", modelKey).
			Where("is_delete", 0).
			Where("status", 1).
			First()
		if err != nil {
			logger.Error(err.Error())
			return EasyModel{}, errors.New("系统运行错误！")
		}
		if modelInfo == nil {
			return EasyModel{}, errors.New("模型" + modelKey + "不存在！")
		}
		//!初始化模型
		easyModel := EasyModel{
			ModelKey:     modelInfo["model_key"].(string),
			ModelName:    modelInfo["model_name"].(string),
			TableName:    modelInfo["table_name"].(string),
			AllowCreate:  modelInfo["allow_create"].(int64) == 1,
			AllowUpdate:  modelInfo["allow_update"].(int64) == 1,
			AllowDelete:  modelInfo["allow_delete"].(int64) == 1,
			AllowStatus:  modelInfo["allow_status"].(int64) == 1,
			Fields:       []EasyModelField{},
			ListTabs:     []EasyModelListTab{},
			OrderType:    modelInfo["order_type"].(string),
			PageSize:     util.Int642Int(modelInfo["page_size"].(int64)),
			PageNotice:   modelInfo["page_notice"].(string),
			Buttons:      map[string]builder.Button{},
			TopButtons:   []string{},
			RightButtons: []string{},
			UrlParams:    []EasyModelUrlParam{},
			LevelIndent:  modelInfo["level_indent"].(string),
			BatchAction:  util.Is(modelInfo["batch_action"].(int64) == 1, true, false).(bool),
		}
		//格式化多tab页
		if modelInfo["tabs_for_list"].(string) != "" {
			tabs := strings.Split(modelInfo["tabs_for_list"].(string), "\n")
			for _, tab := range tabs {
				Opts := strings.Split(tab, "|")
				OptsLength := len(Opts)
				newTab := EasyModelListTab{
					TabName:          "标签页",
					SelectCondition:  "",
					SearchFormParams: map[string]interface{}{},
				}
				if OptsLength > 0 && Opts[0] != "" {
					newTab.TabName = Opts[0]
				}
				if OptsLength > 1 && Opts[1] != "" {
					//!对自定义查询条件进行数据库兼容处理
					newTab.SelectCondition = util.Sql(Opts[1])
				}
				//! tab和搜索表单的联动参数设置
				if OptsLength > 2 && Opts[2] != "" {
					searchFormPs := strings.Split(Opts[2], "&")
					for _, ps := range searchFormPs {
						if ps != "" {
							params := strings.Split(ps, "=")
							if len(params) == 2 {
								newTab.SearchFormParams[params[0]] = params[1]
							}
						}
					}
				}
				easyModel.ListTabs = append(easyModel.ListTabs, newTab)
			}
		}
		//!格式化按钮列表
		var allButton []interface{}
		var topButtons []map[string]interface{}
		if modelInfo["top_buttons"].(string) != "" {
			util.JsonDecode(modelInfo["top_buttons"].(string), &topButtons)
		}
		for _, btn := range topButtons {
			easyModel.TopButtons = append(easyModel.TopButtons, btn["text"].(string))
			allButton = append(allButton, btn["text"])
		}
		var rightButtons []map[string]interface{}
		if modelInfo["right_buttons"].(string) != "" {
			util.JsonDecode(modelInfo["right_buttons"].(string), &rightButtons)
		}
		for _, btn := range rightButtons {
			easyModel.RightButtons = append(easyModel.RightButtons, btn["text"].(string))
			allButton = append(allButton, btn["text"])
		}
		//!获取自定义按钮列表
		easyModel.Buttons, err = pageBuilder.GetEasyModelsButtons(allButton)
		if err != nil {
			return EasyModel{}, err
		}
		//!格式化url参数
		urlParams := strings.Split(modelInfo["url_params"].(string), "\n")
		for _, urlParam := range urlParams {
			params := strings.Split(urlParam, ":")
			if len(params) == 3 {
				easyModel.UrlParams = append(easyModel.UrlParams, EasyModelUrlParam{
					ParamKey:     params[0],
					FieldKey:     params[1],
					DefaultValue: params[2],
				})
			} else if len(params) == 2 {
				easyModel.UrlParams = append(easyModel.UrlParams, EasyModelUrlParam{
					ParamKey:     params[0],
					FieldKey:     params[1],
					DefaultValue: "",
				})
			}
		}

		//!模型字段信息
		modelFields, err := db.New().Table("tb_easy_models_fields").
			Where("is_delete", 0).
			Where("status", 1).
			Where("model_id", modelInfo["id"].(int64)).
			Order("index_num asc,id asc").
			Get()
		if err != nil {
			logger.Error(err.Error())
			return EasyModel{}, errors.New("系统运行错误！")
		}
		for _, field := range modelFields {
			modelField := EasyModelField{
				FieldKey:                     field["field_key"].(string),
				FieldName:                    field["field_name"].(string),
				FieldNameReset:               field["field_name_reset"].(string),
				FieldStyleReset:              field["field_style_reset"].(string),
				FieldNotice:                  field["field_notice"].(string),
				IsShowOnList:                 field["is_show_on_list"].(int64) == 1,
				DataTypeOnList:               field["data_type_on_list"].(string),
				DataTypeCommandOnList:        field["data_type_command_on_list"].(string),
				DataTypeOnCreate:             field["data_type_on_create"].(string),
				DataTypeOnUpdate:             field["data_type_on_update"].(string),
				AllowCreate:                  field["allow_create"].(int64) == 1,
				AllowUpdate:                  field["allow_update"].(int64) == 1,
				IsMust:                       field["is_must"].(int64) == 1,
				DefaultValue:                 field["default_value"].(string),
				FieldOptions:                 nil,
				OptionIndent:                 field["option_indent"].(int64) == 1,
				ExpandIf:                     field["expand_if"].(string),
				ExpandWatchFields:            field["watch_fields"].(string),
				ExpandDynamicOptionModelsKey: field["dynamic_option_models_key"].(string),
				GroupTitle:                   field["group_title"].(string),
				FieldAugment:                 field["field_augment"].(string),
				Attach2Field:                 field["attach_to_field"].(string),
				SaveTransRule:                field["save_trans_rule"].(string),
			}
			if field["option_models_key"].(string) != "" {
				if actionName == "list" {
					modelField.FieldOptions = Models.OptionModels{}.ByKey(field["option_models_key"], field["option_beautify"].(int64) > 0)
				} else {
					modelField.FieldOptions = Models.OptionModels{}.ByKey(field["option_models_key"], field["option_beautify"].(int64) == 1)
				}
				if field["set_as_tabs"].(int64) == 1 {
					//!!!!!!!! 将某个字段和字段的选项集设置为列表tab页 !!!!!!!!!!!
					for _, options := range modelField.FieldOptions {
						newTab := EasyModelListTab{
							TabName:         options["name"].(string),
							SelectCondition: field["field_key"].(string) + "='" + util.Interface2String(options["value"]) + "'",
							SearchFormParams: map[string]interface{}{
								field["field_key"].(string): util.Interface2String(options["value"]),
							},
						}
						easyModel.ListTabs = append(easyModel.ListTabs, newTab)
					}
				}
			}
			easyModel.Fields = append(easyModel.Fields, modelField)
		}

		//!模型搜索表单信息
		searchForm, err := db.New().Table("tb_easy_models_search_form").
			Where("is_delete", 0).
			Where("status", 1).
			Where("model_id", modelInfo["id"].(int64)).
			Order("index_num asc,id asc").
			Get()
		if err != nil {
			logger.Error(err.Error())
			return EasyModel{}, errors.New("系统运行错误！")
		}
		for _, item := range searchForm {
			var searchFields []string
			if item["search_fields"].(string) != "" {
				util.JsonDecode(item["search_fields"].(string), &searchFields)
			}
			//# 格式化追加选项集
			var OptionModelsAdds []map[string]interface{}
			optionModelsAdd := strings.Split(item["option_models_add"].(string), "\n")
			for _, omAdd := range optionModelsAdd {
				omPs := strings.Split(omAdd, "|")
				if len(omPs) == 2 {
					OptionModelsAdds = append(OptionModelsAdds, map[string]interface{}{
						"value": omPs[0],
						"name":  omPs[1],
						"sql":   "",
					})
				} else if len(omPs) == 3 {
					OptionModelsAdds = append(OptionModelsAdds, map[string]interface{}{
						"value": omPs[0],
						"name":  omPs[1],
						"sql":   omPs[2],
					})
				}
			}
			searchFormItem := SearchFromItem{
				DataType:        item["data_type"].(string),
				SearchKey:       item["search_key"].(string),
				SearchName:      item["search_name"].(string),
				Placeholder:     item["placeholder"].(string),
				OptionModelsKey: item["option_models_key"].(string),
				OptionModelsAdd: OptionModelsAdds,
				SearchFields:    searchFields,
				MatchType:       item["match_type"].(string),
				DefaultValue:    item["default_value"].(string),
				SubQuery:        item["sub_query"].(string),
			}
			easyModel.SearchForm = append(easyModel.SearchForm, searchFormItem)
		}

		//!保存并返回
		easyModelList[saveKey] = easyModel
		//!定时删除数据
		go func() {
			t := time.After(time.Second * 10) //十秒钟后删除
			_, _ = <-t
			//logger.Info("推迟一分钟，删除EasyModelList 的" + modelKey)
			easyModelListLock.Lock()
			defer easyModelListLock.Unlock()
			delete(easyModelList, saveKey)
		}()
		return easyModel, nil
	}
}
