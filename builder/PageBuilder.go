/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2022/8/21 22:01
 */

package builder

import (
	"encoding/json"
	"errors"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/config"
	"github.com/yqstech/gef/util"
	"html"
	"net/http"
	"strings"
)

// PageBuilder 页面构建器
type PageBuilder struct {
	//!请求ID
	RequestID string
	////!页面对象
	//ActivePage NodePage
	////!页面对象的方法列表
	//ActionList map[string]PageAction //页面的处理方法列表

	//! http响应
	HttpResponseWriter http.ResponseWriter
	//! http请求
	HttpRequest *http.Request
	//! http参数
	HttpParams httprouter.Params

	// ====================  主数据  ================================
	//页面样式
	pageStyle string
	// 标题
	title string
	// 列表模板名称
	listTplName string
	// 数据表名称
	tbName string
	//数据表主键
	tbPK string
	//数据库自动设置ID
	isAutoID bool
	//标记删除的字段
	deleteField string
	//插入自动完成的字段列表
	insertAutoFields []string
	//修改自动完成的字段列表
	updateAutoFields []string
	// 页面公告
	pageNotice string
	// 页面名称
	pageName string
	// 操作名称（例如列表，新增，修改等）
	actionName string
	//Tab选项卡列表
	pageTabs []pageTab
	//Tab选项卡选中指示
	pageTabSelect int

	// 列表页数据查询地址
	listDataUrl string

	//是否隐藏分页组件
	listPageHide bool

	//列表是否支持批量操作
	listBatchAction bool

	//列表当前页码
	listPage int
	// 列表页数据分页大小
	listPageSize int

	// 列表查询字段
	listFields string
	// 列表数据中需要删除的字段
	listFieldsRemove []string
	// 列表查询sql排序
	listOrder string
	// 列表页查询条件
	listCondition [][]interface{}

	// 按钮列表
	buttons map[string]Button

	// 列表右侧按钮集合
	listRightBtns []string

	// 列表顶部按钮集合
	listTopBtns []string

	// 列表列集合
	listColumns []ListColumn
	// 列表列定义样式
	listColumnsStyles map[string]interface{}
	// 列表页 搜索表单项
	listSearchFields []ListSearchField

	//新增数据地址
	addDataUrl string
	//新增页面模板名称
	addTplName string

	//查询数据地址
	findDataUrl string
	//修改页面保存数据地址
	editDataUrl string
	//修改页面模板名称
	editTplName string

	formFields      []FormField
	formFieldKeys   []interface{}
	formData        gorose.Data
	formSubmitTitle string
	formSubmitHide  int
	uploadImageUrl  string
}

// =============    数据管理方法   ================

