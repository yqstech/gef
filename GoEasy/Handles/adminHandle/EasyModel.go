/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyModel
 * @Version: 1.0.0
 * @Date: 2022/5/2 10:54
 */

package adminHandle

import (
	"errors"
	"github.com/gohouse/gorose/v2"
	"github.com/yqstech/gef/GoEasy/EasyApp"
	EasyModel2 "github.com/yqstech/gef/GoEasy/EasyModel"
	"github.com/yqstech/gef/GoEasy/Models"
	"github.com/yqstech/gef/GoEasy/Utils/util"
	"strings"
)

type EasyModel struct {
	Base
	ModelKey string //模型关键字
}

// NodeBegin 开始
func (that EasyModel) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	easyModel, err := EasyModel2.GetEasyModelInfo(that.ModelKey, "begin")
	if err != nil {
		return err, 500
	} else {
		pageData.SetTitle(easyModel.ModelName + "管理")
		pageData.SetPageName(easyModel.ModelName)
		pageData.SetTbName(easyModel.TableName)
		return nil, 0
	}
}

// NodeList 初始化列表
func (that EasyModel) NodeList(pageData *EasyApp.PageData) (error, int) {
	easyModel, err := EasyModel2.GetEasyModelInfo(that.ModelKey, "list")
	if err != nil {
		return err, 500
	} else {
		//!排序、分页、列表页面提示
		pageData.SetListOrder(easyModel.OrderType)
		if easyModel.PageSize > 0 {
			pageData.SetListPageSize(easyModel.PageSize)
		}
		if easyModel.PageNotice != "" {
			pageData.SetPageNotice(easyModel.PageNotice)
		}
		//!设置tabs列表和选中项
		validUrl := util.UrlScreenParam(pageData.GetHttpRequest(), []string{"id"}, false, true)
		for index, tab := range easyModel.ListTabs {
			pageData.PageTabAdd(tab.TabName, validUrl+"tab="+util.Int2String(index))
		}
		//获取第几页
		tabIndex := that.GetTabIndex(pageData, "tab")
		pageData.SetPageTabSelect(tabIndex)

		//!获取自定义按钮列表
		for btnName, Btn := range easyModel.Buttons {
			pageData.SetButton(btnName, Btn)
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
		pageData.SetListTopBtns(topBtns...)
		//!URL参数透传到顶部按钮上
		for _, urlParam := range easyModel.UrlParams {
			for _, btn := range topBtns {
				urlAppend := urlParam.FieldKey + "=" + util.GetValue(pageData.GetHttpRequest(), urlParam.ParamKey)
				pageData.SetButtonActionUrl(btn, urlAppend, true)
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
		pageData.SetListRightBtns(rightBtns...)
		//ID列同步到字段管理了，统一不显示默认id列了
		pageData.ListColumnClear()
		if easyModel.BatchAction {
			pageData.SetListBatchAction(true)
		}
		//增加列表列
		for _, field := range easyModel.Fields {
			if field.IsShowOnList {
				//自定义列名称
				if field.FieldNameReset != "" {
					field.FieldName = field.FieldNameReset
				}
				pageData.ListColumnAdd(field.FieldKey, field.FieldName, field.DataTypeOnList+"::"+field.DataTypeCommandOnList, field.FieldOptions)
				//自定义列样式
				if field.FieldStyleReset != "" {
					pageData.SetListColumnStyle(field.FieldKey, field.FieldStyleReset)
				}
			}
		}
		return nil, 0
	}
}

// NodeListCondition 修改查询条件
func (that EasyModel) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	easyModel, err := EasyModel2.GetEasyModelInfo(that.ModelKey, "list")
	if err != nil {
		return condition, nil, 0
	} else {
		//设置url参数查询条件
		for _, urlParam := range easyModel.UrlParams {
			if urlParam.ParamKey != "" && urlParam.FieldKey != "" {
				urlValue := util.GetValue(pageData.GetHttpRequest(), urlParam.ParamKey)
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
		tabIndex := that.GetTabIndex(pageData, "tab")
		if len(easyModel.ListTabs) > tabIndex {
			if easyModel.ListTabs[tabIndex].SelectCondition != "" {
				condition = append(condition, []interface{}{
					easyModel.ListTabs[tabIndex].SelectCondition,
				})
			}
		}
	}
	//追加查询条件
	return condition, nil, 0
}

// NodeListData 重写列表数据
func (that EasyModel) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	easyModel, err := EasyModel2.GetEasyModelInfo(that.ModelKey, "list")
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
func (that EasyModel) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	easyModel, err := EasyModel2.GetEasyModelInfo(that.ModelKey, util.Is(id == 0, "add", "edit").(string))
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
						pageData.FormFieldsAdd("", "block", field.GroupTitle, "", "", false, nil, "", nil)
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
					pageData.FormFieldsAdd(field.FieldKey, field.DataTypeOnCreate, field.FieldName, field.FieldNotice, field.DefaultValue, field.IsMust, FieldOptions, "", expand)
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
						pageData.FormFieldsAdd("", "block", field.GroupTitle, "", "", false, nil, "", nil)
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
					pageData.FormFieldsAdd(field.FieldKey, field.DataTypeOnUpdate, field.FieldName, field.FieldNotice, field.DefaultValue, field.IsMust, FieldOptions, "", expand)
				}
			}
		}

		return nil, 0
	}
}

// NodeFormData 表单显示前修改数据
func (that EasyModel) NodeFormData(pageData *EasyApp.PageData, data gorose.Data, id int64) (gorose.Data, error, int) {
	easyModel, err := EasyModel2.GetEasyModelInfo(that.ModelKey, util.Is(id == 0, "add", "edit").(string))
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
			}
		} else {
			//新增页链接参数透传
			for _, urlParam := range easyModel.UrlParams {
				field := urlParam.FieldKey
				value := util.GetValue(pageData.GetHttpRequest(), field)
				data[field] = value
			}
		}
	}
	return data, nil, 0
}

// NodeSaveData 表单保存数据前使用
func (that EasyModel) NodeSaveData(pageData *EasyApp.PageData, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	easyModel, err := EasyModel2.GetEasyModelInfo(that.ModelKey, util.Is(oldData["id"].(int64) == 0, "add", "edit").(string))
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
		}
	}
	return postData, nil, 0
}
