package config

import (
	"fmt"
	"gin_demo/global"

	"github.com/spf13/viper"
)

func ViperConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("配置文件读取时出错：%v", err))
		//panic(fmt.Errorf("fatal error config file: %w", err))
	}
	global.Logger.Info("初始化viper Success")
	fmt.Println(viper.GetString("server.port"))
}
