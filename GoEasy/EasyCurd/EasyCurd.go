/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: EasyCurd
 * @Version: 1.0.0
 * @Date: 2022/3/30 9:30 下午
 */

package EasyCurd

import (
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/GoEasy/Utils/db"
	"github.com/yqstech/gef/GoEasy/Utils/util"
	"net/http"
	"strings"
)

// SelectAction 定义操作别名
const SelectAction = "select"
const FindAction = "find"
const CreateAction = "create"
const UpdateAction = "update"
const DeleteAction = "delete"

// DbModel 数据库模型
type DbModel struct {
	TableName          string //数据库名称
	Select             bool   //查询权限
	Create             bool   //新增权限
	Update             bool   //更新权限
	Delete             bool   //删除权限
	InsertId           bool   //允许插入Id
	SoftDeleteDisable  bool   //软删除是否禁用，禁用标识真删除
	CheckLogin         bool   //是否校验登录
	SelectWithDisabled bool   //是否校验登录

	Condition map[string]interface{} //查询或更新的sql条件
	Data      map[string]interface{} //新增或更新数据
	UkName    string                 //用户键字段
	PkName    string                 //主键
	Order     string                 //排序 默认id desc

	Fields        string   //公开查询字段，不填为*
	PrivateFields []string //私密的数据库字段，不对外展示，查询出的数据做删除处理
	LockFields    []string //锁定字段，锁定的字段不许修改，假如是更新操作，删除此字段

	FmtRule map[string]FmtRuleFunc //对定义某个键的数据值进行格式化

	ResultExtend map[string]interface{} //返回数据拓展数据

	CompleteData CompleteDataFunc //完善select单项或者find单项数据

	CompleteResultData CompleteResultDataFunc //对查询结果进行转换
}

// EasyCurd 增删改查处理对象
type EasyCurd struct {
	Page     int     //第几页
	PageSize int     //每页数量
	DbModel  DbModel //数据库模型
	Action   string  //当前操作
}

func (that *EasyCurd) Select(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.Action = SelectAction
	if !that.DbModel.Select {
		//校验权限
		that.ApiResult(w, 101, "无权操作！", nil)
		return
	}
	//查询数据
	conn := that.getOrm(ps)
	if that.DbModel.Fields == "" {
		that.DbModel.Fields = "*"
	}
	if that.DbModel.Order == "" {
		that.DbModel.Order = "id desc"
	}
	data, err := conn.Fields(that.DbModel.Fields).Limit(that.PageSize).Page(that.Page).Order(that.DbModel.Order).Get()
	if err != nil {
		logger.Error(err.Error(), conn.LastSql())
		that.ApiResult(w, 500, "数据库操作错误！", nil)
		return
	}
	//logger.Debug(conn.LastSql(),data)
	total, err := that.getOrm(ps).Count()
	if err != nil {
		logger.Error(err.Error(), conn.LastSql())
		that.ApiResult(w, 500, "数据库操作错误！", nil)
		return
	}
	//格式化数据
	for k, v := range data {
		data[k] = that.FmtData(v)
		if that.DbModel.CompleteData != nil {
			data[k] = that.DbModel.CompleteData(data[k])
		}
	}

	resultData := map[string]interface{}{
		"data":    data,
		"total":   total,
		"_extend": that.DbModel.ResultExtend,
	}
	if that.DbModel.CompleteResultData != nil {
		that.ApiResult(w, 200, "success", that.DbModel.CompleteResultData(resultData, SelectAction))
		return
	}
	that.ApiResult(w, 200, "success", resultData)
}

