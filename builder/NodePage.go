/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: NodePage
 * @Version: 1.0.0
 * @Date: 2022/8/22 09:59
 */

package builder

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/builder/adminTemplates"
	"io/fs"
	"net/http"
)

// NodePage 节点页的实现
type NodePage struct {
	NodePage        NodePager                    //当前页面对象
	NodePageActions map[string]httprouter.Handle //当前页面方法对象集
	PageBuilder     *PageBuilder                 //页面构建器指针
}

// Run 节点页面运行入口
// 初始化节点页面数据
func (that *NodePage) Run(pageBuilder *PageBuilder, nodePage NodePager) {
	//记录当前页面
	that.NodePage = nodePage
	//记录构建器地址
	that.PageBuilder = pageBuilder
	//初始化内置的方法集
	that.NodePageActions = map[string]httprouter.Handle{
		"index":              nodePage.Index,
		"add":                nodePage.Add,
		"edit":               nodePage.Edit,
		"status":             nodePage.Status,
		"delete":             nodePage.Delete,
		"upload":             nodePage.Upload,
		"wang_editor_upload": nodePage.WangEditorUpload,
	}
	// NodeInit 调用页面初始化节点
	// 在这里可以注册自定义方法
	that.NodePage.NodeInit(pageBuilder)

	// 查找页面方法
	actionName := pageBuilder.HttpParams.ByName("actionName")
	if pageAction, ok := that.NodePageActions[actionName]; ok {
		//执行页面方法
		pageAction(pageBuilder.HttpResponseWriter, pageBuilder.HttpRequest, pageBuilder.HttpParams)
	} else {
		//执行空白方法
		that.Empty(pageBuilder.HttpResponseWriter, pageBuilder.HttpRequest, pageBuilder.HttpParams)
	}
}

// Empty 空操作方法，报404错误
func (that *NodePage) Empty(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.ErrResult(w, r, 404, "您访问的页面不存在！", "")
}

// ActShow 渲染器+构建器，执行页面渲染并输出
func (that *NodePage) ActShow(w http.ResponseWriter, display Displayer, pageBuilder *PageBuilder) {
	//页面渲染成html
	html, err := that.ActDisplay(display, pageBuilder)
	if err != nil {
		logger.Error(err.Error())
		fmt.Fprint(w, "模板解析出错!")
		return
	}
	//输出页面
	_, err = fmt.Fprint(w, html)
	if err != nil {
		return
	}
}

// ActDisplay 渲染器+构建器，执行页面渲染
func (that *NodePage) ActDisplay(display Displayer, pageBuilder *PageBuilder) (string, error) {
	//设置模板
	for _, FilesAdd := range adminTemplates.AdminTemplatesAdd {
		fd, err := fs.ReadDir(FilesAdd, "admin")
		if err == nil {
			if len(fd) > 0 {
				display.Templates = append(display.Templates, TemplateParseFS{
					Fsys:     FilesAdd,
					Patterns: []string{"admin/*.html"},
				})
			}
		}
	}
	//构建器 导出模板格式的数据
	data := pageBuilder.TemplateData()
	for k, v := range data {
		//渲染器设置参数
		if value, ok := display.DisplayData[k]; !ok || value == "" {
			display.SetDate(k, v)
		}
	}
	//渲染器 渲染生成html
	return display.Display()
}
