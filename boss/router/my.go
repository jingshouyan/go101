package router

import (
	"go101/g"
	"go101/model"
	"go101/response"
	"go101/util"
	"net/http"

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
		response.BizError(c, response.ErrorNotExist)
		return
	}
	if a.DisabledAt > 0 {
		response.BizError(c, response.AccountDisabled)
		return
	}
	check := util.CheckPasswordHash(req.Password, a.PwdHash)
	if !check {
		response.BizError(c, response.PasswordWrong)
		return
	}
	util.SaveAdminIDToSession(c, a.ID)
	response.OK(c, a)

}

func logout(c *gin.Context) {
	util.ClearSession(c)
	response.OK(c)
}

type changePwdReq struct {
	OldPwd string `json:"oldPwd" binding:"required"`
	NewPwd string `json:"newPwd" binding:"required"`
}

func changePwd(c *gin.Context) {
	var req changePwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	aid := c.GetUint(g.AdminIdKey)
	admin, err := model.GetAdminById(aid)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	check := util.CheckPasswordHash(req.OldPwd, admin.PwdHash)
	if !check {
		response.BizError(c, response.PasswordWrong)
		return
	}
	hash, err := util.HashPassword(req.NewPwd)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	a := &model.Admin{
		Model:   model.Model{ID: admin.ID},
		PwdHash: hash,
	}
	model.UpdateAdmin(a)
	response.OK(c)

}

func getProfile(c *gin.Context) {
	aid := c.GetUint(g.AdminIdKey)
	admin, err := model.GetAdminById(aid)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, admin)
}

type updateProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func updateProfile(c *gin.Context) {
	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	aid := c.GetUint(g.AdminIdKey)
	a := &model.Admin{
		Model:    model.Model{ID: aid},
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}
	model.UpdateAdmin(a)
	response.OK(c)
}

func ping(c *gin.Context) {
	response.OK(c, "pong")
}