func (that *EasyCurd) Find(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.Action = FindAction
	if !that.DbModel.Select {
		//校验权限
		that.ApiResult(w, 101, "无权操作！", nil)
		return
	}
	//查询数据
	conn := that.getOrm(ps)
	if that.DbModel.Fields == "" {
		that.DbModel.Fields = "*"
	}
	if that.DbModel.Order == "" {
		that.DbModel.Order = "id desc"
	}
	data, err := conn.Fields(that.DbModel.Fields).Order(that.DbModel.Order).First()
	if err != nil {
		logger.Error(err.Error())
		that.ApiResult(w, 500, "数据库操作错误！", nil)
		return
	}
	if data != nil {
		//格式化数据
		data = that.FmtData(data)
		if that.DbModel.CompleteData != nil {
			data = that.DbModel.CompleteData(data)
		}
	}
	//返回数据
	resultData := map[string]interface{}{
		"data":       data,
		"_extend":    that.DbModel.ResultExtend,
		"_time":      util.TimeNow(),
		"_unix_time": util.Str2UnixTime(util.TimeNow()),
	}
	if that.DbModel.CompleteResultData != nil {
		that.ApiResult(w, 200, "success", that.DbModel.CompleteResultData(resultData, FindAction))
		return
	}
	that.ApiResult(w, 200, "success", resultData)
}
func (that *EasyCurd) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.Action = CreateAction
	if !that.DbModel.Create {
		//校验权限
		that.ApiResult(w, 101, "无权操作！", nil)
		return
	}
	//确认数据
	that.verifyData(false)

	//插入数据
	conn := that.getOrm(ps)
	insertId, err := conn.InsertGetId(that.DbModel.Data)
	if err != nil {
		logger.Error(err.Error())
		that.ApiResult(w, 500, "插入数据出错！", nil)
		return
	}
	if insertId == 0 {
		that.ApiResult(w, 500, "插入数据失败！", nil)
		return
	}

	//返回结果
	resultData := map[string]interface{}{
		"id":         insertId,
		"_extend":    that.DbModel.ResultExtend,
		"_time":      util.TimeNow(),
		"_unix_time": util.Str2UnixTime(util.TimeNow()),
	}
	that.ApiResult(w, 200, "success", resultData)

}
func (that *EasyCurd) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.Action = UpdateAction
	if !that.DbModel.Update {
		//校验权限
		that.ApiResult(w, 101, "无权操作！", nil)
		return
	}

	//格式化要更新的数据
	that.verifyData(true)

	//更新数据
	conn := that.getOrm(ps)
	update, err := conn.Update(that.DbModel.Data)
	if err != nil {
		logger.Error(err.Error())
		that.ApiResult(w, 500, "更新出错！", nil)
		return
	}
	if update == 0 {
		that.ApiResult(w, 201, "数据未修改！", nil)
		return
	}
	//返回结果
	resultData := map[string]interface{}{
		"update":     update,
		"_extend":    that.DbModel.ResultExtend,
		"_time":      util.TimeNow(),
		"_unix_time": util.Str2UnixTime(util.TimeNow()),
	}

	that.ApiResult(w, 200, "success", resultData)

}
func (that *EasyCurd) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.Action = DeleteAction
	if !that.DbModel.Delete {
		//校验权限
		that.ApiResult(w, 101, "无权操作！", nil)
		return
	}
	conn := that.getOrm(ps)
	if !that.DbModel.SoftDeleteDisable {
		//使用软删除
		update, err := conn.Update(map[string]interface{}{
			"update_time": util.TimeNow(),
			"is_delete":   1,
		})
		logger.Alert(conn.LastSql())
		if err != nil {
			logger.Error(err.Error())
			that.ApiResult(w, 500, "删除出错！", nil)
			return
		}
		if update == 0 {
			that.ApiResult(w, 500, "删除失败！", nil)
			return
		}
		//返回结果
		resultData := map[string]interface{}{
			"delete":     update, //数量
			"_extend":    that.DbModel.ResultExtend,
			"_time":      util.TimeNow(),
			"_unix_time": util.Str2UnixTime(util.TimeNow()),
		}

		that.ApiResult(w, 200, "success", resultData)
	} else {
		delete, err := conn.Delete()
		if err != nil {
			logger.Error(err.Error(), conn.LastSql())
			that.ApiResult(w, 500, "删除出错！", nil)
			return
		}
		if delete == 0 {
			that.ApiResult(w, 500, "删除失败！", nil)
			return
		}
		//返回结果
		resultData := map[string]interface{}{
			"delete":     delete, //数量
			"_extend":    that.DbModel.ResultExtend,
			"_time":      util.TimeNow(),
			"_unix_time": util.Str2UnixTime(util.TimeNow()),
		}
		that.ApiResult(w, 200, "success", resultData)
	}
}
func (that *EasyCurd) Error(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	that.ApiResult(w, 404, "未知操作！", nil)
	return
}

