package common

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mproxy/config"
)

var GetMysqlDB = sync.OnceValue[*gorm.DB](func() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config.GetConf().Mysql.Source,
	}), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(fmt.Errorf("初始化MySQL失败error:%+v", err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("初始化MySQL失败error:%+v", err))
	}
	sqlDB.SetMaxIdleConns(config.GetConf().Mysql.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(config.GetConf().Mysql.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetConf().Mysql.SetConnMaxLifetime) * time.Second)
	return db
})