// DataReset 数据重置和初始化
// 所有操作方法运行前，都需要重置数据，否则会遗留之前的内容
func (builder *PageBuilder) DataReset() {
	builder.pageStyle = ""
	builder.title = ""
	builder.listTplName = "list.html"
	builder.tbName = ""
	builder.tbPK = "id"
	builder.isAutoID = true
	builder.deleteField = "is_delete"
	builder.insertAutoFields = []string{"create_time", "update_time"}
	builder.updateAutoFields = []string{"update_time"}
	builder.pageNotice = ""
	builder.pageName = "未知模块"
	builder.actionName = "未知操作"
	builder.pageTabs = []pageTab{}
	builder.pageTabSelect = 0
	builder.listFields = "*"
	builder.listFieldsRemove = []string{"is_delete"}
	builder.listOrder = "id desc"
	builder.listDataUrl = ""
	builder.listPage = 1
	builder.listPageSize = 20
	builder.listCondition = [][]interface{}{}
	builder.buttons = map[string]Button{
		"add": {
			ButtonName: "新增",
			Action:     "add",
			ActionType: 2,
			ActionUrl:  "add",
			Class:      "def",
			Icon:       "ri-add-circle-line",
			Display:    "",
			Expand: map[string]string{
				"w": "98%",
				"h": "98%",
			},
		},
		"edit": {
			ButtonName: "编辑",
			Action:     "edit",
			ActionType: 2,
			ActionUrl:  "edit",
			Class:      "layui-btn-normal",
			Icon:       "ri-edit-box-line",
			Display:    "!item.btn_edit || item.btn_edit!='hide'",
			Expand: map[string]string{
				"w": "98%",
				"h": "98%",
			},
		},
		"disable": {
			ButtonName:  "禁用",
			Action:      "status",
			ActionType:  1,
			ActionUrl:   "status?status=0",
			Class:       "layui-btn-warm",
			Icon:        "ri-forbid-line",
			Display:     "(!item.btn_status || item.btn_status!='hide') && item.status==1",
			BatchAction: true,
		},
		"enable": {
			ButtonName:  "启用",
			Action:      "status",
			ActionType:  1,
			ActionUrl:   "status?status=1",
			Class:       "",
			Icon:        "ri-checkbox-circle-line",
			Display:     "(!item.btn_status || item.btn_status!='hide') && item.status==0",
			BatchAction: true,
		},
		"delete": {
			ButtonName:  "删除",
			Action:      "delete",
			ActionType:  1,
			ActionUrl:   "delete",
			ConfirmMsg:  "确定要删除此项？",
			Class:       "layui-btn-danger",
			Icon:        "ri-delete-bin-2-line",
			Display:     "!item.btn_delete || item.btn_delete!='hide'",
			BatchAction: true,
		},
	}
	builder.listRightBtns = []string{"edit", "disable", "enable", "delete"}
	builder.listTopBtns = []string{"add"}
	//重置列表项
	builder.listColumns = []ListColumn{
		{
			builder.tbPK, "ID", "text", nil, nil,
		},
	}
	builder.listColumnsStyles = map[string]interface{}{}
	builder.listSearchFields = []ListSearchField{}
	builder.listBatchAction = false

	builder.addDataUrl = ""
	builder.addTplName = "add.html"

	builder.findDataUrl = "find"
	builder.editDataUrl = ""
	builder.editTplName = "edit.html"

	builder.formFields = []FormField{}
	builder.formFieldKeys = []interface{}{}
	builder.formSubmitTitle = ""
	builder.formSubmitHide = 0
	builder.uploadImageUrl = "upload"
}

// SetStyle 设置样式
func (builder *PageBuilder) SetStyle(style string) {
	builder.pageStyle = builder.pageStyle + style
}

func (builder *PageBuilder) GetStyle() string {
	return builder.pageStyle
}

// SetTitle 设置title
func (builder *PageBuilder) SetTitle(tit string) {
	builder.title = tit
}

// SetListTplName 设置列表模板名称
func (builder *PageBuilder) SetListTplName(tit string) {
	builder.listTplName = tit
}
func (builder *PageBuilder) GetListTplName() string {
	return builder.listTplName
}

// SetTbName 设置列表模板名称
func (builder *PageBuilder) SetTbName(tbName string) {
	builder.tbName = tbName
}

// SetTbName 设置列表模板名称
func (builder *PageBuilder) GetTbName() string {
	return builder.tbName
}

// SetPK 设置数据表主键
func (builder *PageBuilder) SetPK(pk string) {
	builder.tbPK = pk
}

// GetPK 设置数据表主键
func (builder *PageBuilder) GetPK() string {
	return builder.tbPK
}

func (builder *PageBuilder) SetIsAutoID(isauto bool) {
	builder.isAutoID = isauto
}

// SetInsertAutoFields 设置插入时自动赋值字段
func (builder *PageBuilder) SetInsertAutoFields(fields ...string) {
	builder.insertAutoFields = []string{}
	builder.insertAutoFields = append(builder.insertAutoFields, fields...)
}

// SetUpdateAutoFields 设置更新时自动赋值字段
func (builder *PageBuilder) SetUpdateAutoFields(fields ...string) {
	builder.updateAutoFields = []string{}
	builder.updateAutoFields = append(builder.updateAutoFields, fields...)
}

// SetDeleteField 设置删除标记字段
func (builder *PageBuilder) SetDeleteField(field string) {
	builder.deleteField = field
}

// SetPageNotice 设置列表模板名称
func (builder *PageBuilder) SetPageNotice(str string) {
	builder.pageNotice = str
}

// SetPageName 设置结构体名称(控制器名称)
func (builder *PageBuilder) SetPageName(str string) {
	builder.pageName = str
}

// SetActionName 设置action名称(动作名称)
func (builder *PageBuilder) SetActionName(str string) {
	builder.actionName = str
}

// PageTabAdd 增加一个Tab选项卡
func (builder *PageBuilder) PageTabAdd(title, href string) {
	pt := pageTab{
		Href:  href,
		Title: title,
	}
	builder.pageTabs = append(builder.pageTabs, pt)
}

