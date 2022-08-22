package builder

import (
	"net/http"
	
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
)

// NodePager 节点页接口
type NodePager interface {
	// Run 服务入口
	Run(*PageBuilder, NodePager)
	// Index 必须的handles
	Index(http.ResponseWriter, *http.Request, httprouter.Params)
	Add(http.ResponseWriter, *http.Request, httprouter.Params)
	Edit(http.ResponseWriter, *http.Request, httprouter.Params)
	Status(http.ResponseWriter, *http.Request, httprouter.Params)
	Delete(http.ResponseWriter, *http.Request, httprouter.Params)
	Upload(http.ResponseWriter, *http.Request, httprouter.Params)
	WangEditorUpload(http.ResponseWriter, *http.Request, httprouter.Params)
	ApiResult(w http.ResponseWriter, code int, msg string, data interface{})
	// NodeInit 节点
	NodeInit(*PageBuilder)
	NodeBegin(*PageBuilder) (error, int)
	NodeList(*PageBuilder) (error, int)
	NodeListCondition(*PageBuilder, [][]interface{}) ([][]interface{}, error, int)
	NodeAutoCondition(*PageBuilder, [][]interface{}) ([][]interface{}, error, int)
	NodeListData(*PageBuilder, []gorose.Data) ([]gorose.Data, error, int)
	NodeListOrm(*PageBuilder, *gorose.IOrm) (error, int)
	NodeCheckAuth(*PageBuilder, string, int) (bool, error)
	NodeForm(*PageBuilder, int64) (error, int)
	NodeFormData(*PageBuilder, gorose.Data, int64) (gorose.Data, error, int)
	NodeSaveData(*PageBuilder, gorose.Data, map[string]interface{}) (map[string]interface{}, error, int)
	NodeAutoData(*PageBuilder, map[string]interface{}, string) (map[string]interface{}, error, int)
	NodeSaveSuccess(*PageBuilder, map[string]interface{}, int64) (bool, error, int)
	NodeAddSuccess(*PageBuilder, map[string]interface{}, int64) (bool, error, int)
	NodeUpdateSuccess(*PageBuilder, map[string]interface{}, int64) (bool, error, int)
	NodeStatusSuccess(*PageBuilder, int64, int64) (error, int)
	NodeDeleteBefore(*PageBuilder, int64) (error, int)
	NodeDeleteSuccess(*PageBuilder, int64, string) (error, int)
	// ActDisplay 渲染方法
	ActDisplay(Displayer, *PageBuilder) (string, error)
	ActShow(http.ResponseWriter, Displayer, *PageBuilder)
}
