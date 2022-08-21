/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Handles
 * @Version: 1.0.0
 * @Date: 2021/12/7 8:23 下午
 */

package EasyApp

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/GoEasy/Models"
	"github.com/yqstech/gef/GoEasy/Utils/db"
	"github.com/yqstech/gef/GoEasy/Utils/util"
	"github.com/yqstech/gef/config"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// ====================  默认内置的操作方法 ============================

//
// Index
//  @Description: 默认Index方法,即列表页
//
func (nf Page) Index(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageData.httpW = w
	pageData.httpR = r
	pageData.httpPS = ps
	pageData.actionName = "列表"
	pageData.listPageHide = false
	//! pageData.ActivePage
	if pageData.ActivePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}
	//logger.Debug("NodeBegin 调用前", pageData)
	//! NodeBegin()
	err, code := pageData.ActivePage.NodeBegin(pageData)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	//logger.Debug("NodeBegin 调用后", pageData)
	err, code = pageData.ActivePage.NodeList(pageData)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}

	//! 校验权限隐藏不可用按钮
	pageName := ps.ByName("pageName")                     //结构名
	accountID := util.String2Int(ps.ByName("account_id")) //账户ID
	btnHide := []string{}                                 //需要隐藏的按钮
	btnShow := []string{}
	for _, btnName := range pageData.listRightBtns {
		//根据按钮名称查信息
		if btnInfo, ok := pageData.buttons[btnName]; ok {
			//得到按钮rule
			btnRule := ""
			//统一去掉开头的斜杠
			if btnInfo.Action[0:1] == "/" {
				btnInfo.Action = btnInfo.Action[1:]
			}
			ruleArr := strings.Split(btnInfo.Action, "/")
			if len(ruleArr) > 1 {
				//有多项参数，取最后两项作为model和action
				btnRule = "/" + ruleArr[len(ruleArr)-2] + "/" + ruleArr[len(ruleArr)-1]
			} else {
				//仅有一项，则model未当前页
				btnRule = "/" + pageName + "/" + btnInfo.Action
			}
			auth, err := pageData.ActivePage.NodeCheckAuth(pageData, btnRule, accountID)
			if err != nil {
				logger.Error(err.Error())
				nf.ApiResult(w, 120, "校验权限失败！", nil)
				return
			}
			if !auth {
				btnHide = append(btnHide, btnInfo.Action)
			} else {
				btnShow = append(btnShow, btnInfo.Action)
			}
		}
	}
	//清空右侧按钮，不展示操作列
	if len(btnShow) == 0 {
		pageData.listRightBtns = []string{}
	}
	if r.Method == "POST" {
		//校验是否设置数据表
		if pageData.tbName == "" {
			nf.ApiResult(w, 120, "数据表初始失败！", nil)
			return
		}
		action := util.PostValue(r, "action")
		if action == "fastUpdate" {
			PostData := util.PostJson(r, "formFields")
			//支持快速修改的组件名称
			dataTypes := []interface{}{"switch", "select", "input"}
			for _, listColumn := range pageData.listColumns {
				if util.IsInArray(listColumn.DataType, dataTypes) {
					//检索固定的类型
					if SetValue, ok := PostData[listColumn.FieldName]; ok {
						//检索到匹配的字段，并获取到设置的值
						if setId, ok2 := PostData["id"]; ok2 {
							//获取到id
							conn := db.New()
							_, err := conn.Table(pageData.tbName).
								Where(pageData.tbPK, setId).
								Update(map[string]interface{}{
									listColumn.FieldName: SetValue,
								})
							if err != nil {
								logger.Error(err.Error(), conn.LastSql())
								nf.ApiResult(w, 500, "修改数据出现错误！", "")
								return
							}
							ok, err, code := pageData.ActivePage.NodeUpdateSuccess(pageData, PostData, int64(util.String2Int(util.Interface2String(setId))))
							if err != nil {
								nf.ErrResult(w, r, code, err.Error(), nil)
								return
							}

							ok, err, code = pageData.ActivePage.NodeSaveSuccess(pageData, PostData, int64(util.String2Int(util.Interface2String(setId))))
							if err != nil {
								nf.ErrResult(w, r, code, err.Error(), nil)
								return
							}
							if ok {
								nf.ApiResult(w, 200, "修改成功！", "success")
							}
							return
						}

					}
				}
			}
			nf.ApiResult(w, 120, "操作失败！", nil)
			return
		}
		//初始化查询条件（form表单）
		search := util.PostJson(r, "search")
		for k, v := range search {
			if util.Interface2String(v) != "" {
				pageData.listCondition = append(pageData.listCondition, []interface{}{k, v})
			}
		}
		//格式化列表查询条件
		pageData.listCondition, err, code = pageData.ActivePage.NodeListCondition(pageData, pageData.listCondition)
		if err != nil {
			nf.ApiResult(w, code, err.Error(), nil)
			return
		}
		//自动添加查询条件（时间、多商户等）
		pageData.listCondition, err, code = pageData.ActivePage.NodeAutoCondition(pageData, pageData.listCondition)
		if err != nil {
			nf.ApiResult(w, code, err.Error(), nil)
			return
		}
		//!开始查询
		//添加查询条件
		conn := db.New().Table(pageData.tbName)
		for _, v := range pageData.listCondition {
			conn = conn.Where(v...)
		}
		//子类中可修改orm对象，追加条件等信息
		err, code = pageData.ActivePage.NodeListOrm(pageData, &conn)
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		//! 执行统计和分页查询
		total, err := conn.Count()
		if err != nil {
			logger.Error(err.Error())
			nf.ApiResult(w, 500, "系统运行出错！", "")
			logger.Alert(conn.LastSql())
			return
		}
		page := util.PostValue(r, "page")
		if page != "" {
			pageData.listPage = util.String2Int(page)
		}

		data, err := conn.Fields(pageData.listFields).Limit(pageData.listPageSize).Page(pageData.listPage).Order(pageData.listOrder).Get()
		if err != nil {
			logger.Error(err.Error())
			nf.ApiResult(w, 500, "系统运行出错！", "")
			logger.Alert(conn.LastSql())
			return
		}
		//隐藏按钮列表和展示按钮列表任意一个为空就没意义
		//展示按钮列表为空则清空listRightBtns
		if len(btnHide) > 0 && len(btnShow) > 0 {
			for k, _ := range data {
				for _, btnName := range btnHide {
					data[k]["btn_"+btnName] = "hide"
				}
			}
		}

		//!自定义列表数据操作
		data, err, code = pageData.ActivePage.NodeListData(pageData, data)
		if err != nil {
			nf.ApiResult(w, code, err.Error(), nil)
			return
		}
		//!最后删除敏感字段数据
		for _, removeField := range pageData.listFieldsRemove {
			for k, v := range data {
				if _, ok := v[removeField]; ok {
					delete(data[k], removeField)
				}
			}
		}
		nf.ApiResult(w, 200, "success", map[string]interface{}{"total": total, "data": data})
		return
	}
	if pageData.listDataUrl == "" {
		pageData.listDataUrl = r.RequestURI
	}
	//logger.Info("ActShow前", pageData)
	nf.ActShow(w, Template{
		TplName: pageData.listTplName,
	}, pageData)
}

