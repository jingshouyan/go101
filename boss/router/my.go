package router

import (
	"go101/model"
	"go101/response"
	"go101/util"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Type     int    `json:"type" binding:"required"`
}

func login(c *gin.Context) {
	var req loginReq
	err := c.BindJSON(&req)
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	a, err := model.GetAdminByAccount(req.Type, req.Account)
	if err != nil {
		response.BizError(c, response.ErrorNotExist, nil)
		return
	}
	check := util.CheckPasswordHash(req.Password, a.PwdHash)
	if !check {
		response.BizError(c, response.PasswordWrong, nil)
		return
	}
	session := sessions.Default(c)
	session.Set("uid", a.ID)
	session.Save()
	response.OK(c, a)

}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	response.OK(c, nil)
}

func ping(c *gin.Context) {
	response.OK(c, "pong")
}
