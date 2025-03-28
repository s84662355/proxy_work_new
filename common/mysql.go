package common

import (
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"mproxy/config"
)

var GetMysqlDB = sync.OnceValue[*gorm.DB](func() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config.ConfData.Mysql.Source,
	}), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(config.ConfData.Mysql.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(config.ConfData.Mysql.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConfData.Mysql.SetConnMaxLifetime) * time.Second)

	return db
})