//
// Add
//  @Description: 默认新增方法
//
func (nf Page) Add(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pageData.httpW = w
	pageData.httpR = r
	pageData.httpPS = ps
	pageData.actionName = "新增"
	//! pageData.ActivePage
	if pageData.ActivePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}
	//! NodeBegin()
	err, code := pageData.ActivePage.NodeBegin(pageData)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}

	//! NodeForm() 节点
	err, code = pageData.ActivePage.NodeForm(pageData, 0)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}

	//! NodeFormData() 数据转换节点
	pageData.formData, err, code = pageData.ActivePage.NodeFormData(pageData, gorose.Data{"id": int64(0)}, 0)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	//提交信息
	if r.Method == "POST" {
		//获取提交的原始信息
		PostData := util.PostJson(r, "formFields")
		//检查多余字段
		PostData, err := nf.postDataCheckSpare(pageData, PostData)
		if err != nil {
			nf.ApiResult(w, 133, err.Error(), "")
			return
		}
		//检查必传字段
		PostData, err = nf.postDataCheckMust(pageData, PostData)
		if err != nil {
			nf.ApiResult(w, 134, err.Error(), "")
			return
		}
		//保存数据前操作
		PostData, err, code = pageData.ActivePage.NodeSaveData(pageData, pageData.formData, PostData)
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if code == -1 {
			return
		}
		//数据自动完成
		PostData, err, code = pageData.ActivePage.NodeAutoData(pageData, PostData, "add")
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		insertId := int64(0)
		insertRst := int64(0)
		if id, ok := PostData[pageData.tbPK]; ok {
			//自定义ID直接插入
			insertId = int64(util.String2Int(util.Interface2String(id)))
			insertRst, err = db.New().Table(pageData.tbName).Insert(PostData)
		} else {
			//使用数据库自增ID
			insertId, err = db.New().Table(pageData.tbName).InsertGetId(PostData)
			insertRst = insertId
		}
		if err != nil {
			logger.Error(err.Error())
			nf.ApiResult(w, 500, "插入数据出现错误！", "")
			return
		}
		if insertRst == 0 {
			logger.Error(err.Error())
			nf.ApiResult(w, 500, "插入数据失败！", "")
			return
		}

		ok, err, code := pageData.ActivePage.NodeAddSuccess(pageData, PostData, insertId)
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}

		ok, err, code = pageData.ActivePage.NodeSaveSuccess(pageData, PostData, insertId)
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if ok {
			nf.ApiResult(w, 200, "插入数据成功!", "success")
		}
		return
	}
	if pageData.addDataUrl == "" {
		pageData.addDataUrl = r.RequestURI
	}
	nf.ActShow(w, Template{
		TplName: pageData.addTplName,
	}, pageData)
}

