package util

import (
	"github.com/yqstech/gef/config"
	"strings"
)

// Sql 数据库的sql语句过滤
func Sql(sql string) string {
	if config.DbType == "sqlite" || config.DbType == "sqlite3" {
		//数据库的方法兼容处理
		sql = strings.Replace(sql, "CURRENT_DATE()", "date()", -1)
	} else if config.DbType == "mysql" {
		//数据库方法兼容处理
		sql = strings.Replace(sql, "date()", "CURRENT_DATE()", -1)
	}
	return sql
}
