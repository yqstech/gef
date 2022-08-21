/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: PageData
 * @Version: 1.0.0
 * @Date: 2021/12/7 9:30 下午
 */

package EasyApp

import (
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/yqstech/gef/Utils/pool"
	"github.com/yqstech/gef/Utils/util"
	"net/http"
	"strings"
	"time"
)

// PageData 流程数据
type PageData struct {
	RequestID  string                //请求ID
	ActivePage AppPage               //真实结构体对象
	ActionList map[string]PageAction //页面的处理方法列表

	//http请求
	httpW  http.ResponseWriter
	httpR  *http.Request
	httpPS httprouter.Params
	// ====================  主数据  ================================
	// 标题
	title string
	// 列表模板名称
	listTplName string
	// 数据表名称
	tbName string
	//数据表主键
	tbPK string
	//数据库自动设置ID
	isAutoID    bool
	deleteField string
	//插入自动完成
	insertAutoFields []string
	//修改自动完成
	updateAutoFields []string
	PageStyle        string
	// 页面公告
	pageNotice string
	// 结构体名称（结构体主题）
	pageName string
	// 动作名称（例如列表，新增，修改等）
	actionName string
	// pageTab Tab选项卡
	pageTabs []pageTab
	// pageTabSelect Tab选项卡选中第几个
	pageTabSelect int

	// 列表查询字段，默认全查
	listFields string

	// 列表数据中需要删除的字段
	listFieldsRemove []string

	// 列表页排序
	listOrder string

	// 列表页数据查询地址
	listDataUrl string

	listPage int

	//隐藏分页组件
	listPageHide bool

	//列表是否支持批量操作
	listBatchAction bool

	// 列表页数据分页大小
	listPageSize int

	// 列表页查询条件
	listCondition [][]interface{}

	// 全部按钮列表
	buttons map[string]Button

	// 列表页右侧按钮组
	listRightBtns []string

	// 列表顶部按钮组
	listTopBtns []string

	// 列表展示列列表
	listColumns []ListColumn
	//列表展示列样式
	listColumnsStyles map[string]interface{}
	// 列表页 搜索条件表单项 列表
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
	editTplName     string
	formFields      []FormField
	formFieldKeys   []interface{}
	formData        gorose.Data
	formSubmitTitle string
	formSubmitHide  int
	uploadImageUrl  string
}
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

func (fd *PageData) ActionAdd(actionName string, action PageAction) {
	fd.ActionList[actionName] = action
}

func (fd *PageData) ActionClear() {
	fd.ActionList = map[string]PageAction{}
}

// =============    数据管理方法   ================

