/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Manager
 * @Version: 1.0.0
 * @Date: 2022/8/16 13:47
 */

package database

import (
	"github.com/gef/GoEasy/Utils/gdb"
	"github.com/wonderivan/logger"
)

// AutoMigrate 自动维护表结构
func AutoMigrate() {
	err := gdb.New().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(Tables...)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	DataInit()
}

// DataInit 检查数据表并初始化数据
func DataInit() {
	//设置默认数据
	DefaultData(defaultData)
	//自动维护权限表
	AutoAdminRules(adminRules)
}
