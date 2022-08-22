/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Nodes
 * @Version: 1.0.0
 * @Date: 2021/12/7 8:34 下午
 */

package builder

import (
	"github.com/gohouse/gorose/v2"
	"github.com/yqstech/gef/Utils/util"
)

// ================  默认节点 ====================================

// NodeInit 页面初始化节点
func (that NodePage) NodeInit(pageBuilder *PageBuilder) {}

// NodeBegin 开始节点
func (that NodePage) NodeBegin(pageBuilder *PageBuilder) (error, int) {
	return nil, 0
}

// NodeList 列表开始节点
func (that NodePage) NodeList(pageBuilder *PageBuilder) (error, int) {
	//logger.Alert("NodeList")
	return nil, 0
}

// NodeStatusSuccess 修改状态成功节点
func (that NodePage) NodeStatusSuccess(pageBuilder *PageBuilder, id, status int64) (error, int) {
	//logger.Alert("NodeStatusSuccess")
	return nil, 0
}

// NodeListCondition 列表查询条件修改节点
func (that NodePage) NodeListCondition(pageBuilder *PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {
	//logger.Alert("NodeListCondition")
	return condition, nil, 0
}

// NodeAutoCondition 自动设置查询条件节点
func (that NodePage) NodeAutoCondition(pageBuilder *PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {
	//软删除
	if pageBuilder.deleteField != "" {
		condition = append(condition, []interface{}{pageBuilder.deleteField, "=", 0})
	}
	return condition, nil, 0
}

// NodeListOrm 设置Orm对象节点
func (that NodePage) NodeListOrm(pageBuilder *PageBuilder, orm *gorose.IOrm) (error, int) {
	return nil, 0
}

// NodeListData 列表查询数据节点
func (that NodePage) NodeListData(pageBuilder *PageBuilder, data []gorose.Data) ([]gorose.Data, error, int) {
	//logger.Alert("NodeListData")
	return data, nil, 0
}

// NodeCheckAuth 权限校验接口
func (that NodePage) NodeCheckAuth(pageBuilder *PageBuilder, btnRule string, accountID int) (bool, error) {
	return true, nil
}

// NodeForm 表单构建节点
func (that NodePage) NodeForm(pageBuilder *PageBuilder, id int64) (error, int) {
	return nil, 0
}

// NodeFormData 表单查询出数据节点
func (that NodePage) NodeFormData(pageBuilder *PageBuilder, data gorose.Data, id int64) (gorose.Data, error, int) {
	return data, nil, 0
}

// NodeSaveData 表单保存数据前节点
func (that NodePage) NodeSaveData(pageBuilder *PageBuilder, oldData gorose.Data, postData map[string]interface{}) (map[string]interface{}, error, int) {

	return postData, nil, 0
}

// NodeAutoData 编辑数据自动完成节点
func (that NodePage) NodeAutoData(pageBuilder *PageBuilder, postData map[string]interface{}, action string) (map[string]interface{}, error, int) {
	if action == "add" {
		//非数据库自增，设置int64 ID
		if !pageBuilder.isAutoID {
			if id, ok := postData[pageBuilder.tbPK]; ok {
				if id == "" || util.Interface2String(id) == "0" {
					postData[pageBuilder.tbPK] = util.GetOnlyID()
				}
			}
		} else {
			//数据库自增，清理无用ID
			if _, ok := postData[pageBuilder.tbPK]; ok {
				delete(postData, pageBuilder.tbPK)
			}
		}
	}
	if action == "edit" {
		//自动删除ID
		if _, ok := postData[pageBuilder.tbPK]; ok {
			delete(postData, pageBuilder.tbPK)
		}
	}
	//数组类型的自动转json（checkbox组件）
	for key, value := range postData {
		postData[key] = util.Array2String(value)
	}

	//其他字段自动完成
	for _, fieldKey := range pageBuilder.insertAutoFields {
		if fieldKey == "create_time" && action == "add" {
			postData["create_time"] = util.TimeNow()
		}
		if fieldKey == "update_time" {
			postData["update_time"] = util.TimeNow()
		}
	}
	for _, fieldKey := range pageBuilder.updateAutoFields {
		if fieldKey == "update_time" {
			postData["update_time"] = util.TimeNow()
		}
	}

	return postData, nil, 0
}

// NodeSaveSuccess 保存成功节点
func (that NodePage) NodeSaveSuccess(pageBuilder *PageBuilder, postData map[string]interface{}, id int64) (bool, error, int) {
	return true, nil, 0
}

// NodeUpdateSuccess 修改成功节点
func (that NodePage) NodeUpdateSuccess(pageBuilder *PageBuilder, postData map[string]interface{}, id int64) (bool, error, int) {
	return true, nil, 0
}

// NodeAddSuccess 添加成功节点
func (that NodePage) NodeAddSuccess(pageBuilder *PageBuilder, postData map[string]interface{}, id int64) (bool, error, int) {
	return true, nil, 0
}

// NodeDeleteBefore 删除前节点
func (that NodePage) NodeDeleteBefore(pageBuilder *PageBuilder, id int64) (error, int) {
	return nil, 0
}

// NodeDeleteSuccess 删除成功节点
func (that NodePage) NodeDeleteSuccess(pageBuilder *PageBuilder, id int64, deleteField string) (error, int) {
	return nil, 0
}
