package pool

import (
	"github.com/gef/config"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

// Redis 连接池
var Redis *redis.Client

func RedisInit() {
	//创建Redis连接池
	if config.RedisOpen == "TRUE" {
		Redis = redis.NewClient(&redis.Options{
			Addr:     config.RedisHost + ":" + config.RedisPort,
			Password: config.RedisPwd,
			DB:       config.String2Int(config.RedisDb),
		})
	}
}
