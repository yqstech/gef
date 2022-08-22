/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Handles
 * @Version: 1.0.0
 * @Date: 2021/12/7 8:23 下午
 */

package builder

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Models"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/Utils/util"
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
func (that *NodePage) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.PageBuilder.SetActionName("列表")
	//! that.NodePage
	if that.NodePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}
	
	//! NodeBegin()
	err, code := that.NodePage.NodeBegin(that.PageBuilder)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	err, code = that.NodePage.NodeList(that.PageBuilder)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	//! 校验权限隐藏不可用按钮
	pageName := ps.ByName("pageName")                     //结构名
	accountID := util.String2Int(ps.ByName("account_id")) //账户ID
	var btnHide []string                                  //需要隐藏的按钮
	var btnShow []string
	for _, btnName := range that.PageBuilder.listRightBtns {
		//根据按钮名称查信息
		if btnInfo, ok := that.PageBuilder.buttons[btnName]; ok {
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
			auth, err := that.NodePage.NodeCheckAuth(that.PageBuilder, btnRule, accountID)
			if err != nil {
				logger.Error(err.Error())
				that.ApiResult(w, 120, "校验权限失败！", nil)
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
		that.PageBuilder.ListRightBtnsClear()
	}
	if r.Method == "POST" {
		//校验是否设置数据表
		if that.PageBuilder.tbName == "" {
			that.ApiResult(w, 120, "数据表初始失败！", nil)
			return
		}
		action := util.PostValue(r, "action")
		if action == "fastUpdate" {
			PostData := util.PostJson(r, "formFields")
			//支持快速修改的组件名称
			dataTypes := []interface{}{"switch", "select", "input"}
			for _, listColumn := range that.PageBuilder.listColumns {
				if util.IsInArray(listColumn.DataType, dataTypes) {
					//检索固定的类型
					if SetValue, ok := PostData[listColumn.FieldName]; ok {
						//检索到匹配的字段，并获取到设置的值
						if setId, ok2 := PostData["id"]; ok2 {
							//获取到id
							conn := db.New()
							_, err := conn.Table(that.PageBuilder.tbName).
								Where(that.PageBuilder.tbPK, setId).
								Update(map[string]interface{}{
									listColumn.FieldName: SetValue,
								})
							if err != nil {
								logger.Error(err.Error(), conn.LastSql())
								that.ApiResult(w, 500, "修改数据出现错误！", "")
								return
							}
							ok, err, code := that.NodePage.NodeUpdateSuccess(that.PageBuilder, PostData, int64(util.String2Int(util.Interface2String(setId))))
							if err != nil {
								that.ErrResult(w, r, code, err.Error(), nil)
								return
							}
							
							ok, err, code = that.NodePage.NodeSaveSuccess(that.PageBuilder, PostData, int64(util.String2Int(util.Interface2String(setId))))
							if err != nil {
								that.ErrResult(w, r, code, err.Error(), nil)
								return
							}
							if ok {
								that.ApiResult(w, 200, "修改成功！", "success")
							}
							return
						}
						
					}
				}
			}
			that.ApiResult(w, 120, "操作失败！", nil)
			return
		}
		//初始化查询条件（form表单）
		search := util.PostJson(r, "search")
		for k, v := range search {
			if util.Interface2String(v) != "" {
				that.PageBuilder.ListConditionAdd([]interface{}{k, v})
			}
		}
		//格式化列表查询条件
		listCondition, err, code := that.NodePage.NodeListCondition(that.PageBuilder, that.PageBuilder.listCondition)
		if err != nil {
			that.ApiResult(w, code, err.Error(), nil)
			return
		}
		//自动添加查询条件（时间、多商户等）
		listCondition, err, code = that.NodePage.NodeAutoCondition(that.PageBuilder, listCondition)
		if err != nil {
			that.ApiResult(w, code, err.Error(), nil)
			return
		}
		that.PageBuilder.SetListCondition(listCondition)
		//!开始查询
		//添加查询条件
		conn := db.New().Table(that.PageBuilder.tbName)
		for _, v := range listCondition {
			conn = conn.Where(v...)
		}
		//子类中可修改orm对象，追加条件等信息
		err, code = that.NodePage.NodeListOrm(that.PageBuilder, &conn)
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		//! 执行统计和分页查询
		total, err := conn.Count()
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "系统运行出错！", "")
			logger.Alert(conn.LastSql())
			return
		}
		page := util.PostValue(r, "page")
		if page != "" {
			that.PageBuilder.SetListPage(util.String2Int(page))
		}
		
		data, err := conn.Fields(that.PageBuilder.GetListFields()).Limit(that.PageBuilder.GetListPageSize()).Page(that.PageBuilder.GetListPage()).Order(that.PageBuilder.GetListOrder()).Get()
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "系统运行出错！", "")
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
		data, err, code = that.NodePage.NodeListData(that.PageBuilder, data)
		if err != nil {
			that.ApiResult(w, code, err.Error(), nil)
			return
		}
		//!最后删除敏感字段数据
		for _, removeField := range that.PageBuilder.GetListFieldsRemove() {
			for k, v := range data {
				if _, ok := v[removeField]; ok {
					delete(data[k], removeField)
				}
			}
		}
		that.ApiResult(w, 200, "success", map[string]interface{}{"total": total, "data": data})
		return
	}
	if that.PageBuilder.GetListDataURL() == "" {
		that.PageBuilder.SetListDataURL(r.RequestURI)
	}
	
	tpl := Displayer{
		TplName: that.PageBuilder.GetListTplName(),
	}
	that.ActShow(w, tpl, that.PageBuilder)
}

