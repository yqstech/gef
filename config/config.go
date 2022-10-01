package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// 系统
var (
	GOOS     = "linux"
	WorkPath = ""
	AppPath  = ""
	Debug    = "" //是否开启调试模式
)

// 数据库
var (
	DbType         = "mysql"       //mysql sqlite3
	DbFile         = "./db.sqlite" //针对于sqlite
	DbHost         = "localhost:3306"
	DbName         = "db-name"
	DbUser         = "db-user"
	DbPwd          = "db-password"
	DbMaxOpenConns = 100
	DbMaxIdleConns = 50
)

// Redis 缓存
var (
	RedisOpen = "FALSE"
	RedisHost = "localhost"
	RedisPort = "6379"
	RedisPwd  = ""
	RedisDb   = "0"
)

// web服务
var (
	AdminPort = "9001"     //后台端口号
	AdminPath = "/admin"   //后台路径前缀
	FrontPort = "9002"     //前台服务端口号
	ApiPath   = "/app_api" //应用接口路径前缀
)

// 上传参数
var (
	UploadTableName = "tb_attachment"
	UploadPath      = "./data/uploads"
	UploadUrl       = "/uploads"
)

// Init 在环境变量中读取配置信息
func Init() error {
	var err error
	//! 获取运行的系统类型
	GOOS = runtime.GOOS

	//! 获取工作目录
	WorkPath, err = os.Getwd()
	if err != nil {
		return err
	}
	//! 获取应用目录（编译后的程序目录）
	AppPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	//# 读取应用配置文件
	if _, err = os.Stat(WorkPath + "/configs/.env"); err == nil {
		err = godotenv.Load(WorkPath + "/configs/.env")
		if err != nil {
			return err
		}
	} else {
		if _, err = os.Stat(AppPath + "/configs/.env"); err == nil {
			err = godotenv.Load(WorkPath + "/configs/.env")
			if err != nil {
				return err
			}
		} else {
			return errors.New("获取系统配置信息失败！")
		}
	}
	//!开启调试模式
	Debug = os.Getenv("Debug")
	//# 数据库
	//! 新增支持自定义数据库类型
	if os.Getenv("DbType") != "" {
		DbType = os.Getenv("DbType")
	}
	if os.Getenv("DbFile") != "" {
		DbFile = os.Getenv("DbFile")
	}
	DbHost = os.Getenv("DbHost")
	DbName = os.Getenv("DbName")
	DbUser = os.Getenv("DbUser")
	DbPwd = os.Getenv("DbPwd")
	DbMaxOpenConns = String2Int(os.Getenv("DbMaxOpenConns"))
	DbMaxIdleConns = String2Int(os.Getenv("DbMaxIdleConns"))
	//# Redis信息 缓存
	RedisOpen = os.Getenv("RedisOpen")
	RedisHost = os.Getenv("RedisHost")
	RedisPort = os.Getenv("RedisPort")
	RedisPwd = os.Getenv("RedisPwd")
	RedisDb = os.Getenv("RedisDb")
	//# web服务
	AdminPort = os.Getenv("AdminPort")
	AdminPath = os.Getenv("AdminPath")
	FrontPort = os.Getenv("FrontPort")
	ApiPath = os.Getenv("ApiPath")
	//# 上传配置
	UploadTableName = os.Getenv("UploadTableName")
	UploadPath = os.Getenv("UploadPath")
	UploadUrl = os.Getenv("UploadUrl")
	return nil
}
func String2Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
