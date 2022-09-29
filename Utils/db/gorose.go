package db

import (
	"github.com/yqstech/gef/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
)

var err error
var engin *gorose.Engin

func Init() {
	if config.DbType == "mysql" {
		mysqlDsn := config.DbUser + ":" + config.DbPwd + "@tcp(" + config.DbHost + ")/" + config.DbName + "?charset=utf8"
		engin, err = gorose.Open(&gorose.Config{
			Driver:          config.DbType,
			Dsn:             mysqlDsn,
			SetMaxOpenConns: config.DbMaxOpenConns,
			SetMaxIdleConns: config.DbMaxIdleConns,
		})
	} else if config.DbType == "sqlite3" {
		engin, err = gorose.Open(&gorose.Config{
			Driver:          config.DbType,
			Dsn:             config.DbFile,
			SetMaxOpenConns: config.DbMaxOpenConns,
			SetMaxIdleConns: config.DbMaxIdleConns,
		})
	} else {
		panic("数据库类型不支持！")
	}

}
func New() gorose.IOrm {
	return engin.NewOrm()
}