//
// Add
//  @Description: 默认新增方法
//
func (that *NodePage) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//开始
	that.PageBuilder.SetActionName("新增")
	//! that.NodePage
	if that.NodePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}
	//! NodeBegin()
	err, code := that.NodePage.NodeBegin(that.PageBuilder)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	//! NodeForm() 节点
	err, code = that.NodePage.NodeForm(that.PageBuilder, 0)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	//! NodeFormData() 数据转换节点
	that.PageBuilder.formData, err, code = that.NodePage.NodeFormData(that.PageBuilder, gorose.Data{"id": int64(0)}, 0)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	//提交信息
	if r.Method == "POST" {
		//获取提交的原始信息
		PostData := util.PostJson(r, "formFields")
		//检查多余字段
		PostData, err := that.postDataCheckSpare(PostData)
		if err != nil {
			that.ApiResult(w, 133, err.Error(), "")
			return
		}
		//检查必传字段
		PostData, err = that.postDataCheckMust(PostData)
		if err != nil {
			that.ApiResult(w, 134, err.Error(), "")
			return
		}
		//保存数据前操作
		PostData, err, code = that.NodePage.NodeSaveData(that.PageBuilder, that.PageBuilder.formData, PostData)
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if code == -1 {
			return
		}
		//数据自动完成
		PostData, err, code = that.NodePage.NodeAutoData(that.PageBuilder, PostData, "add")
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		insertId := int64(0)
		insertRst := int64(0)
		if id, ok := PostData[that.PageBuilder.GetPK()]; ok {
			//自定义ID直接插入
			insertId = int64(util.String2Int(util.Interface2String(id)))
			insertRst, err = db.New().Table(that.PageBuilder.tbName).Insert(PostData)
		} else {
			//使用数据库自增ID
			insertId, err = db.New().Table(that.PageBuilder.tbName).InsertGetId(PostData)
			insertRst = insertId
		}
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "插入数据出现错误！", "")
			return
		}
		if insertRst == 0 {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "插入数据失败！", "")
			return
		}
		
		ok, err, code := that.NodePage.NodeAddSuccess(that.PageBuilder, PostData, insertId)
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		
		ok, err, code = that.NodePage.NodeSaveSuccess(that.PageBuilder, PostData, insertId)
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if ok {
			that.ApiResult(w, 200, "插入数据成功!", "success")
		}
		return
	}
	if that.PageBuilder.GetAddDataUrl() == "" {
		that.PageBuilder.SetAddDataUrl(r.RequestURI)
	}
	that.ActShow(w, Displayer{
		TplName: that.PageBuilder.GetAddTplName(),
	}, that.PageBuilder)
}