// SetPageTabSelect 设置Tab选项卡的选中
func (builder *PageBuilder) SetPageTabSelect(index int) {
	builder.pageTabSelect = index
}

// SetListFields 设置列表查询数据库字段
func (builder *PageBuilder) SetListFields(fields string) {
	builder.listFields = fields
}
func (builder *PageBuilder) GetListFields() string {
	return builder.listFields
}

// SetListFieldsRemove 设置列表需要去掉的字段数据
func (builder *PageBuilder) SetListFieldsRemove(fields ...string) {
	builder.listFieldsRemove = fields
}
func (builder *PageBuilder) GetListFieldsRemove() []string {
	return builder.listFieldsRemove
}

// SetListOrder 设置列表查询排序方式
func (builder *PageBuilder) SetListOrder(order string) {
	builder.listOrder = order
}
func (builder *PageBuilder) GetListOrder() string {
	return builder.listOrder
}

// SetListDataURL 设置列表查询数据地址，默认为空，代表当前地址
func (builder *PageBuilder) SetListDataURL(url string) {
	builder.listDataUrl = url
}
func (builder *PageBuilder) GetListDataURL() string {
	return builder.listDataUrl
}

// SetListPageSize 设置列表页分页大小
func (builder *PageBuilder) SetListPageSize(size int) {
	builder.listPageSize = size
}
func (builder *PageBuilder) GetListPageSize() int {
	return builder.listPageSize
}

// SetListPage 设置列表页码
func (builder *PageBuilder) SetListPage(page int) {
	builder.listPage = page
}
func (builder *PageBuilder) GetListPage() int {
	return builder.listPage
}

// SetListPageHide 隐藏分页
func (builder *PageBuilder) SetListPageHide() {
	builder.listPageHide = true
}

// SetListBatchAction 开启批量操作
func (builder *PageBuilder) SetListBatchAction(isOpen bool) {
	builder.listBatchAction = isOpen
}

// SetListCondition 设置列表页查询条件
func (builder *PageBuilder) SetListCondition(c [][]interface{}) {
	builder.listCondition = c
}

// GetListCondition 设置列表页查询条件
func (builder *PageBuilder) GetListCondition() [][]interface{} {
	return builder.listCondition
}

// ListConditionAdd 增加列表页查询条件
func (builder *PageBuilder) ListConditionAdd(c []interface{}) {
	builder.listCondition = append(builder.listCondition, c)
}

func (builder *PageBuilder) SetButton(btnName string, btn Button) {
	builder.buttons[btnName] = btn
}
func (builder *PageBuilder) GetButtons() map[string]Button {
	return builder.buttons
}

// SetListRightBtns 设置列表右侧使用的按钮
func (builder *PageBuilder) SetListRightBtns(btns ...string) {
	builder.listRightBtns = btns
}

// ListRightBtnsClear 清除右侧按钮
func (builder *PageBuilder) ListRightBtnsClear() {
	builder.listRightBtns = []string{}
}

// ListRightBtnsIconClear 清除右侧按钮的图标
func (builder *PageBuilder) ListRightBtnsIconClear() {
	for _, btnName := range builder.listRightBtns {
		if btn, ok := builder.buttons[btnName]; ok {
			btn.Icon = ""
			builder.buttons[btnName] = btn
		}
	}
}

// SetButtonIcon 设置按钮图标
func (builder *PageBuilder) SetButtonIcon(btnKey, icon string) {
	if btn, ok := builder.buttons[btnKey]; ok {
		btn.Icon = icon
		builder.buttons[btnKey] = btn
	}
}

// SetButtonName 设置按钮名称
func (builder *PageBuilder) SetButtonName(btnKey, ButtonName string) {
	if btn, ok := builder.buttons[btnKey]; ok {
		btn.ButtonName = ButtonName
		builder.buttons[btnKey] = btn
	}
}

// SetButtonActionUrl 设置按钮链接地址 是否是追加
func (builder *PageBuilder) SetButtonActionUrl(btnName, url string, isAddend bool) {
	if btn, ok := builder.buttons[btnName]; ok {
		if isAddend {
			if strings.Contains(btn.ActionUrl, "?") {
				btn.ActionUrl = btn.ActionUrl + "&" + url
			} else {
				btn.ActionUrl = btn.ActionUrl + "?" + url
			}
		} else {
			btn.ActionUrl = url
		}
		builder.buttons[btnName] = btn
	}
}

