package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	// 密钥
	secretKey = []byte("lucaya")
)

type Claims struct {
	Uid uint `json:"uid"`
	jwt.StandardClaims
}

// 生成JWT
func GenerateJwt(uid uint) (string string, err error) {
	//过期时间 3天
	expireToken := time.Now().Add(time.Hour * 24 * 3).Unix()
	claims := Claims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	}
	//创建jwt令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//生成token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("生成token失败:", err.Error())
		return "", err
	}
	return tokenString, nil
}

// 验证JWT
func ParseJwt(token string) (*jwt.Token, error) {
	claims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
