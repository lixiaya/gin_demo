package main

import (
	"fmt"
	"gin_demo/config"
	"gin_demo/global"
	"gin_demo/internal/router"
)

func main() {

	//初始化zap
	global.Logger = config.InitLogger()
	global.Logger.Info("----------start-----------")
	//err := util.SendEmail("2439438173@qq.com", "yzm")
	//err := util.SendEmail("trhxnlove@yeah.net", "yzm")
	//if err != nil {
	//	fmt.Println("send email err:", err)
	//}
	//viper读取配置文件
	config.ViperConfig()
	//redis初始化
	global.Rdb, _ = config.InitRedis()
	fmt.Println(global.Rdb)
	//global.Rdb.Set(context.Background(), "1", "21", 5*time.Minute)
	//数据库初始化连接
	global.DB = config.InitMysql()
	//jwt, err := util.GenerateJwt(1)
	//if err != nil {
	//	return
	//}
	//fmt.Println(jwt)
	//gin路由初始化
	router.InitRouter()
}