//检查是否有多余项（提交的字段在formFields内不存在）
func (nf Page) postDataCheckSpare(pageData *PageData, PostData map[string]interface{}) (map[string]interface{}, error) {
	//查询所有表单项的KEY
	for _, formField := range pageData.formFields {
		pageData.formFieldKeys = append(pageData.formFieldKeys, formField.Key)
	}
	for postKey, _ := range PostData {
		if !util.IsInArray(postKey, pageData.formFieldKeys) {
			//删除未定义的下划线开头的字段
			if postKey[0:1] == "_" {
				delete(PostData, postKey)
				continue
			}
			return nil, errors.New("字段" + postKey + "无效！")
		}
	}
	return PostData, nil
}

//检查必传项目必须上传
func (nf Page) postDataCheckMust(pageData *PageData, PostData map[string]interface{}) (map[string]interface{}, error) {
	mustFieldVlues := map[string][]interface{}{}
	mustFieldName := map[string]string{}
	mustFieldType := map[string]string{}
	for _, formField := range pageData.formFields {
		if formField.IsMust {
			mustFieldName[formField.Key] = formField.Title
			mustFieldType[formField.Key] = formField.Type
			if formField.Data == nil {
				//数据为空，判断数据不能为空
				mustFieldVlues[formField.Key] = []interface{}{}
			} else {
				//有待选数据，判断必是其一
				mustFieldVlues[formField.Key] = []interface{}{}

				//传入数据的格式修改为[]map[string]interface{} {name:"",value:""}
				for _, dvalue := range formField.Data {
					mustFieldVlues[formField.Key] = append(mustFieldVlues[formField.Key],
						util.Interface2String(dvalue["value"]))
				}
			}
		}
	}

	for postKey, postValue := range PostData {
		if options, ok := mustFieldVlues[postKey]; ok {
			//没有待选数据的，判断不可为空
			if len(options) == 0 {
				if postValue == "" {
					return nil, errors.New(mustFieldName[postKey] + "不可为空！")
				}
			} else if mustFieldType[postKey] == "checkbox" || mustFieldType[postKey] == "tags" {
				//checkbox单独对比xzl
				if len(postValue.([]interface{})) == 0 {
					return nil, errors.New("请勾选" + mustFieldName[postKey] + "！")
				}
			} else {
				if !util.IsInArray(util.Interface2String(postValue), options) {
					logger.Alert(mustFieldName[postKey], postValue, options)
					return nil, errors.New("请选择" + mustFieldName[postKey] + "！")
				}
			}
		}
	}
	return PostData, nil
}

