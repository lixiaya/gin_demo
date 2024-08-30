package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// 加密密码
func BcryptHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// 验证密码
func BcryptVerify(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Sprintln("两次密码不一致")
	}
	return true
}
