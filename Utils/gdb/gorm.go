/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: gorm
 * @Version: 1.0.0
 * @Date: 2022/8/16 12:03
 */

package gdb

import (
	"github.com/yqstech/gef/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var err error
var db *gorm.DB

func Init() {
	if config.DbType == "mysql" {
		dsn := config.DbUser + ":" + config.DbPwd + "@tcp(" + config.DbHost + ")/" + config.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("数据库链接失败！")
		}
		sqlDB, _ := db.DB()
		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(config.DbMaxIdleConns)
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(config.DbMaxOpenConns)
	} else if config.DbType == "sqlite" || config.DbType == "sqlite3" {
		db, err = gorm.Open(sqlite.Open(config.DbFile), &gorm.Config{})
		if err != nil {
			panic("数据库链接失败！")
		}
	} else {
		panic("数据库类型不支持！")
	}
}
func New() *gorm.DB {
	return db
}