// ListRightBtnsAdd 右侧按钮新增
func (builder *PageBuilder) ListRightBtnsAdd(btns ...string) {
	builder.listRightBtns = append(builder.listRightBtns, btns...)
}

func (builder *PageBuilder) GetListRightBtns() []string {
	return builder.listRightBtns
}

// SetListTopBtns 设置列表顶部使用的列表
func (builder *PageBuilder) SetListTopBtns(btns ...string) {
	builder.listTopBtns = btns
}

// ListTopBtnsClear 清除顶部按钮
func (builder *PageBuilder) ListTopBtnsClear() {
	builder.listTopBtns = []string{}
}

// SetListColumns 设置列表列信息
func (builder *PageBuilder) SetListColumns(data []ListColumn) {
	builder.listColumns = data
}

func (builder *PageBuilder) GetListColumns() []ListColumn {
	return builder.listColumns
}

// SetListColumnStyle 设置列表列信息
func (builder *PageBuilder) SetListColumnStyle(key, value string) {
	builder.listColumnsStyles[key] = value
}

// ListColumnAdd 列表新增一列信息
func (builder *PageBuilder) ListColumnAdd(FieldName, ColumnName, DataType string, Data []map[string]interface{}) {
	//数据自动转map
	MapData := map[interface{}]interface{}{}
	for _, v := range Data {
		//第一种写法，含有name和value参数
		if _, ok := v["value"]; ok {
			if _, ok2 := v["name"]; ok2 {
				MapData[v["value"]] = v["name"]
				continue
			}
		}
		//第二种直接是map写法
		for dataK, dataV := range Data {
			MapData[dataK] = dataV
		}
	}

	//DataType参数支持::指令写法,可以设置组件参数
	Options := map[string]interface{}{}
	if strings.Contains(DataType, "::") {
		//解析指令
		options := strings.Split(DataType, "::")
		//::前面的为组件名称
		DataType = options[0]
		//::后面的为指令参数
		if len(options) > 1 && len(options[1]) > 0 {
			//指令参数支持json格式或url格式
			if options[1][0:1] == "{" {
				json.Unmarshal([]byte(options[1]), &Options)
			} else {
				//url链接写法
				opts := strings.Split(options[1], "&")
				for _, opt := range opts {
					optKv := strings.Split(opt, "=")
					if len(optKv) == 1 {
						optKv = strings.Split(opt, ":")
					}
					if len(optKv) == 1 {
						continue
					}
					Options[optKv[0]] = optKv[1]
				}
			}
		}
	}
	builder.listColumns = append(builder.listColumns, ListColumn{
		FieldName:  FieldName,
		ColumnName: ColumnName,
		DataType:   DataType,
		Data:       MapData,
		Options:    Options,
	})
}

// ListColumnClear 列表清除
func (builder *PageBuilder) ListColumnClear() {
	builder.listColumns = []ListColumn{}
}

// ListSearchFieldAdd 增加列表搜索项
func (builder *PageBuilder) ListSearchFieldAdd(fkey, ftype, ftitle string, defvalue interface{}, value interface{}, data []map[string]interface{},
	style string,
	expand map[string]interface{}) {

	//数据自动转map
	MapData := map[interface{}]interface{}{}
	for _, v := range data {
		MapData[v["value"]] = v["name"]
	}

	builder.listSearchFields = append(builder.listSearchFields, ListSearchField{
		Key:      fkey,
		Type:     ftype,
		Title:    ftitle,
		DefValue: defvalue,
		Value:    value,
		Data:     MapData,
		Style:    style,
		Expand:   expand,
	})
}

// FormFieldsAdd 表单项增加
func (builder *PageBuilder) FormFieldsAdd(fkey, ftype, ftitle, fnotice, fvalue string, ismust bool,
	dataOrOptions []map[string]interface{},
	style string,
	expand map[string]interface{}) {
	builder.formFields = append(builder.formFields, FormField{
		Key:    fkey,
		Type:   ftype,
		Title:  ftitle,
		Notice: fnotice,
		Value:  fvalue,
		IsMust: ismust,
		Data:   dataOrOptions,
		Style:  style,
		Expand: expand,
	})
}