//
// Edit
//  @Description: 默认编辑数据方法
//
func (nf Page) Edit(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pageData.httpW = w
	pageData.httpR = r
	pageData.httpPS = ps

	pageData.actionName = "修改"
	//! pageData.ActivePage
	if pageData.ActivePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}

	//获取ID
	id := int64(util.String2Int(util.GetValue(r, "id")))
	if id <= 0 {
		nf.ErrResult(w, r, 103, "页面获取ID失败！", nil)
		return
	}

	//! NodeBegin()
	err, code := pageData.ActivePage.NodeBegin(pageData)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}

	//!NodeForm 初始化字段
	err, code = pageData.ActivePage.NodeForm(pageData, id)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}

	//!查询数据信息
	editCondition := [][]interface{}{}
	editCondition, err, code = pageData.ActivePage.NodeAutoCondition(pageData, editCondition)
	if err != nil {
		nf.ApiResult(w, code, err.Error(), nil)
		return
	}
	conn := db.New().Table(pageData.tbName)
	for _, v := range editCondition {
		conn = conn.Where(v...)
	}
	info, err := conn.Where(pageData.tbPK, id).First()
	if err != nil {
		logger.Error(err.Error())
		nf.ErrResult(w, r, 500, "程序出错了！", "")
		return
	}
	if info == nil {
		nf.ErrResult(w, r, 404, "数据不存在或已删除！", "")
		return
	}
	//logger.Alert("Edit查询数据", conn.LastSql())

	//原始表信息转换
	pageData.formData, err, code = pageData.ActivePage.NodeFormData(pageData, info, id)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	if r.Method == "POST" {
		//修改数据
		//获取提交的原始信息
		PostData := util.PostJson(r, "formFields")
		//检查多余字段
		PostData, err := nf.postDataCheckSpare(pageData, PostData)
		if err != nil {
			nf.ApiResult(w, 133, err.Error(), "")
			return
		}
		//检查必传字段
		PostData, err = nf.postDataCheckMust(pageData, PostData)
		if err != nil {
			nf.ApiResult(w, 134, err.Error(), "")
			return
		}
		//保存数据前操作
		PostData, err, code = pageData.ActivePage.NodeSaveData(pageData, pageData.formData, PostData)
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if code == -1 {
			return
		}
		//数据自动完成
		PostData, err, code = pageData.ActivePage.NodeAutoData(pageData, PostData, "edit")
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		update, err := db.New().Table(pageData.tbName).Where(pageData.tbPK, id).Update(PostData)
		if err != nil {
			logger.Error(err.Error())
			nf.ApiResult(w, 500, "插入数据出现错误！", "")
			return
		}
		if update != 1 {
			nf.ApiResult(w, 201, "数据未修改！", "")
			return
		}
		ok, err, code := pageData.ActivePage.NodeUpdateSuccess(pageData, PostData, id)
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}

		ok, err, code = pageData.ActivePage.NodeSaveSuccess(pageData, PostData, id)
		if err != nil {
			nf.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if ok {
			nf.ApiResult(w, 200, "插入数据成功!", "success")
		}
		return
	}

	//POST地址 含有链接内含有id
	if pageData.editDataUrl == "" {
		pageData.editDataUrl = r.RequestURI
	}

	nf.ActShow(w, Template{
		TplName: pageData.editTplName,
	}, pageData)

}

//
// Status
//  @Description: 默认设置状态方法
//
func (nf Page) Status(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pageData.httpW = w
	pageData.httpR = r
	pageData.httpPS = ps
	//logger.Alert("默认Status方法")
	if pageData.ActivePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}
	status := int64(util.String2Int(util.GetValue(r, "status")))
	id := int64(util.String2Int(util.PostValue(r, "id")))

	//批量操作
	var ids []int64
	if id <= 0 {
		idsJson := r.PostFormValue("ids")
		util.JsonDecode(idsJson, &ids)
		if len(ids) == 0 {
			nf.ApiResult(w, 103, "页面获取ID失败！", nil)
			return
		}
	} else {
		ids = append(ids, id)
	}

	//! NodeBegin()
	err, code := pageData.ActivePage.NodeBegin(pageData)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}

	if status != 1 {
		status = 0
	}
	//格式化权限
	editCondition := [][]interface{}{}
	editCondition, err, code = pageData.ActivePage.NodeAutoCondition(pageData, editCondition)
	if err != nil {
		nf.ApiResult(w, code, err.Error(), nil)
		return
	}
	conn := db.New().Table(pageData.tbName)
	for _, v := range editCondition {
		conn = conn.Where(v...)
	}

	//!转一下格式
	var Ids []interface{}
	for _, v := range ids {
		Ids = append(Ids, v)
	}

	rst, err := conn.WhereIn(pageData.tbPK, Ids).Update(map[string]interface{}{
		"status": status,
	})

	if err != nil {
		logger.Error(err.Error())
		nf.ApiResult(w, 500, "修改数据出错！", nil)
		return
	}
	if rst == 0 {
		nf.ApiResult(w, 500, "数据不存在或未修改！", nil)
		return
	}
	err, code = pageData.ActivePage.NodeStatusSuccess(pageData, id, status)
	if err != nil {
		nf.ApiResult(w, code, err.Error(), nil)
		return
	}
	if code == -1 {
		return
	}
	nf.ApiResult(w, 200, "数据修改成功！", util.JsonEncode(ids))

}

