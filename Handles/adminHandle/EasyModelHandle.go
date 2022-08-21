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
	"github.com/yqstech/gef/EasyApp"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/util"
	"strings"
	"sync"
	"time"
)

type EasyModelHandle struct {
	Base
	ModelKey string //模型关键字
}

// NodeBegin 开始
func (that EasyModelHandle) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	easyModel, err := GetEasyModelInfo(that.ModelKey, "begin")
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
func (that EasyModelHandle) NodeList(pageData *EasyApp.PageData) (error, int) {
	easyModel, err := GetEasyModelInfo(that.ModelKey, "list")
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
func (that EasyModelHandle) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	easyModel, err := GetEasyModelInfo(that.ModelKey, "list")
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
func (that EasyModelHandle) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	easyModel, err := GetEasyModelInfo(that.ModelKey, "list")
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
func (that EasyModelHandle) NodeForm(pageData *EasyApp.PageData, id int64) (error, int) {
	easyModel, err := GetEasyModelInfo(that.ModelKey, util.Is(id == 0, "add", "edit").(string))
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
func (that EasyModelHandle) NodeFormData(pageData *EasyApp.PageData, data gorose.Data, id int64) (gorose.Data, error, int) {
	easyModel, err := GetEasyModelInfo(that.ModelKey, util.Is(id == 0, "add", "edit").(string))
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
func (that EasyModelHandle) NodeSaveData(pageData *EasyApp.PageData, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {
	easyModel, err := GetEasyModelInfo(that.ModelKey, util.Is(oldData["id"].(int64) == 0, "add", "edit").(string))
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
	OrderType    string                    //排序方式
	PageSize     int                       //分页大小
	PageNotice   string                    //页面备注
	Buttons      map[string]EasyApp.Button //页面
	TopButtons   []string                  //顶部按钮
	RightButtons []string                  //右侧操作按钮
	UrlParams    []EasyModelUrlParam       //链接参数
	LevelIndent  string                    //按级缩进
	BatchAction  bool                      //是否支持批量操作
}
type EasyModelListTab struct {
	TabName         string //tab页名称
	SelectCondition string //查询条件
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

// easyModelList 存储的模型信息列表
var easyModelList = map[string]EasyModel{}
var easyModelListLock sync.Mutex

// GetEasyModelInfo 获取模型信息
func GetEasyModelInfo(modelKey string, actionName string) (EasyModel, error) {
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
			Buttons:      map[string]EasyApp.Button{},
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
					TabName:         "标签页",
					SelectCondition: "",
				}
				if OptsLength > 0 && Opts[0] != "" {
					newTab.TabName = Opts[0]
				}
				if OptsLength > 1 && Opts[1] != "" {
					newTab.SelectCondition = Opts[1]
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
		if len(allButton) > 0 {
			selfButtonList, err := db.New().Table("tb_easy_models_buttons").
				Where("is_delete", 0).
				Where("status", 1).
				WhereIn("button_key", allButton).
				Get()
			if err != nil {
				logger.Error(err.Error())
				return EasyModel{}, errors.New("系统运行错误！")
			}
			for _, btnInfo := range selfButtonList {
				easyModel.Buttons[btnInfo["button_key"].(string)] = EasyApp.Button{
					ButtonName: btnInfo["button_name"].(string),
					Action:     btnInfo["action"].(string),
					ActionType: util.Int642Int(btnInfo["action_type"].(int64)),
					ConfirmMsg: btnInfo["confirm_msg"].(string),
					LayerTitle: btnInfo["layer_title"].(string),
					ActionUrl:  btnInfo["action_url"].(string),
					Class:      btnInfo["class_name"].(string),
					Icon:       btnInfo["button_icon"].(string),
					Display:    btnInfo["display"].(string),
					Expand: map[string]string{
						"w": btnInfo["layer_width"].(string),
						"h": btnInfo["layer_height"].(string),
					},
					BatchAction: util.Is(btnInfo["batch_action"].(int64) == 1, true, false).(bool),
				}
			}
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
						}
						easyModel.ListTabs = append(easyModel.ListTabs, newTab)
					}
				}
			}
			easyModel.Fields = append(easyModel.Fields, modelField)
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