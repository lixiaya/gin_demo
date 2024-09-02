package middleware

import (
	"gin_demo/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带token，无权限访问"})
			ctx.Abort()
			return
		}
		//去除Bearer和两边空格
		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
		//验证token签名
		parseToken, err := util.ParseJwt(tokenString)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"解析token error": err.Error()})
			}
			ctx.Abort()
			return
		}
		if claims, ok := parseToken.Claims.(*util.Claims); ok && parseToken.Valid {
			ctx.Set("uid", claims.Uid)
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token验证失败"})
		}

	}
}
