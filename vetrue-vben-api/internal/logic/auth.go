package logic

import (
	"context"
	"errors"
	"vetrue-vben-api/internal/database"
	"vetrue-vben-api/internal/dto"
	"vetrue-vben-api/models"
	"vetrue-vben-api/pkg/jwt"
	"vetrue-vben-api/pkg/utils"
	"log"
)

type AuthLogic struct {
}

func NewAuthLogic() *AuthLogic {
	return &AuthLogic{}
}

// @Summary SystemAuthLogin
func (self *AuthLogic) SystemAuthLogin(ctx context.Context, params *dto.LoginReq) (resp *dto.LoginResp, err error) {
	user := new(models.SysUsers)
	if err = database.GetDB().Model(&models.SysUsers{}).Where("username = ?", params.Username).First(user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status != models.SysUserStatusNormal {
		return nil, errors.New("用户已被禁用")
	}

	if !utils.ValidatePasswd(params.Password, user.Salt, user.Password) {
		return nil, errors.New("密码错误")
	}

	token, err := jwt.GenerateJwtToken(map[string]interface{}{
		"id":       user.Id,
		"username": user.Username,
		"role_id":  user.RoleId,
	})
	if err != nil {
		log.Printf("生成 token 失败: %v", err)
		return nil, errors.New("生成 token 失败")
	}

	resp = &dto.LoginResp{
		AccessToken: token,
		Id:          user.Id,
		Password:    "",
		RealName:    user.Username,
		Roles:       []string{""},
		Username:    user.Username,
	}

	return
}

// @Summary SystemAuthLogout
func (self *AuthLogic) SystemAuthLogout(ctx context.Context, params *dto.LogoutReq) (err error) {
	// TODO implement

	return
}