//检查是否有多余项（提交的字段在formFields内不存在）
func (that *NodePage) postDataCheckSpare(PostData map[string]interface{}) (map[string]interface{}, error) {
	//查询所有表单项的KEY
	for _, formField := range that.PageBuilder.formFields {
		that.PageBuilder.formFieldKeys = append(that.PageBuilder.formFieldKeys, formField.Key)
	}
	for postKey, _ := range PostData {
		if !util.IsInArray(postKey, that.PageBuilder.formFieldKeys) {
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
func (that *NodePage) postDataCheckMust(PostData map[string]interface{}) (map[string]interface{}, error) {
	mustFieldVlues := map[string][]interface{}{}
	mustFieldName := map[string]string{}
	mustFieldType := map[string]string{}
	for _, formField := range that.PageBuilder.formFields {
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
func (that *NodePage) Edit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//开始
	that.PageBuilder.SetActionName("修改")
	//! that.NodePage
	if that.NodePage == nil {
		logger.Error("运行出错请刷新页面！")
		return
	}
	
	//获取ID
	id := int64(util.String2Int(util.GetValue(r, "id")))
	if id <= 0 {
		that.ErrResult(w, r, 103, "页面获取ID失败！", nil)
		return
	}
	
	//! NodeBegin()
	err, code := that.NodePage.NodeBegin(that.PageBuilder)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	//!NodeForm 初始化字段
	err, code = that.NodePage.NodeForm(that.PageBuilder, id)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	//!查询数据信息
	editCondition := [][]interface{}{}
	editCondition, err, code = that.NodePage.NodeAutoCondition(that.PageBuilder, editCondition)
	if err != nil {
		that.ApiResult(w, code, err.Error(), nil)
		return
	}
	conn := db.New().Table(that.PageBuilder.tbName)
	for _, v := range editCondition {
		conn = conn.Where(v...)
	}
	info, err := conn.Where(that.PageBuilder.GetPK(), id).First()
	if err != nil {
		logger.Error(err.Error())
		that.ErrResult(w, r, 500, "程序出错了！", "")
		return
	}
	if info == nil {
		that.ErrResult(w, r, 404, "数据不存在或已删除！", "")
		return
	}
	//logger.Alert("Edit查询数据", conn.LastSql())
	
	//原始表信息转换
	that.PageBuilder.formData, err, code = that.NodePage.NodeFormData(that.PageBuilder, info, id)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	if r.Method == "POST" {
		//修改数据
		//获取提交的原始信息
		PostData := util.PostJson(r, "formFields")
		//检查多余字段
		PostData, err := that.postDataCheckSpare(PostData)
		if err != nil {
			that.ApiResult(w, 133, err.Error(), "")
			return
		}
		//检查必传字段
		PostData, err = that.postDataCheckMust(PostData)
		if err != nil {
			that.ApiResult(w, 134, err.Error(), "")
			return
		}
		//保存数据前操作
		PostData, err, code = that.NodePage.NodeSaveData(that.PageBuilder, that.PageBuilder.formData, PostData)
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if code == -1 {
			return
		}
		//数据自动完成
		PostData, err, code = that.NodePage.NodeAutoData(that.PageBuilder, PostData, "edit")
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		update, err := db.New().Table(that.PageBuilder.tbName).Where(that.PageBuilder.GetPK(), id).Update(PostData)
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "插入数据出现错误！", "")
			return
		}
		if update != 1 {
			that.ApiResult(w, 201, "数据未修改！", "")
			return
		}
		ok, err, code := that.NodePage.NodeUpdateSuccess(that.PageBuilder, PostData, id)
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		
		ok, err, code = that.NodePage.NodeSaveSuccess(that.PageBuilder, PostData, id)
		if err != nil {
			that.ErrResult(w, r, code, err.Error(), nil)
			return
		}
		if ok {
			that.ApiResult(w, 200, "插入数据成功!", "success")
		}
		return
	}
	
	//POST地址 含有链接内含有id
	if that.PageBuilder.editDataUrl == "" {
		that.PageBuilder.editDataUrl = r.RequestURI
	}
	
	that.ActShow(w, Displayer{
		TplName: that.PageBuilder.editTplName,
	}, that.PageBuilder)
	
}

//
// Status
//  @Description: 默认设置状态方法
//
func (that *NodePage) Status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//开始
	//logger.Alert("默认Status方法")
	if that.NodePage == nil {
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
			that.ApiResult(w, 103, "页面获取ID失败！", nil)
			return
		}
	} else {
		ids = append(ids, id)
	}
	
	//! NodeBegin()
	err, code := that.NodePage.NodeBegin(that.PageBuilder)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	if status != 1 {
		status = 0
	}
	//格式化权限
	editCondition := [][]interface{}{}
	editCondition, err, code = that.NodePage.NodeAutoCondition(that.PageBuilder, editCondition)
	if err != nil {
		that.ApiResult(w, code, err.Error(), nil)
		return
	}
	conn := db.New().Table(that.PageBuilder.tbName)
	for _, v := range editCondition {
		conn = conn.Where(v...)
	}
	
	//!转一下格式
	var Ids []interface{}
	for _, v := range ids {
		Ids = append(Ids, v)
	}
	
	rst, err := conn.WhereIn(that.PageBuilder.GetPK(), Ids).Update(map[string]interface{}{
		"status": status,
	})
	
	if err != nil {
		logger.Error(err.Error())
		that.ApiResult(w, 500, "修改数据出错！", nil)
		return
	}
	if rst == 0 {
		that.ApiResult(w, 500, "数据不存在或未修改！", nil)
		return
	}
	err, code = that.NodePage.NodeStatusSuccess(that.PageBuilder, id, status)
	if err != nil {
		that.ApiResult(w, code, err.Error(), nil)
		return
	}
	if code == -1 {
		return
	}
	that.ApiResult(w, 200, "数据修改成功！", util.JsonEncode(ids))
	
}