// FormFieldKeysAdd 新增表单字段key列表
func (builder *PageBuilder) FormFieldKeysAdd(fieldKeys ...interface{}) {
	builder.formFieldKeys = append(builder.formFieldKeys, fieldKeys)
}

// SetAddDataUrl 设置新增页提交地址
func (builder *PageBuilder) SetAddDataUrl(url string) {
	builder.addDataUrl = url
}
func (builder *PageBuilder) GetAddDataUrl() string {
	return builder.addDataUrl
}

// SetAddTplName 设置新增页模板名称
func (builder *PageBuilder) SetAddTplName(tplName string) {
	builder.addTplName = tplName
}
func (builder *PageBuilder) GetAddTplName() string {
	return builder.addTplName
}

// SetFindTplName 设置新增页提交地址
func (builder *PageBuilder) SetFindTplName(url string) {
	builder.findDataUrl = url
}
func (builder *PageBuilder) GetFindTplName() string {
	return builder.findDataUrl
}

// SetEditDataUrl 设置新增页提交地址
func (builder *PageBuilder) SetEditDataUrl(url string) {
	builder.editDataUrl = url
}
func (builder *PageBuilder) GetEditDataUrl() string {
	return builder.editDataUrl
}

// SetEditTplName 设置新增页模板名称
func (builder *PageBuilder) SetEditTplName(tplName string) {
	builder.editTplName = tplName
}
func (builder *PageBuilder) GetEditTplName() string {
	return builder.editTplName
}

// SetFormData 设置表单默认数据
func (builder *PageBuilder) SetFormData(formData gorose.Data) {
	builder.formData = formData
}

// SetFormSubmitTitle 设置表单提交按钮文字
func (builder *PageBuilder) SetFormSubmitTitle(title string) {
	builder.formSubmitTitle = title
}

// SetFormSubmitHide 设置表单按钮隐藏
func (builder *PageBuilder) SetFormSubmitHide() {
	builder.formSubmitHide = 1
}

// SetUploadImageUrl 设置上传图片的地址
func (builder *PageBuilder) SetUploadImageUrl(url string) {
	builder.uploadImageUrl = url
}

func (builder *PageBuilder) SetHttpWriter(w http.ResponseWriter) {
	builder.HttpResponseWriter = w
}
func (builder *PageBuilder) SetHttpRequest(r *http.Request) {
	builder.HttpRequest = r
}
func (builder *PageBuilder) SetHttpParams(ps httprouter.Params) {
	builder.HttpParams = ps
}

func (builder *PageBuilder) GetHttpWriter() http.ResponseWriter {
	return builder.HttpResponseWriter
}
func (builder *PageBuilder) GetHttpRequest() *http.Request {
	return builder.HttpRequest
}
func (builder *PageBuilder) GetHttpParams() httprouter.Params {
	return builder.HttpParams
}

