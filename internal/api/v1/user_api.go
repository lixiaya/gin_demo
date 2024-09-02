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
	"strconv"
	"time"
)

type UserApi struct {
}

// LoginSuccessData 登录成功返回数据
type LoginSuccessData struct {
	Uid   uint   `json:"uid"`
	Token string `json:"token"`
}

type capCode struct {
	Id   string `json:"id"`
	B64  string `json:"b64"`
	Code string `json:"code"`
}

// Login 登录
func (u *UserApi) Login(ctx *gin.Context) {
	var l model.User
	err := ctx.ShouldBind(&l)
	if err != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//查询数据是否存在email，存在则验证密码，不存在则返回失败
	loginUser := dto.UserLogin{
		Uid:      l.UID,
		Email:    l.Email,
		Password: l.Password,
	}
	if global.DB.Where("email = ?", loginUser.Email).First(&l) == nil {
		util.ResponseErr(ctx, http.StatusBadRequest, "email已存在")
	}
	fmt.Println(l)
	//验证密码
	match := util.BcryptVerify(l.Password, loginUser.Password)
	if !match {
		util.ResponseErr(ctx, http.StatusBadRequest, "密码错误")
		return
	}
	//生成jwt——token ,保存到redis，过期时间为三天
	generateJwt, err := util.GenerateJwt(l.UID)
	if err != nil {
		util.ResponseErr(ctx, http.StatusInternalServerError, "生成token失败")
		return
	}
	data := LoginSuccessData{
		Uid:   l.UID,
		Token: generateJwt,
	}
	global.Rdb.Set(context.Background(), strconv.Itoa(int(l.UID)), generateJwt, 3*24*time.Hour)
	util.ResponseOk(ctx, http.StatusOK, "login success", data)
}

// Register 注册
func (u *UserApi) Register(ctx *gin.Context) {
	var user dto.UserRegister
	if err := ctx.ShouldBind(&user); err != nil {
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

// GenerateCaptcha 发送验证码请求
func (u *UserApi) GenerateCaptcha(ctx *gin.Context) {
	var user dto.UserRegister
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println("失败")
		return
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

// GetUserAllInfo 获取所有用户信息
func (u *UserApi) GetUserAllInfo(ctx *gin.Context) {
	//从上下文获取用户id
	var l []model.User
	uid := ctx.MustGet("uid")
	fmt.Printf("uid:%v ", uid)
	//查询所有用户信息
	global.DB.Find(&l)
	newUsers := make([]model.User, 0)
	for _, user := range l {
		newUsers = append(newUsers, user)
	}
	//记录那位用户访问
	global.Logger.Info(fmt.Sprintf("UID: %v 查询了所有用户信息", uid))
	util.ResponseOk(ctx, http.StatusOK, "获取所有用户信息成功", newUsers)
}

// GetUserInfo 查询单个用户信息
func (u *UserApi) GetUserInfo(ctx *gin.Context) {
	//通过email 查询用户信息
	var l model.User
	email := ctx.Param("email")
	fmt.Println("email:", email)
	err := global.DB.Where("email = ?", email).First(&l).Error
	if err != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, "用户不存在")
		return
	}
	util.ResponseOk(ctx, http.StatusOK, "查询用户信息成功", l)
}

// UpdateUserInfo 修改用户信息
// 目前字段较少 只能修改昵称和性别
func (u *UserApi) UpdateUserInfo(ctx *gin.Context) {
	var l model.User
	err := ctx.ShouldBindJSON(&l)
	fmt.Printf("l:", l)
	if err != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, err.Error)
		return
	}
	type updataUser struct {
		Nickname  string
		Gender    string
		UpdatedAt time.Time
	}
	newl := updataUser{
		Nickname:  l.Nickname,
		Gender:    l.Gender,
		UpdatedAt: time.Now(),
	}

	resultErr := global.DB.Model(&l).Where("email=?", l.Email).Updates(newl).Error
	fmt.Println(l.Email)
	fmt.Println(resultErr)
	if resultErr != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, "更新用户信息失败")
		return
	}
	util.ResponseOk(ctx, http.StatusOK, "更新用户信息成功", newl)
}

// DeleteUserInfo 删除用户信息
func (u *UserApi) DeleteUserInfo(ctx *gin.Context) {
	var email string
	email = ctx.Param("email")
	resultErr := global.DB.Where("email=?", email).Delete(&model.User{}).Error
	if resultErr != nil {
		util.ResponseErr(ctx, http.StatusBadRequest, "删除用户信息失败")
		return
	}
	util.ResponseOk(ctx, http.StatusOK, "删除用户信息成功", email)

}
