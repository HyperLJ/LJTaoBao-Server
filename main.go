package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"itStudioTB/common"
	"os"
)

func main() {
	// 读取数据库配置
	InitConfig()
	// 连接数据库
	db := common.InitDB()
	defer db.Close()

	r := gin.Default()
	r = CollectRoute(r)

	port := viper.GetString("server.port")
	if port != "" {
		// panic 抛出异常
		panic(r.Run(":" + port))
	}

	panic(r.Run())
}

// 配置数据库相关信息
func InitConfig() {
	workDir, _ := os.Getwd()

	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")

	err := viper.ReadInConfig()

	if err != nil {
		panic(err.Error())
	}
}