// GetEasyModelsButtons 获取按钮列表
func (builder *PageBuilder) GetEasyModelsButtons(allButton []interface{}) (map[string]Button, error) {
	buttons := map[string]Button{}
	if len(allButton) > 0 {
		selfButtonList, err := db.New().Table("tb_easy_models_buttons").
			Where("is_delete", 0).
			Where("status", 1).
			WhereIn("button_key", allButton).
			Get()
		if err != nil {
			logger.Error(err.Error())
			return buttons, errors.New("系统运行错误！")
		}
		for _, btnInfo := range selfButtonList {
			//!按钮地址补充后台路径前缀
			if len(btnInfo["action_url"].(string)) > 1 && btnInfo["action_url"].(string)[0:1] == "/" {
				btnInfo["action_url"] = config.AdminPath + btnInfo["action_url"].(string)
			}
			buttons[btnInfo["button_key"].(string)] = Button{
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
	return buttons, nil
}

// TemplateData 设置模板数据
func (builder *PageBuilder) TemplateData() map[string]interface{} {
	data := map[string]interface{}{}
	//页面样式
	data["pageStyle"] = builder.GetStyle()
	//设置标题
	if data["title"] == nil {
		data["title"] = builder.title
	}
	data["pageName"] = builder.pageName
	data["pageNotice"] = builder.pageNotice
	//列表列
	data["listColumns"] = builder.listColumns
	data["listColumnsStyles"] = builder.listColumnsStyles
	//列表数据地址
	data["listDataUrl"] = builder.listDataUrl
	//列表按钮
	rightBtns := []Button{}
	for _, button_name := range builder.listRightBtns {
		if btn, ok := builder.buttons[button_name]; ok {
			rightBtns = append(rightBtns, btn)
		}
	}
	data["rightBtns"] = rightBtns
	//列表顶部按钮
	topBtns := []Button{}
	for _, button_name := range builder.listTopBtns {
		if btn, ok := builder.buttons[button_name]; ok {
			topBtns = append(topBtns, btn)
		}
	}
	data["topBtns"] = topBtns
	//列表搜索项
	data["listSearchFields"] = builder.listSearchFields
	//列表Tab
	data["pageTabs"] = builder.pageTabs
	data["pageTabSelect"] = builder.pageTabSelect
	data["pageTabsLength"] = len(builder.pageTabs)
	//分页
	data["listPageHide"] = builder.listPageHide //是否隐藏分页
	data["listPageSize"] = builder.listPageSize //分页
	//批量操作
	data["listBatchAction"] = builder.listBatchAction //批量操作

	//表单
	//提交地址
	data["editDataUrl"] = builder.editDataUrl
	data["addDataUrl"] = builder.addDataUrl
	//表单项列表
	for k, v := range builder.formFields {
		fkey := v.Key
		ftype := v.Type
		if fkey != "" {
			//初始化表单项的原始数据值
			if dbvalue, ok := builder.formData[fkey]; ok {
				builder.formFields[k].Value = Interface2String(dbvalue)
			}
			//修正checkbox默认值
			if IsInArray(ftype, []interface{}{"checkbox", "checkbox_level", "tags", "images"}) {
				if builder.formFields[k].Value == "" {
					builder.formFields[k].Value = "[]"
				}
			}
			//文本类组件需要转义双引号
			if IsInArray(ftype, []interface{}{"text", "textarea", "text-disabled", "textarea-disabled", "wangEditor"}) {
				builder.formFields[k].Value = strings.ReplaceAll(builder.formFields[k].Value, "\"", "\\\"")
			}
			//文本域,换行\n 避免被html解析造成换行错乱
			if ftype == "textarea" || ftype == "textarea-disabled" {
				builder.formFields[k].Value = strings.ReplaceAll(builder.formFields[k].Value, "\n", "\\n")
			}
			if ftype == "code" || ftype == "codeEditor" {
				builder.formFields[k].Value = html.EscapeString(builder.formFields[k].Value)
			}

		}
	}
	data["formFields"] = builder.formFields

	if builder.formSubmitTitle == "" {
		builder.formSubmitTitle = builder.actionName + builder.pageName
	}
	data["formSubmitTitle"] = builder.formSubmitTitle

	data["formSubmitHide"] = builder.formSubmitHide

	//组件需要
	//上传图片的地址
	data["uploadImageUrl"] = builder.uploadImageUrl

	return data
}

// 页面选项卡结构
type pageTab struct {
	Href  string
	Title string
}

// Button 定义按钮结构体
type Button struct {
	ButtonName  string //按钮名称
	Action      string //权限校验规则,add 或 a/b
	ActionType  int    //操作类型 1、ajax操作 2、弹出页面 3、javascript
	ConfirmMsg  string //确认对话框信息，ActionType=1有效
	LayerTitle  string //弹出窗口标题
	ActionUrl   string //操作地址
	Class       string //样式-类
	Icon        string //icon class
	Display     string //展示条件
	Expand      map[string]string
	BatchAction bool //是否支持批量操作
}

// ListColumn 列表展示列数据结构
type ListColumn struct {
	FieldName  string
	ColumnName string
	DataType   string
	Data       map[interface{}]interface{}
	Options    map[string]interface{} //参数
}

// ListSearchField 列表搜索表单项
type ListSearchField struct {
	Key      string                      //项目名称 key
	Type     string                      //项目类型 text select等
	Title    string                      //项目表标题
	DefValue interface{}                 //默认值空值
	Value    interface{}                 //预设值
	Data     map[interface{}]interface{} //数据列表
	Style    string                      //样式
	Expand   map[string]interface{}      //拓展数据
}

type FormField struct {
	Key    string
	Type   string                   //项目类型 text select等
	Title  string                   //项目表标题
	Notice string                   //项目备注
	Value  string                   //默认值
	IsMust bool                     //是否必传
	Data   []map[string]interface{} //数据列表
	Style  string                   //样式
	Expand map[string]interface{}   //拓展数据
}
