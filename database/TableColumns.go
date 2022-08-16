/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: TableColumns
 * @Version: 1.0.0
 * @Date: 2022/8/16 13:46
 */

package database

import "time"

type ID struct {
	ID int `gorm:"column:id;type:int(11);AUTO_INCREMENT;primary_key" json:"id"`
}
type CUSD struct {
	CreateTime time.Time `gorm:"column:create_time;type:DATETIME;default:CURRENT_TIMESTAMP;NOT NULL;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:DATETIME;default:CURRENT_TIMESTAMP;NOT NULL;comment:更新时间" json:"update_time"`
	Status     int       `gorm:"column:status;type:tinyint(1);default:1;NOT NULL;comment:状态" json:"status"`
	IsDelete   int       `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL;comment:是否删除" json:"is_delete"`
}
type CD struct {
	CreateTime time.Time `gorm:"column:create_time;type:DATETIME;default:CURRENT_TIMESTAMP;NOT NULL;comment:创建时间" json:"create_time"`
	IsDelete   int       `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL;comment:是否删除" json:"is_delete"`
}

