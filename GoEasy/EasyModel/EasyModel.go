/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyModel
 * @Version: 1.0.0
 * @Date: 2022/5/2 16:15
 */

package EasyModel

import (
	"errors"
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Models"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/wonderivan/logger"
	"strings"
	"sync"
	"time"
)

// EasyModel 定义模型信息结构体
type EasyModel struct {
	ModelKey     string                    //模型关键字
	ModelName    string                    //模型名称
	TableName    string                    //数据表名称
	AllowCreate  bool                      //允许新增
	AllowUpdate  bool                      //允许新增
	AllowDelete  bool                      //允许删除
	AllowStatus  bool                      //允许改状态
	Fields       []ModelField              //字段列表
	ListTabs     []Tab                     //列表多Tab页
	OrderType    string                    //排序方式
	PageSize     int                       //分页大小
	PageNotice   string                    //页面备注
	Buttons      map[string]EasyApp.Button //页面
	TopButtons   []string                  //顶部按钮
	RightButtons []string                  //右侧操作按钮
	UrlParams    []UrlParam                //链接参数
}
type Tab struct {
	TabName         string //tab页名称
	SelectCondition string //查询条件
}
type UrlParam struct {
	ParamKey     string //参数key
	FieldKey     string //数据库字段
	DefaultValue string //默认值
}

// ModelField 模型字段信息
type ModelField struct {
	FieldKey              string                   //字段关键字
	FieldName             string                   //字段名称
	FieldNameReset        string                   //字段重置名称（列表）
	FieldStyleReset       string                   //字段重设样式（列表）
	FieldNotice           string                   //字段提示信息
	IsShowOnList          bool                     //是否在列表页展示该字段
	DataTypeOnList        string                   //列表页展示时的数据类型
	DataTypeCommandOnList string                   //列表页展示时数据类型的指令
	DataTypeOnCreate      string                   //表单页展示时的数据类型
	DataTypeOnUpdate      string                   //表单页展示时的数据类型
	AllowCreate           bool                     //新增页显示
	AllowUpdate           bool                     //修改页显示
	IsMust                bool                     //是否必填
	DefaultValue          string                   //默认值
	FieldOptions          []map[string]interface{} //待选数据
	OptionIndent          bool                     //选项按照上下级缩进
	ExpandIf              string                   //拓展数据-if条件
	GroupTitle            string                   //创建一个新分组的标题
	FieldAugment          string                   //数据增强转换规则
	Attach2Field          string                   //附加到其他字段
	SaveTransRule         string                   //保存时数据变换规则
}

// easyModelList 存储的模型信息列表
var easyModelList = map[string]EasyModel{}
var easyModelListLock sync.Mutex

// GetEasyModelInfo 获取模型信息
func GetEasyModelInfo(modelKey string) (EasyModel, error) {
	if easyModel, ok := easyModelList[modelKey]; ok {
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
			Fields:       []ModelField{},
			ListTabs:     []Tab{},
			OrderType:    modelInfo["order_type"].(string),
			PageSize:     util.Int642Int(modelInfo["page_size"].(int64)),
			PageNotice:   modelInfo["page_notice"].(string),
			Buttons:      map[string]EasyApp.Button{},
			TopButtons:   []string{},
			RightButtons: []string{},
			UrlParams:    []UrlParam{},
		}
		//格式化多tab页
		if modelInfo["tabs_for_list"].(string) != "" {
			tabs := strings.Split(modelInfo["tabs_for_list"].(string), "\n")
			for _, tab := range tabs {
				Opts := strings.Split(tab, "|")
				OptsLength := len(Opts)
				newTab := Tab{
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
			}
		}
		//!格式化url参数
		urlParams := strings.Split(modelInfo["url_params"].(string), "\n")
		for _, urlParam := range urlParams {
			params := strings.Split(urlParam, ":")
			if len(params) == 3 {
				easyModel.UrlParams = append(easyModel.UrlParams, UrlParam{
					ParamKey:     params[0],
					FieldKey:     params[1],
					DefaultValue: params[2],
				})
			} else if len(params) == 2 {
				easyModel.UrlParams = append(easyModel.UrlParams, UrlParam{
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
			modelField := ModelField{
				FieldKey:              field["field_key"].(string),
				FieldName:             field["field_name"].(string),
				FieldNameReset:        field["field_name_reset"].(string),
				FieldStyleReset:       field["field_style_reset"].(string),
				FieldNotice:           field["field_notice"].(string),
				IsShowOnList:          field["is_show_on_list"].(int64) == 1,
				DataTypeOnList:        field["data_type_on_list"].(string),
				DataTypeCommandOnList: field["data_type_command_on_list"].(string),
				DataTypeOnCreate:      field["data_type_on_create"].(string),
				DataTypeOnUpdate:      field["data_type_on_update"].(string),
				AllowCreate:           field["allow_create"].(int64) == 1,
				AllowUpdate:           field["allow_update"].(int64) == 1,
				IsMust:                field["is_must"].(int64) == 1,
				DefaultValue:          field["default_value"].(string),
				FieldOptions:          nil,
				OptionIndent:          field["option_indent"].(int64) == 1,
				ExpandIf:              field["expand_if"].(string),
				GroupTitle:            field["group_title"].(string),
				FieldAugment:          field["field_augment"].(string),
				Attach2Field:          field["attach_to_field"].(string),
				SaveTransRule:         field["save_trans_rule"].(string),
			}
			if field["option_models_id"].(int64) > 0 {
				modelField.FieldOptions = Models.OptionModels{}.ById(util.Int642Int(field["option_models_id"].(int64)), field["option_beautify"].(int64) == 1)
				if field["set_as_tabs"].(int64) == 1 {
					//!!!!!!!! 将某个字段和字段的选项集设置为列表tab页 !!!!!!!!!!!
					for _, options := range modelField.FieldOptions {
						newTab := Tab{
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
		easyModelList[modelKey] = easyModel
		//!定时删除数据
		go func() {
			t := time.After(time.Second * 10) //十秒钟后删除
			_, _ = <-t
			//logger.Info("推迟一分钟，删除EasyModelList 的" + modelKey)
			easyModelListLock.Lock()
			defer easyModelListLock.Unlock()
			delete(easyModelList, modelKey)
		}()
		return easyModel, nil
	}
}

// selectDataList 待选择数据
var selectDataList = map[int64][]map[string]interface{}{}
var selectDataListLock sync.Mutex

func GetSelectData(id int64) []map[string]interface{} {
	if selectData, ok := selectDataList[id]; ok {
		return selectData
	} else {
		selectDataListLock.Lock()
		defer selectDataListLock.Unlock()
		//通过数据库查询
		data, err := db.New().Table("tb_easy_models_fields_select_data").
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
			//通过数据表获取列表
			arrOptions, err, _ := Models.Model{}.SelectOptionsData(data["table_name"].(string), map[string]string{
				data["value_field"].(string): "value",
				data["name_field"].(string):  "name",
			}, "0", "未知", "", "")
			if err != nil {
				logger.Error(err.Error())
				return nil
			}
			selectData = arrOptions
		}
		//!定时删除
		go func() {
			t := time.After(time.Second * 10) //十秒钟后删除
			_, _ = <-t
			selectDataListLock.Lock()
			defer selectDataListLock.Unlock()
			delete(selectDataList, id)
		}()
		
		selectDataList[id] = selectData
		return selectData
	}
}