//
// Delete
//  @Description: 默认删除方法
//
func (nf Page) Delete(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pageData.httpW = w
	pageData.httpR = r
	pageData.httpPS = ps

	if pageData.ActivePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}
	id := int64(util.String2Int(util.PostValue(r, "id")))

	//批量操作
	var ids []int64
	if id <= 0 {
		idsJson := r.PostFormValue("ids")
		util.JsonDecode(idsJson, &ids)
		if len(ids) == 0 {
			nf.ApiResult(w, 103, "页面获取ID失败！", nil)
			return
		}
	} else {
		ids = append(ids, id)
	}

	//! NodeBegin()
	err, code := pageData.ActivePage.NodeBegin(pageData)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}

	editCondition := [][]interface{}{}
	editCondition, err, code = pageData.ActivePage.NodeAutoCondition(pageData, editCondition)
	if err != nil {
		nf.ApiResult(w, code, err.Error(), nil)
		return
	}
	conn := db.New().Table(pageData.tbName)
	for _, v := range editCondition {
		conn = conn.Where(v...)
	}

	err, code = pageData.ActivePage.NodeDeleteBefore(pageData, id)
	if err != nil {
		nf.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	if code == -1 {
		return
	}

	//!转一下格式
	var Ids []interface{}
	for _, v := range ids {
		Ids = append(Ids, v)
	}

	if pageData.deleteField == "" {
		rst, err := conn.WhereIn(pageData.tbPK, Ids).Delete()
		if err != nil {
			logger.Error(err.Error())
			nf.ApiResult(w, 500, "操作失败！", nil)
			return
		}
		err, code = pageData.ActivePage.NodeDeleteSuccess(pageData, id, pageData.deleteField)
		if err != nil {
			nf.ApiResult(w, code, err.Error(), nil)
			return
		}
		if rst == 0 {
			nf.ApiResult(w, 500, "数据不存在或已删除！", nil)
			return
		}
		nf.ApiResult(w, 200, "操作成功！", "1")
	} else {
		rst, err := conn.WhereIn(pageData.tbPK, Ids).Update(map[string]interface{}{
			pageData.deleteField: 1,
		})
		if err != nil {
			logger.Error(err.Error())
			nf.ApiResult(w, 500, "操作失败！", nil)
			return
		}
		err, code = pageData.ActivePage.NodeDeleteSuccess(pageData, id, pageData.deleteField)
		if err != nil {
			nf.ApiResult(w, code, err.Error(), nil)
			return
		}
		if rst == 0 {
			nf.ApiResult(w, 500, "数据不存在或已删除！", nil)
			return
		}
		nf.ApiResult(w, 200, "操作成功！", "1")
	}

}

// Upload 默认上传图片处理
func (nf Page) Upload(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pageData.httpW = w
	pageData.httpR = r
	pageData.httpPS = ps
	//logger.Alert("默认Upload方法")
	rst, err := nf.doUpload(pageData, w, r, ps)
	if err != nil {
		nf.ApiResult(w, 100, err.Error(), nil)
		return
	}
	if rst != nil {
		nf.ApiResult(w, 200, "S", rst)
		return
	}
	nf.ApiResult(w, 100, "未知错误！", nil)
}
func (nf Page) WangEditorUpload(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pageData.httpW = w
	pageData.httpR = r
	pageData.httpPS = ps
	//logger.Alert("WangEditorUpload方法")
	rst, err := nf.doUpload(pageData, w, r, ps)
	if err != nil {
		data := map[string]interface{}{
			"errno":   1,
			"message": err.Error(),
			"data":    []map[string]interface{}{},
		}
		fmt.Fprint(w, util.JsonEncode(data))
		return
	}
	if rst != nil {
		data := map[string]interface{}{
			"errno": 0,
			"data": map[string]interface{}{
				"url":  rst["url"],
				"alt":  "",
				"href": "#",
			},
		}
		fmt.Fprint(w, util.JsonEncode(data))
		return
	}
	data := map[string]interface{}{
		"errno":   1,
		"message": "上传失败！",
		"data":    []map[string]interface{}{},
	}
	fmt.Fprint(w, util.JsonEncode(data))
}

