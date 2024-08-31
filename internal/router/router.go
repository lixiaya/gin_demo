package router

import (
	"fmt"
	"gin_demo/config"
	"gin_demo/global"
	api "gin_demo/internal/api/v1"
	"gin_demo/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouter() {
	r := gin.Default()

	//跨域
	corsconfig := config.ConfigCors()
	r.Use(cors.New(corsconfig))
	//r.Use()
	r.Use(middleware.ZapMiddleware(global.Logger))
	userApi := api.UserApi{}
	v1 := r.Group("/api/v1")
	{

		v1.POST("/generate-captcha", userApi.GenerateCaptcha)
		v1.GET("/login", userApi.Login)
		v1.POST("/register", userApi.Register)
	}

	err := r.Run(fmt.Sprintf(":%v", viper.GetString("server.port")))
	if err != nil {
		global.Logger.Error(fmt.Sprintln("gin  server  start error:", err))
		panic(err.Error())
	}
	global.Logger.Info(fmt.Sprintln("server run success on ", viper.GetString("server.port")))
}