// DataReset 数据重置和初始化
// 所有操作方法运行前，都需要重置数据，否则会遗留之前的内容
func (fd *PageData) DataReset() {
	fd.PageStyle = ""
	fd.title = ""
	fd.listTplName = "list.html"
	fd.tbName = ""
	fd.tbPK = "id"
	fd.isAutoID = true
	fd.deleteField = "is_delete"
	fd.insertAutoFields = []string{"create_time", "update_time"}
	fd.updateAutoFields = []string{"update_time"}
	fd.pageNotice = ""
	fd.pageName = "未知模块"
	fd.actionName = "未知操作"
	fd.pageTabs = []pageTab{}
	fd.pageTabSelect = 0
	fd.listFields = "*"
	fd.listFieldsRemove = []string{"is_delete"}
	fd.listOrder = "id desc"
	fd.listDataUrl = ""
	fd.listPage = 1
	fd.listPageSize = 20
	fd.listCondition = [][]interface{}{}
	fd.buttons = map[string]Button{
		"add": {
			ButtonName: "新增",
			Action:     "add",
			ActionType: 2,
			ActionUrl:  "add",
			Class:      "def",
			Icon:       "layui-icon-add-circle",
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
	fd.listRightBtns = []string{"edit", "disable", "enable", "delete"}
	fd.listTopBtns = []string{"add"}
	//重置列表项
	fd.listColumns = []ListColumn{
		{
			fd.tbPK, "ID", "text", nil, nil,
		},
	}
	fd.listColumnsStyles = map[string]interface{}{}
	fd.listSearchFields = []ListSearchField{}
	fd.listBatchAction = false

	fd.addDataUrl = ""
	fd.addTplName = "add.html"

	fd.findDataUrl = "find"
	fd.editDataUrl = ""
	fd.editTplName = "edit.html"

	fd.formFields = []FormField{}
	fd.formFieldKeys = []interface{}{}
	fd.formSubmitTitle = ""
	fd.formSubmitHide = 0
	fd.uploadImageUrl = "upload"
}

// SetStyle 设置样式
func (fd *PageData) SetStyle(style string) {
	fd.PageStyle = fd.PageStyle + style
}

// SetTitle 设置title
func (fd *PageData) SetTitle(tit string) {
	fd.title = tit
}

// SetListTplName 设置列表模板名称
func (fd *PageData) SetListTplName(tit string) {
	fd.listTplName = tit
}

// SetTbName 设置列表模板名称
func (fd *PageData) SetTbName(tbName string) {
	fd.tbName = tbName
}

// SetPK 设置数据表主键
func (fd *PageData) SetPK(pk string) {
	fd.tbPK = pk
}

func (fd *PageData) SetIsAutoID(isauto bool) {
	fd.isAutoID = isauto
}

// SetInsertAutoFields 设置插入时自动赋值字段
func (fd *PageData) SetInsertAutoFields(fields ...string) {
	fd.insertAutoFields = []string{}
	fd.insertAutoFields = append(fd.insertAutoFields, fields...)
}

// SetUpdateAutoFields 设置更新时自动赋值字段
func (fd *PageData) SetUpdateAutoFields(fields ...string) {
	fd.updateAutoFields = []string{}
	fd.updateAutoFields = append(fd.updateAutoFields, fields...)
}

// SetDeleteField 设置删除标记字段
func (fd *PageData) SetDeleteField(field string) {
	fd.deleteField = field
}

// SetPageNotice 设置列表模板名称
func (fd *PageData) SetPageNotice(str string) {
	fd.pageNotice = str
}

// SetPageName 设置结构体名称(控制器名称)
func (fd *PageData) SetPageName(str string) {
	fd.pageName = str
}

// SetActionName 设置action名称(动作名称)
func (fd *PageData) SetActionName(str string) {
	fd.actionName = str
}

// PageTabAdd 增加一个Tab选项卡
func (fd *PageData) PageTabAdd(title, href string) {
	pt := pageTab{
		Href:  href,
		Title: title,
	}
	fd.pageTabs = append(fd.pageTabs, pt)
}

// SetPageTabSelect 设置Tab选项卡的选中
func (fd *PageData) SetPageTabSelect(index int) {
	fd.pageTabSelect = index
}

// SetListFields 设置列表查询数据库字段
func (fd *PageData) SetListFields(fields string) {
	fd.listFields = fields
}

// SetListFieldsRemove 设置列表需要去掉的字段数据
func (fd *PageData) SetListFieldsRemove(fields ...string) {
	fd.listFieldsRemove = fields
}

// SetListOrder 设置列表查询排序方式
func (fd *PageData) SetListOrder(order string) {
	fd.listOrder = order
}

// SetListDataURL 设置列表查询数据地址，默认为空，代表当前地址
func (fd *PageData) SetListDataURL(url string) {
	fd.listDataUrl = url
}

// SetListPageSize 设置列表页分页大小
func (fd *PageData) SetListPageSize(size int) {
	fd.listPageSize = size
}

// SetListPage 设置列表页码
func (fd *PageData) SetListPage(page int) {
	fd.listPage = page
}

// SetListPageHide 隐藏分页
func (fd *PageData) SetListPageHide() {
	fd.listPageHide = true
}

// SetListBatchAction 开启批量操作
func (fd *PageData) SetListBatchAction(isOpen bool) {
	fd.listBatchAction = isOpen
}

// SetListCondition 设置列表页查询条件
func (fd *PageData) SetListCondition(c [][]interface{}) {
	fd.listCondition = c
}

// ListConditionAdd 增加列表页查询条件
func (fd *PageData) ListConditionAdd(c []interface{}) {
	fd.listCondition = append(fd.listCondition, c)
}

func (fd *PageData) SetButton(btnName string, btn Button) {
	fd.buttons[btnName] = btn
}

// SetListRightBtns 设置列表右侧使用的按钮
func (fd *PageData) SetListRightBtns(btns ...string) {
	fd.listRightBtns = btns
}

// ListRightBtnsClear 清除右侧按钮
func (fd *PageData) ListRightBtnsClear() {
	fd.listRightBtns = []string{}
}

// ListRightBtnsIconClear 清除右侧按钮的图标
func (fd *PageData) ListRightBtnsIconClear() {
	for _, btnName := range fd.listRightBtns {
		if btn, ok := fd.buttons[btnName]; ok {
			btn.Icon = ""
			fd.buttons[btnName] = btn
		}
	}
}

//SetButtonIcon 设置按钮图标
func (fd *PageData) SetButtonIcon(btnName, icon string) {
	if btn, ok := fd.buttons[btnName]; ok {
		btn.Icon = icon
		fd.buttons[btnName] = btn
	}
}

//SetButtonActionUrl 设置按钮链接地址 是否是追加
func (fd *PageData) SetButtonActionUrl(btnName, url string, isAddend bool) {
	if btn, ok := fd.buttons[btnName]; ok {
		if isAddend {
			if strings.Contains(btn.ActionUrl, "?") {
				btn.ActionUrl = btn.ActionUrl + "&" + url
			} else {
				btn.ActionUrl = btn.ActionUrl + "?" + url
			}
		} else {
			btn.ActionUrl = url
		}
		fd.buttons[btnName] = btn
	}
}

// ListRightBtnsAdd 右侧按钮新增
func (fd *PageData) ListRightBtnsAdd(btns ...string) {
	fd.listRightBtns = append(fd.listRightBtns, btns...)
}

// SetListTopBtns 设置列表顶部使用的列表
func (fd *PageData) SetListTopBtns(btns ...string) {
	fd.listTopBtns = btns
}

// ListTopBtnsClear 清除顶部按钮
func (fd *PageData) ListTopBtnsClear() {
	fd.listTopBtns = []string{}
}

// SetListColumns 设置列表列信息
func (fd *PageData) SetListColumns(data []ListColumn) {
	fd.listColumns = data
}

// SetListColumnStyle 设置列表列信息
func (fd *PageData) SetListColumnStyle(key, value string) {
	fd.listColumnsStyles[key] = value
}

// ListColumnAdd 列表新增一列信息
func (fd *PageData) ListColumnAdd(FieldName, ColumnName, DataType string, Data []map[string]interface{}) {
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
				util.JsonDecode(options[1], &Options)
			} else {
				//url链接写法
				opts := strings.Split(options[1], "&")
				for _, opt := range opts {
					optKv := strings.Split(opt, "=")
					Options[optKv[0]] = optKv[1]
				}
			}
		}
	}
	fd.listColumns = append(fd.listColumns, ListColumn{
		FieldName:  FieldName,
		ColumnName: ColumnName,
		DataType:   DataType,
		Data:       MapData,
		Options:    Options,
	})
}

