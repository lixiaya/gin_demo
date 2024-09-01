package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func ZapMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		// 执行时间
		latency := end.Sub(start)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqURI := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		msg := []string{"Code: ", strconv.Itoa(statusCode), "Method:", reqMethod, "latency:", strconv.FormatInt(int64(latency), 10), "Url:", reqURI, "ClientIP:", clientIP}
		logger.Info("请求信息：", strings.Join(msg, " "))
	}
}
