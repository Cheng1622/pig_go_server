package controllers

import (
	"errors"
	"go_server/dao/mysql"
	"go_server/logic"
	"go_server/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	p := new(models.ParamLogin)

	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("LoginHandler with invalid param", zap.Error(err))
		// 因为有的错误 比如json格式不对的错误 是不属于validator错误的 自然无法翻译，所以这里要做类型判断
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
		} else {
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		}
		return
	}

	// 业务处理
	token, err := logic.Login(p)
	if err != nil {
		// 可以在日志中 看出 到底是哪些用户不存在
		zap.L().Error("login failed", zap.String("Email", p.Email), zap.Error(err))
		if errors.Is(err, mysql.WrongPassword) {
			ResponseError(c, CodeInvalidPassword)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	ResponseSuccess(c, token)
}

func RegisterHandler(c *gin.Context) {
	// 获取参数和参数校验
	p := new(models.ParamRegister)
	// 这里只能校验下 是否是标准的json格式 之类的 比较简单
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("RegisterHandler with invalid param", zap.Error(err))
		// 因为有的错误 比如json格式不对的错误 是不属于validator错误的 自然无法翻译，所以这里要做类型判断
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
		} else {
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		}
		return
	}
	// 业务处理
	err := logic.Register(p)
	if err != nil {
		zap.L().Error("register failed", zap.String("email", p.Email), zap.Error(err))
		if errors.Is(err, mysql.UserAleadyExists) {
			ResponseError(c, CodeUserExist)
		} else {
			ResponseError(c, CodeInvalidParam)
		}
		return
	}
	// 返回响应
	ResponseSuccess(c, "register success")
}

func ProfileHandler(c *gin.Context) {
	user := models.User{
		Email: c.GetString("email"),
	}
	data, err := mysql.GetUserByEmail(&user)
	if err != nil {
		zap.L().Error("ProfileHandler", zap.Error(err))
		ResponseError(c, CodeTokenInvalid)
		return
	}
	ResponseSuccess(c, data)

}