// ListColumnClear 列表新增一列信息
func (fd *PageData) ListColumnClear() {
	fd.listColumns = []ListColumn{}
}

// ListSearchFieldAdd 增加列表搜索项
func (fd *PageData) ListSearchFieldAdd(fkey, ftype, ftitle string, defvalue interface{}, value interface{}, data []map[string]interface{},
	style string,
	expand map[string]interface{}) {

	//数据自动转map
	MapData := map[interface{}]interface{}{}
	for _, v := range data {
		MapData[v["value"]] = v["name"]
	}

	fd.listSearchFields = append(fd.listSearchFields, ListSearchField{
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
func (fd *PageData) FormFieldsAdd(fkey, ftype, ftitle, fnotice, fvalue string, ismust bool,
	dataOrOptions []map[string]interface{},
	style string,
	expand map[string]interface{}) {
	fd.formFields = append(fd.formFields, FormField{
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
func (fd *PageData) FormFieldKeysAdd(fieldKeys ...interface{}) {
	fd.formFieldKeys = append(fd.formFieldKeys, fieldKeys)
}

// SetAddDataUrl 设置新增页提交地址
func (fd *PageData) SetAddDataUrl(url string) {
	fd.addDataUrl = url
}

// SetAddTplName 设置新增页模板名称
func (fd *PageData) SetAddTplName(tplName string) {
	fd.addTplName = tplName
}

// SetFindTplName 设置新增页提交地址
func (fd *PageData) SetFindTplName(url string) {
	fd.findDataUrl = url
}

// SetEditDataUrl 设置新增页提交地址
func (fd *PageData) SetEditDataUrl(url string) {
	fd.editDataUrl = url
}

// SetEditTplName 设置新增页模板名称
func (fd *PageData) SetEditTplName(tplName string) {
	fd.editTplName = tplName
}

// SetFormData 设置表单默认数据
func (fd *PageData) SetFormData(formData gorose.Data) {
	fd.formData = formData
}

// SetFormSubmitTitle 设置表单提交按钮文字
func (fd *PageData) SetFormSubmitTitle(title string) {
	fd.formSubmitTitle = title
}

// SetFormSubmitHide 设置表单按钮隐藏
func (fd *PageData) SetFormSubmitHide() {
	fd.formSubmitHide = 1
}

// SetUploadImageUrl 设置上传图片的地址
func (fd *PageData) SetUploadImageUrl(url string) {
	fd.uploadImageUrl = url
}

func (fd *PageData) GetHttpWriter() http.ResponseWriter {
	return fd.httpW
}
func (fd *PageData) GetHttpRequest() *http.Request {
	return fd.httpR
}
func (fd *PageData) GetHttpParams() httprouter.Params {
	return fd.httpPS
}

// ================ 拓展方法 ============================================

// GetTempData 获取数据
func (fd *PageData) GetTempData(keyName string) interface{} {
	requestID := fd.GetHttpParams().ByName("requestID")
	rst, ok := pool.Gocache.Get(requestID + "_" + keyName)
	if !ok {
		return nil
	}
	return rst
}

// SetTempData 暂存数据
func (fd *PageData) SetTempData(keyName string, value interface{}) {
	requestID := fd.GetHttpParams().ByName("requestID")
	pool.Gocache.Set(requestID+"_"+keyName, value, time.Second*10)
}
