/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Nodes
 * @Version: 1.0.0
 * @Date: 2021/12/7 8:34 下午
 */

package EasyApp

import (
	"github.com/gohouse/gorose/v2"
	util2 "github.com/yqstech/gef/Utils/util"
)

// ================  默认节点 ====================================

//! getHandle内部节点

// PageInit 初始化节点（用于设置action列表等信息）
func (nf Page) PageInit(pageData *PageData) {

}

//! realHandle内部节点

// NodeBegin 开始节点
func (nf Page) NodeBegin(pageData *PageData) (error, int) {
	//logger.Alert("默认NodeBegin")
	return nil, 0
}

// NodeList 列表开始节点
func (nf Page) NodeList(pageData *PageData) (error, int) {
	//logger.Alert("NodeList")
	return nil, 0
}

func (nf Page) NodeStatusSuccess(pageData *PageData, id, status int64) (error, int) {
	//logger.Alert("NodeStatusSuccess")
	return nil, 0
}

// NodeListCondition 列表查询条件修改节点
func (nf Page) NodeListCondition(pageData *PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	//logger.Alert("NodeListCondition")
	return condition, nil, 0
}

// NodeListCondition 列表查询条件修改节点
func (nf Page) NodeAutoCondition(pageData *PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	//软删除
	if pageData.deleteField != "" {
		condition = append(condition, []interface{}{pageData.deleteField, "=", 0})
	}

	return condition, nil, 0
}

//
func (nf Page) NodeListOrm(pageData *PageData, orm *gorose.IOrm) (error, int) {
	return nil, 0
}

// NodeListData 列表查询出数据节点
func (nf Page) NodeListData(pageData *PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	//logger.Alert("NodeListData")
	return data, nil, 0
}

func (nf Page) NodeCheckAuth(pageData *PageData, btnRule string, accountID int) (bool, error) {
	return true, nil
}

func (nf Page) NodeForm(pageData *PageData, id int64) (error, int) {
	return nil, 0
}

func (nf Page) NodeFormData(pageData *PageData, data gorose.Data, id int64) (gorose.Data, error, int) {
	return data, nil, 0
}

func (nf Page) NodeSaveData(pageData *PageData, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {

	return postData, nil, 0
}
func (nf Page) NodeDeleteBefore(pageData *PageData, id int64) (error, int) {
	return nil, 0
}

// NodeAutoData 增改数据自动完成
func (nf Page) NodeAutoData(pageData *PageData, postData map[string]interface{}, action string) (map[string]interface{}, error, int) {
	if action == "add" {
		//非数据库自增，设置int64 ID
		if !pageData.isAutoID {
			if id, ok := postData[pageData.tbPK]; ok {
				if id == "" || util2.Interface2String(id) == "0" {
					postData[pageData.tbPK] = util2.GetOnlyID()
				}
			}
		} else {
			//数据库自增，清理无用ID
			if _, ok := postData[pageData.tbPK]; ok {
				delete(postData, pageData.tbPK)
			}
		}
	}
	if action == "edit" {
		//自动删除ID
		if _, ok := postData[pageData.tbPK]; ok {
			delete(postData, pageData.tbPK)
		}
	}
	//数组类型的自动转json（checkbox组件）
	for key, value := range postData {
		postData[key] = util2.Array2String(value)
	}

	//其他字段自动完成
	for _, fieldKey := range pageData.insertAutoFields {
		if fieldKey == "create_time" && action == "add" {
			postData["create_time"] = util2.TimeNow()
		}
		if fieldKey == "update_time" {
			postData["update_time"] = util2.TimeNow()
		}
	}
	for _, fieldKey := range pageData.updateAutoFields {
		if fieldKey == "update_time" {
			postData["update_time"] = util2.TimeNow()
		}
	}

	return postData, nil, 0
}

func (nf Page) NodeSaveSuccess(pageData *PageData, postData map[string]interface{}, id int64) (bool, error, int) {
	return true, nil, 0
}
func (nf Page) NodeUpdateSuccess(pageData *PageData, postData map[string]interface{}, id int64) (bool, error, int) {
	return true, nil, 0
}
func (nf Page) NodeAddSuccess(pageData *PageData, postData map[string]interface{}, id int64) (bool, error, int) {
	return true, nil, 0
}

func (nf Page) NodeDeleteSuccess(pageData *PageData, id int64, deleteField string) (error, int) {
	return nil, 0
}
