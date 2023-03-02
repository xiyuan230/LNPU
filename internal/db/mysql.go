package db

import (
	"MyLNPU/conf"
	"MyLNPU/internal/log"
	"MyLNPU/internal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func Init() {
	config := conf.GetConfig().MySql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s", config.UserName, config.UserPassword, config.Host, config.Port, config.Database)
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    false, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{})

	if err != nil {
		log.Errorf("数据库连接失败... ERROR: %s", err)
		os.Exit(-1)
	}
	db.AutoMigrate(&model.User{}, &model.Notice{})
	defaultData()
	log.Println("mysql初始化成功...")
}

func GetDB() *gorm.DB {
	return db
}

func defaultData() {
	db.FirstOrCreate(&model.Notice{ID: 1}, &model.Notice{})
}
