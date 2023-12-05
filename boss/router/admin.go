package router

import (
	"go101/model"
	"go101/response"
	"go101/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func pageQueryAdmin(c *gin.Context) {
	limit, offset := util.GetLimitOffset(c)
	admins, err := model.GetAdmins(limit, offset, nil)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, nil)
		return
	}
	total, err := model.CountAdmins(nil)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, nil)
		return
	}
	page := util.PageResult(total, admins)
	response.OK(c, page)
}

func getAdminById(c *gin.Context) {
	id, err := com.StrTo(c.Param("id")).Int64()
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	admin, err := model.GetAdminById(uint(id))
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, admin)
}

type editAdminReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	RoleID   uint   `json:"roleId"`
}

func editAdmin(c *gin.Context) {
	id, err := com.StrTo(c.Param("id")).Int64()
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	var req editAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	a := &model.Admin{
		Model:    model.Model{ID: uint(id)},
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		RoleID:   req.RoleID,
	}
	model.UpdateAdmin(a)

}

type addAdminReq struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	RoleID   uint   `json:"roleId" binding:"required"`
}

func addAdmin(c *gin.Context) {
	var req addAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}

	a := &model.Admin{
		Username: req.Username,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		RoleID:   req.RoleID,
	}
	if err := model.AddAdmin(a); err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, a)
}

func deleteAdmin(c *gin.Context) {
	id, err := com.StrTo(c.Param("id")).Int64()
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := model.DeleteAdmin(uint(id)); err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, nil)
}

type changeStatusReq struct {
	Disable int64 `json:"disable" binding:"min=-1,max=1"`
}

func changeAdminState(c *gin.Context) {
	id, err := com.StrTo(c.Param("id")).Int64()
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	var req changeStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	disabledAt := req.Disable
	if req.Disable > 0 {
		disabledAt = util.GetCurrentTimeMs()
	}
	a := &model.Admin{
		Model:      model.Model{ID: uint(id)},
		DisabledAt: disabledAt,
	}
	model.UpdateAdmin(a)
}
