/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2021-10-15 19:47:32
 * @LastEditTime: 2021-10-20 11:28:25
 * @Description: NodeFlow接口
 */

package EasyApp

import (
	"net/http"

	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
)

// AppPage
// 定义应用页面接口
type AppPage interface {
	// Run 服务入口
	Run(obj *AppPage, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	// PageInit 前置操作
	PageInit(pageData *PageData)

	// Index 必须的handles
	Index(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Add(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Edit(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Status(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Delete(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Upload(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	WangEditorUpload(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ApiResult(w http.ResponseWriter, code int, msg string, data interface{})

	// NodeBegin 必须的节点
	NodeBegin(pageData *PageData) (error, int)
	NodeList(pageData *PageData) (error, int)
	NodeListCondition(*PageData, [][]interface{}) ([][]interface{}, error, int)
	NodeAutoCondition(*PageData, [][]interface{}) ([][]interface{}, error, int)
	NodeListData(pageData *PageData, data []gorose.Data) ([]gorose.Data, error, int)
	NodeListOrm(pageData *PageData, data *gorose.IOrm) (error, int)
	NodeCheckAuth(*PageData, string, int) (bool, error)
	NodeForm(pageData *PageData, id int64) (error, int)
	NodeFormData(pageData *PageData, data gorose.Data, id int64) (gorose.Data, error, int)
	NodeSaveData(pageData *PageData, formData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int)
	NodeAutoData(pageData *PageData, postData map[string]interface{}, action string) (map[string]interface{}, error, int)
	NodeSaveSuccess(pageData *PageData, postData map[string]interface{}, id int64) (bool, error, int)
	NodeAddSuccess(pageData *PageData, postData map[string]interface{}, id int64) (bool, error, int)
	NodeUpdateSuccess(pageData *PageData, postData map[string]interface{}, id int64) (bool, error, int)
	NodeStatusSuccess(*PageData, int64, int64) (error, int)
	NodeDeleteBefore(*PageData, int64) (error, int)
	NodeDeleteSuccess(*PageData, int64, string) (error, int)

	// ActDisplay 渲染方法
	ActDisplay(tpl Template, pageData *PageData) (string, error)
	ActShow(w http.ResponseWriter, tpl Template, pageData *PageData)
}
