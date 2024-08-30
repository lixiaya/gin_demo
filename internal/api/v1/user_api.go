package v1

import (
	"context"
	"fmt"
	"gin_demo/global"
	"gin_demo/internal/model"
	"gin_demo/internal/service/dto"
	"gin_demo/util"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserApi struct {
}

type capCode struct {
	Id   string `json:"id"`
	B64  string `json:"b64"`
	Code string `json:"code"`
}

func (u *UserApi) Login(ctx *gin.Context) {
	var l model.User
	err := ctx.ShouldBind(&l)
	//fmt.Println(l)
	if err != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//查询数据是否存在email，存在则验证密码，不存在则返回失败
	loginUser := dto.UserLogin{
		Email:    l.Email,
		Password: l.Password,
	}
	fmt.Println(loginUser)
	if global.DB.Where("email = ?", loginUser.Email).First(&l) == nil {
		util.ResponseErr(ctx, http.StatusBadRequest, "email已存在")
	}
	//验证密码
	match := util.BcryptVerify(l.Password, loginUser.Password)

	if !match {
		util.ResponseErr(ctx, http.StatusBadRequest, "密码错误")
		return
	}

	util.ResponseOk(ctx, http.StatusOK, "login success", "")
}

func (u *UserApi) Register(ctx *gin.Context) {
	var user dto.UserRegister
	err := ctx.ShouldBind(&user)
	if err != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, err.Error())
	}
	//检查email是否重复
	if global.DB.Where("email = ?", user.Email).First(&user) == nil {
		util.ResponseErr(ctx, http.StatusBadRequest, "username 重复")
		return
	}
	//检查验证码是否正确
	key, err := global.Rdb.Get(context.Background(), user.Email).Result()
	if err != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, "获取验证码失败")
		return
	}
	if key != user.Captcha {
		util.ResponseErr(ctx, http.StatusBadRequest, "验证码错误")
		return
	}

	//加密后的密码
	bcryptPassword := util.BcryptHash(user.Password)
	newUser := &model.User{
		Email:     user.Email,
		Password:  bcryptPassword,
		Gender:    user.Gender,
		Nickname:  user.Nickname,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
	//注册信息插入数据库
	if global.DB.Create(newUser).Error != nil {
		util.ResponseErr(ctx, http.StatusInternalServerError, "注册失败")
		return
	}

	util.ResponseOk(ctx, http.StatusOK, "注册success", newUser)
}

/*
*
发送验证码请求
*/
func (u *UserApi) GenerateCaptcha(ctx *gin.Context) {
	var user dto.UserRegister
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println("失败")
	}

	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	store := base64Captcha.DefaultMemStore
	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, code, err := captcha.Generate()
	newCode := &capCode{
		id,
		b64s,
		code,
	}
	if err != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//发送验证码id和code到redis 过期时间五分钟
	//将邮箱地址（key）和生成的验证码（value）存入redis
	global.Rdb.Set(context.Background(), user.Email, newCode.Code, 5*time.Minute)
	//发送验证码
	err = util.SendEmail(user.Email, newCode.Code)
	if err != nil {
		fmt.Println("send email fail")
	}
	util.ResponseOk(ctx, http.StatusOK, "生成验证码成功", newCode)

}