//
// Delete
//  @Description: 默认删除方法
//
func (that *NodePage) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//开始
	if that.NodePage == nil {
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
			that.ApiResult(w, 103, "页面获取ID失败！", nil)
			return
		}
	} else {
		ids = append(ids, id)
	}
	
	//! NodeBegin()
	err, code := that.NodePage.NodeBegin(that.PageBuilder)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
		return
	}
	
	editCondition := [][]interface{}{}
	editCondition, err, code = that.NodePage.NodeAutoCondition(that.PageBuilder, editCondition)
	if err != nil {
		that.ApiResult(w, code, err.Error(), nil)
		return
	}
	conn := db.New().Table(that.PageBuilder.tbName)
	for _, v := range editCondition {
		conn = conn.Where(v...)
	}
	
	err, code = that.NodePage.NodeDeleteBefore(that.PageBuilder, id)
	if err != nil {
		that.ErrResult(w, r, code, err.Error(), nil)
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
	
	if that.PageBuilder.deleteField == "" {
		rst, err := conn.WhereIn(that.PageBuilder.GetPK(), Ids).Delete()
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "操作失败！", nil)
			return
		}
		err, code = that.NodePage.NodeDeleteSuccess(that.PageBuilder, id, that.PageBuilder.deleteField)
		if err != nil {
			that.ApiResult(w, code, err.Error(), nil)
			return
		}
		if rst == 0 {
			that.ApiResult(w, 500, "数据不存在或已删除！", nil)
			return
		}
		that.ApiResult(w, 200, "操作成功！", "1")
	} else {
		rst, err := conn.WhereIn(that.PageBuilder.GetPK(), Ids).Update(map[string]interface{}{
			that.PageBuilder.deleteField: 1,
		})
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "操作失败！", nil)
			return
		}
		err, code = that.NodePage.NodeDeleteSuccess(that.PageBuilder, id, that.PageBuilder.deleteField)
		if err != nil {
			that.ApiResult(w, code, err.Error(), nil)
			return
		}
		if rst == 0 {
			that.ApiResult(w, 500, "数据不存在或已删除！", nil)
			return
		}
		that.ApiResult(w, 200, "操作成功！", "1")
	}
	
}

// Upload 默认上传图片处理
func (that *NodePage) Upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//logger.Alert("默认Upload方法")
	rst, err := that.doUpload(w, r, ps)
	if err != nil {
		that.ApiResult(w, 100, err.Error(), nil)
		return
	}
	if rst != nil {
		that.ApiResult(w, 200, "S", rst)
		return
	}
	that.ApiResult(w, 100, "未知错误！", nil)
}
func (that *NodePage) WangEditorUpload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//logger.Alert("WangEditorUpload方法")
	rst, err := that.doUpload(w, r, ps)
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
func (that *NodePage) doUpload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]string, error) {
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
func (that NodePage) ApiResult(w http.ResponseWriter, code int, msg string, data interface{}) {
	ApiRst := StructApiResult{Code: code, Msg: msg, Data: data}
	fmt.Fprint(w, util.JsonEncode(ApiRst))
}

// ErrResult 返回错误信息或错误页面
func (that NodePage) ErrResult(w http.ResponseWriter, r *http.Request, code int, msg string, data interface{}) {
	if r.Method == "POST" {
		that.ApiResult(w, code, msg, data)
		return
	}
	//创建错误页的渲染器
	displayer := Displayer{
		DisplayData: map[string]interface{}{
			"error_code": code,
			"error_msg":  msg,
		},
		TplName: "error.html",
	}
	//创建空的构建器
	pageBuilder := PageBuilder{}
	that.PageBuilder.DataReset()
	//渲染错误页
	html, err := that.ActDisplay(displayer, &pageBuilder)
	if err != nil {
		logger.Error(err.Error())
		fmt.Fprint(w, msg+"(code:"+util.Int2String(code)+")")
		return
	}
	fmt.Fprint(w, html)
}