//初始化一个数据库连接,并且格式化条件
func (that *EasyCurd) getOrm(ps httprouter.Params) gorose.IOrm {
	conn := db.New()
	conn = conn.Table(that.DbModel.TableName)
	for key, value := range that.DbModel.Condition {
		switch value.(type) {
		//数组
		case []interface{}:
			if key[0:1] == "_" {
				//下划线开头 则 value 是二级数组
				var orps [][]interface{}
				for _, p := range value.([]interface{}) {
					switch p.(type) {
					case []interface{}:
						orps = append(orps, p.([]interface{}))
					default:
						orps = append(orps, []interface{}{1, 2})
					}
				}
				if len(orps) > 0 {
					conn.Where(func() {
						for i, orp := range orps {
							if i == 0 {
								conn.Where(orp...)
							} else {
								conn.OrWhere(orp...)
							}
						}
					})
				}
			} else {
				//一维数组 key value1 value2
				w := []interface{}{key}
				for _, v := range value.([]interface{}) {
					w = append(w, v)
				}
				// in 处理
				if util.Interface2String(w[1]) == "in" {
					if len(w) > 2 {
						switch w[2].(type) {
						case string:
							ids := strings.Split(w[2].(string), ",")
							var dataIds []interface{}
							for _, x := range ids {
								dataIds = append(dataIds, x)
							}
							conn = conn.WhereIn(key, dataIds)
						default:
							conn = conn.Where(w...)
						}
					}
				} else {
					conn = conn.Where(w...)
				}
			}
		default:
			//单项，key = value
			conn = conn.Where(key, value)
		}
	}
	return conn
}

// FmtData 格式化要更新或新增的数据
func (that *EasyCurd) verifyData(isUpdate bool) {
	for k, v := range that.DbModel.Data {
		//数组转为json数据
		that.DbModel.Data[k] = util.Array2String(v)
	}
	//设置默认主键名称
	if that.DbModel.PkName == "" {
		that.DbModel.PkName = "id"
	}
	//不允许插入ID时
	if !that.DbModel.InsertId {
		if _, ok := that.DbModel.Data[that.DbModel.PkName]; ok {
			delete(that.DbModel.Data, that.DbModel.PkName)
		}
	}
	if isUpdate {
		//更新操作，需要在更新的数据里删除锁定的字段
		for _, f := range that.DbModel.LockFields {
			if _, ok := that.DbModel.Data[f]; ok {
				delete(that.DbModel.Data, f)
			}
		}
	}
}

// FmtData 格式化展示数据
func (that *EasyCurd) FmtData(data gorose.Data) gorose.Data {
	for key, value := range data {
		//删除私有字段，不对外展示
		for _, f := range that.DbModel.PrivateFields {
			if key == f {
				//logger.Alert("删除字段", key)
				delete(data, key)
				//结束本次最外层循环
				goto Eed
			}
		}
		//对值进行格式化操作
		if fn, ok := that.DbModel.FmtRule[key]; ok {
			value = fn(value)
		}
		//保存最新值
		data[key] = value
	Eed:
	}
	return data
}

type ResultData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (that *EasyCurd) ApiResult(w http.ResponseWriter, code int, msg string, data interface{}) {
	resultData := ResultData{Code: code, Msg: msg, Data: data}
	fmt.Fprint(w, util.JsonEncode(resultData))
}
