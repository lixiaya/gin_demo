package config

import "github.com/gin-contrib/cors"

func ConfigCors() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"} //允许的源
	corsConfig.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS",
	} //允许的http请求方式
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // 允许的头部
	corsConfig.AllowCredentials = true                                            // 允许发送Cookie
	return corsConfig
}