//
// doUpload
//  @Description: 默认处理图片上传
//
func (nf Page) doUpload(pageData *PageData, w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]string, error) {

	//!接受文件
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, errors.New("文件上传失败(1)！")
	}
	defer file.Close()

	//获取文件后缀
	uploadFileName := fileHeader.Filename
	fileExt := path.Ext(uploadFileName)
	fileExt = strings.ToLower(fileExt)
	//获取文件大小
	fileSize := fileHeader.Size

	//允许上传的文件后缀,以.开头
	var uploadAllowExts []interface{}
	uploadExt := Models.AppConfigs{}.Value("upload_extension")
	if uploadExt != "" {
		Ext := strings.Split(uploadExt, ",")
		for _, v := range Ext {
			uploadAllowExts = append(uploadAllowExts, "."+v)
		}
	}
	if !util.IsInArray(fileExt, uploadAllowExts) {
		return nil, errors.New("不支持的文件类型！")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, errors.New("拷贝文件出错！")
	}
	//md5
	md5er := md5.New()
	md5er.Write(buf.Bytes())
	md5value := hex.EncodeToString(md5er.Sum(nil))

	//!路径数据
	//上传文件记录数据库
	uploadDBName := config.UploadTableName
	//上传文件相对目录地址
	uploadPath := config.UploadPath
	//上传文件url地址
	uploadUrl := config.UploadUrl
	//中间件给出二级地址
	uploadSubPath := ps.ByName("uploadSubPath")
	if uploadSubPath == "" {
		uploadSubPath = "/" + util.UnixTimeFormat(util.Str2UnixTime(util.TimeNow()), "200601")
	}
	//分组
	uploadGroupID := ps.ByName("uploadGroupID")
	//用户id
	uploadUserID := ps.ByName("uploadUserID")

	//重复文件截停
	oldfile, err := db.New().Table(uploadDBName).
		Where("md5", md5value).
		Where("group_id", uploadGroupID).
		Where("user_id", uploadUserID).
		Where("is_delete", 0).
		First()
	if oldfile != nil {
		//返回图片ID
		return map[string]string{
			"url": oldfile["src"].(string),
			"id":  util.Int642String(oldfile["id"].(int64)),
		}, nil
	}

	//定义文件名称
	fileName := util.MD5(util.TimeNow() + util.GenValidateCode(6))

	//!创建目录
	savePath := uploadPath + uploadSubPath
	os.MkdirAll(savePath, os.ModePerm)

	//src
	fileSrc := uploadUrl + uploadSubPath + "/" + fileName + fileExt
	//文件目录
	saveFile := uploadPath + uploadSubPath + "/" + fileName + fileExt

	//创建图片资源文件
	dst, err := os.Create(saveFile)
	if err != nil {
		util.ErrorLog(err)
		return nil, errors.New("文件上传失败(2)！")
	}
	//保存图片
	_, err = io.Copy(dst, buf)
	if err != nil {
		util.ErrorLog(err)
		return nil, errors.New("文件上传失败(3)！")
	}
	//保存图片地址
	insertId := util.GetOnlyID()
	_, err = db.New().Table(uploadDBName).Insert(map[string]interface{}{
		"id":          insertId,
		"group_id":    uploadGroupID,
		"user_id":     uploadUserID,
		"file_name":   uploadFileName,
		"src":         fileSrc,
		"path":        saveFile,
		"file_size":   fileSize,
		"ext":         fileExt,
		"md5":         md5value,
		"create_time": util.TimeNow(),
		"update_time": util.TimeNow(),
	})
	if err != nil {
		util.ErrorLog(err)
		return nil, errors.New("保存文件信息失败！")
	}
	//返回图片ID
	return map[string]string{"url": fileSrc, "id": util.Int642String(insertId)}, nil

}

// StructApiResult 前台接口标准结构
type StructApiResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ApiResult 接口形式返回方法
func (nf Page) ApiResult(w http.ResponseWriter, code int, msg string, data interface{}) {
	Apirst := StructApiResult{Code: code, Msg: msg, Data: data}
	fmt.Fprint(w, util.JsonEncode(Apirst))
}

// ErrResult 有好的返回错误信息或错误页面
func (nf Page) ErrResult(w http.ResponseWriter, r *http.Request, code int, msg string, data interface{}) {
	if r.Method == "POST" {
		nf.ApiResult(w, code, msg, data)
		return
	}
	//模板对象
	tpl := Template{
		DisplayData: map[string]interface{}{
			"error_code": code,
			"error_msg":  msg,
		},
		TplName: "error.html",
	}
	pageData := PageData{}
	tpl.PageData2Display(&pageData)
	html, err := tpl.Display()
	if err != nil {
		logger.Error(err.Error())
		fmt.Fprint(w, msg+"(code:"+util.Int2String(code)+")")
		return
	}
	fmt.Fprint(w, html)
}
