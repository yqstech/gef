/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Gateway
 * @Version: 1.0.0
 * @Date: 2022/8/22 21:36
 */

package EasyCurd

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wonderivan/logger"
	"github.com/yqstech/gef/Handles/commHandle"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/util"
	"net/http"
	"strings"
)

// Gateway EasyCurd操作的入口
//? 返回一个 httprouter.Handle
//? 内部自动初始化EasyCurd信息并调用增删改查方法
//? 通过识别ps里的userInfo["id"]判断是否登录
func Gateway() httprouter.Handle {
	//!创建增删改查适配器
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//? 模型名称
		modelName := util.PostValue(r, "modelName")
		//? 分页信息
		page := util.PostValueDef(r, "page", "1")
		pageSize := util.PostValueDef(r, "pageSize", "10")
		//? 查询条件
		condition := r.PostFormValue("condition")
		condition = strings.Replace(condition, "now()", util.TimeNow(), -1)
		//? 修改的数据
		data := r.PostFormValue("data")
		
		//? 查询模型信息
		modelInfo, err := db.New().Table("tb_easy_curd_models").
			Where("model_key", modelName).
			Where("is_delete", 0).
			Where("status", 1).First()
		if err != nil {
			logger.Error(err.Error())
			commHandle.Base{}.ApiResult(w, 500, "系统运行错误！", nil)
			return
		}
		if modelInfo == nil {
			commHandle.Base{}.ApiResult(w, 404, "数据模型未配置！", nil)
			return
		}
		
		//? 查询模型字段
		easyModelFields, err := db.New().Table("tb_easy_curd_models_fields").
			Where("model_id", modelInfo["id"]).
			Where("is_delete", 0).
			Get()
		if err != nil {
			logger.Error(err.Error())
			commHandle.Base{}.ApiResult(w, 500, "系统运行错误！", nil)
			return
		}
		//? 私有字段和锁定字段列表
		var PrivateFieldList = []string{"is_delete"}
		var LockFieldList []string
		for _, emf := range easyModelFields {
			if emf["is_private"].(int64) == 1 {
				PrivateFieldList = append(PrivateFieldList, emf["field_key"].(string))
			}
			if emf["is_lock"].(int64) == 1 {
				LockFieldList = append(LockFieldList, emf["field_key"].(string))
			}
		}
		//? 创建 EasyCurd 对象
		easyCurd := EasyCurd{
			Page:     util.String2Int(page),     //列表用
			PageSize: util.String2Int(pageSize), //列表项用
			DbModel: DbModel{
				ID:                 modelInfo["id"].(int64),
				TableName:          modelInfo["table_name"].(string),
				Select:             modelInfo["allow_select"].(int64) == 1,
				Create:             modelInfo["allow_create"].(int64) == 1,
				Update:             modelInfo["allow_update"].(int64) == 1,
				Delete:             modelInfo["allow_delete"].(int64) == 1,
				InsertId:           false,
				SoftDeleteDisable:  modelInfo["soft_delete_disable"].(int64) == 1,
				CheckLogin:         modelInfo["check_login"].(int64) == 1,
				SelectWithDisabled: modelInfo["select_with_disabled"].(int64) == 1,
				Condition:          make(map[string]interface{}),
				Data:               make(map[string]interface{}),
				UkName:             modelInfo["uk_name"].(string),
				PkName:             modelInfo["pk_name"].(string),
				Order:              modelInfo["select_order"].(string),
				Fields:             "*",
				PrivateFields:      PrivateFieldList,
				LockFields:         LockFieldList,
				ResultExtend: map[string]interface{}{
					"home_url": ps.ByName("home_url"),
				},
				FindSuccess:   FindSuccess,
				SelectSuccess: SelectSuccess,
			},
		}
		
		//!绑定前台的查询条件
		Condition := map[string]interface{}{}
		if condition != "" {
			util.JsonDecode(condition, &Condition)
		}
		ConditionNum := 0
		for k, v := range Condition {
			easyCurd.DbModel.Condition[k] = v
			ConditionNum++
		}
		//!设置默认查询条件
		if !easyCurd.DbModel.SoftDeleteDisable {
			easyCurd.DbModel.Condition["is_delete"] = 0
		}
		//!是否可查询禁用数据
		if !easyCurd.DbModel.SelectWithDisabled {
			easyCurd.DbModel.Condition["status"] = 1
		}
		
		//!绑定前台传递的数据
		Data := map[string]interface{}{}
		if data != "" {
			util.JsonDecode(data, &Data)
		}
		for k, v := range Data {
			easyCurd.DbModel.Data[k] = v
		}
		
		//!检查登录
		if easyCurd.DbModel.CheckLogin {
			//? 获取存到参数里的用户信息
			userInfo := ps.ByName("userInfo")
			if userInfo == "" {
				commHandle.Base{}.ApiResult(w, 100, "请登录账户！", nil)
				return
			}
			
			//? 解析出用户信息
			UserInfo := map[string]interface{}{}
			util.JsonDecode(userInfo, &UserInfo)
			
			//? 校验用户信息
			if util.Interface2String(UserInfo["id"]) == "" || util.Interface2String(UserInfo["id"]) == "0" {
				commHandle.Base{}.ApiResult(w, 100, "登录凭据过期，请重新登录！", nil)
				return
			}
			//? 添加用户ID的查询条件
			if easyCurd.DbModel.UkName == "" {
				easyCurd.DbModel.UkName = "user_id"
			}
			easyCurd.DbModel.Condition[easyCurd.DbModel.UkName] = UserInfo["id"]
		}
		
		//? 执行增删改查方法
		switch util.PostValue(r, "actionName") {
		case "select":
			easyCurd.Select(w, r, ps)
		case "find":
			easyCurd.Find(w, r, ps)
		case "add":
			easyCurd.DbModel.Data["create_time"] = util.TimeNow()
			easyCurd.DbModel.Data["update_time"] = util.TimeNow()
			easyCurd.Create(w, r, ps)
		case "update":
			easyCurd.DbModel.Data["update_time"] = util.TimeNow()
			easyCurd.Update(w, r, ps)
		case "delete":
			if ConditionNum == 0 {
				commHandle.Base{}.ApiResult(w, 120, "缺少删除条件！", nil)
				return
			}
			easyCurd.DbModel.Data["update_time"] = util.TimeNow()
			easyCurd.Delete(w, r, ps)
		default:
			easyCurd.Error(w, r, ps)
		}
	}
}
