package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"itStudioTB/model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := "000926" // viper会把0吞掉,上传至服务器时 尽量将密码改为非数字开头
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
	)

	fmt.Println(args)

	// 连接打开数据库
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}

	// 赋值给全局变量
	DB = db

	// 自动生成表
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Good{})
	db.AutoMigrate(&model.Order{})

	return db
}

func GetDB() *gorm.DB {
	return DB
}
