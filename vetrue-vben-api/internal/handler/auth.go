package handler

import (
	"vetrue-vben-api/internal/dto"
	"vetrue-vben-api/internal/logic"
	"vetrue-vben-api/pkg/response"

	"github.com/gin-gonic/gin"
)

var authLogic = logic.NewAuthLogic()

// @Summary SystemAuthLogin
// @Description SystemAuthLogin
// @Accept json
// @Produce json
// @Param body dto.LoginReq
// @Success 200 {object} dto.LoginResp
// @Failure 200 {object} response.Response
// @Router /system/auth/login [post]
func SystemAuthLogin(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}

	resp, err := authLogic.SystemAuthLogin(ctx, &req)
	if err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}
	response.Success(ctx, resp)
}

// @Summary SystemAuthLogout
// @Description SystemAuthLogout
// @Accept json
// @Produce json
// @Param body dto.LogoutReq
// @Success 200 string success
// @Failure 200 {object} response.Response
// @Router /system/auth/logout [post]
func SystemAuthLogout(ctx *gin.Context) {
	var req dto.LogoutReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}

	if err := authLogic.SystemAuthLogout(ctx, &req); err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}
