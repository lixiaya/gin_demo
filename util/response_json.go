package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JsonRespOK struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}
type JsonRespErr struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

// 封装请求成功返回的数据
func ResponseOk(ctx *gin.Context, code int, msg interface{}, data interface{}) {
	json := &JsonRespOK{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	ctx.JSON(http.StatusOK, json)
}

// 封装请求失败返回的数据
func ResponseErr(ctx *gin.Context, code int, msg interface{}) {
	json := &JsonRespErr{
		Code: code,
		Msg:  msg,
	}
	ctx.JSON(http.StatusBadRequest, json)
}
