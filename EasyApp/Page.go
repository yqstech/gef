/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description: 定义基础NodeFlow接口
 * @File: Inodes
 * @Version: 1.0.0
 * @Date: 2021/10/14 5:57 下午
 */

package EasyApp

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
)

// Page 定义基页面
type Page struct {
}

// PageAction 页面方法类型
type PageAction func(*PageData, http.ResponseWriter, *http.Request, httprouter.Params)

// Run 页面执行入口，所有页面都继承此方法
func (nf Page) Run(activePage *AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//初始化页面数据，内置的一些方法名
	pageData := PageData{
		RequestID:  ps.ByName("requestID"),
		ActivePage: *activePage,
		ActionList: map[string]PageAction{
			"index":              (*activePage).Index, //内置调用的这些方法，是需要定义到接口里的，确保子类都继承/实现了此方法
			"add":                (*activePage).Add,
			"edit":               (*activePage).Edit,
			"status":             (*activePage).Status,
			"delete":             (*activePage).Delete,
			"upload":             (*activePage).Upload,
			"wang_editor_upload": (*activePage).WangEditorUpload,
		},
		httpW:  w,
		httpR:  r,
		httpPS: ps,
	}
	//数据重置
	pageData.DataReset()

	//页面初始化节点
	//子类重写可重新设置action列表，可以修改pageData数据
	(*activePage).PageInit(&pageData)

	//获取action名称，从action列表里查找action方法并执行
	actionName := ps.ByName("actionName")
	if pageAction, ok := pageData.ActionList[actionName]; ok {
		pageAction(&pageData, w, r, ps)
	} else {
		nf.Empty(w, r, ps)
	}
}

// ================ 渲染 ============================================

// ActShow 渲染并输出
func (nf Page) ActShow(w http.ResponseWriter, tpl Template, pageData *PageData) {
	//模板渲染得到html
	html, err := nf.ActDisplay(tpl, pageData)
	if err != nil {
		logger.Error(err.Error())
		fmt.Fprint(w, "模板解析出错")
		return
	}
	fmt.Fprint(w, html)
}

// ActDisplay 渲染模板
func (nf Page) ActDisplay(tpl Template, pageData *PageData) (string, error) {
	//流程数据合并到模板中
	tpl.PageData2Display(pageData)
	//渲染模板
	return tpl.Display()
}

// Empty 空操作方法，报404错误
func (nf Page) Empty(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	nf.ErrResult(w, r, 404, "您访问的页面不存在！", "")
}
