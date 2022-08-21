/*
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Date: 2020-11-17 23:03:31
 * @LastEditTime: 2021-10-20 12:36:10
 * @Description:管理员
 */
package Models

import (
	"github.com/yqstech/gef/GoEasy/Utils/db"
	"github.com/yqstech/gef/GoEasy/Utils/util"

	"github.com/wonderivan/logger"
)

type Admin struct {
}

func (that Admin) GetAccountInfoByToken(token string) map[string]interface{} {
	conn := db.New()
	tokenInfo, err := conn.Table("tb_admin_token").
		Where(map[string]interface{}{"token": token, "status": 1, "is_delete": 0}).
		Where("create_time", ">", util.TimeNowFormat("2006-01-02 15:04:05", 0, 0, -15)).
		First()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	if tokenInfo == nil {
		return nil
	}
	if tokenInfo["account_id"].(int64) == 0 {
		return nil
	}
	//获取供货商信息
	accountId := tokenInfo["account_id"].(int64)
	AdminInfo, err := db.New().Table("tb_admin").Where(map[string]interface{}{
		"id":        accountId,
		"is_delete": 0,
		"status":    1,
	}).First()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	if AdminInfo == nil {
		return nil
	}
	return map[string]interface{}{
		"account_id": AdminInfo["id"],
		"group_id":   AdminInfo["group_id"],
		"name":       AdminInfo["name"],
		"account":    AdminInfo["account"],
	}
}

//
// CheckAuth
//  @Description:
//  @receiver that
//  @param rule_name
//  @param accountID
//  @return bool
//
func (that Admin) CheckAuth(rule_name string, accountID int) bool {
	//查询用户信息
	userinfo, err := db.New().Table("tb_admin").
		Where("id", accountID).
		Where("is_delete", 0).
		Where("status", 1).First()
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	if userinfo == nil {
		return false
	}
	if userinfo["id"] == int64(1) {
		return true
	}
	//用户分配了角色，查询角色权限IDs
	ruleIDs := []int64{0}
	if userinfo["group_id"] != int64(0) {
		//查询账户所在角色
		groupInfo, err := db.New().Table("tb_admin_group").
			Where("id", userinfo["group_id"]).
			Where("is_delete", 0).
			Where("status", 1).
			First()
		if err != nil {
			logger.Error(err.Error())
			return false
		}
		if groupInfo != nil {
			util.JsonDecode(groupInfo["rules"].(string), &ruleIDs)
		}
	}
	conn := db.New().Table("tb_admin_rules")
	rule, err := conn.
		Where("is_delete", "=", 0).
		Where("status", "=", 1).
		Where("route", "=", rule_name).
		Where(func() {
			conn.Where("is_compel", 1).OrWhere("id", "in", ruleIDs)
		}).First()
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	if rule == nil {
		return false
	}
	return true
}
