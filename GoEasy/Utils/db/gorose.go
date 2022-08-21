package db

import (
	"github.com/yqstech/gef/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

var err error
var engin *gorose.Engin

func Init() {
	mysqlDsn := config.DbUser + ":" + config.DbPwd + "@tcp(" + config.DbHost + ")/" + config.DbName + "?charset=utf8"
	engin, err = gorose.Open(&gorose.Config{
		Driver:          "mysql",
		Dsn:             mysqlDsn,
		SetMaxOpenConns: config.DbMaxOpenConns,
		SetMaxIdleConns: config.DbMaxIdleConns,
	})
}
func New() gorose.IOrm {
	return engin.NewOrm()
}
