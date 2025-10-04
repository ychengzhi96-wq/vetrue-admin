package logic

import (
	"context"
	"errors"
	"fmt"
	"vetrue-vben-api/internal/database"
	"vetrue-vben-api/internal/dto"
	"vetrue-vben-api/models"
	"vetrue-vben-api/pkg/middleware"
	"vetrue-vben-api/pkg/utils"
	"log"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var apiInfoMap = make(map[string]*models.SysApis)
var recordUseApiInfoMap = make(map[string]*models.SysApis)

type SystemLogic struct {
}

func NewSystemLogic() *SystemLogic {
	return &SystemLogic{}
}

// @Summary ApiAdd
func (self *SystemLogic) ApiAdd(ctx context.Context, params *dto.UpsertApiReq) (err error) {
	var has int64
	database.GetDB().Model(&models.SysApis{}).Where("parent_id != ?", 0).Where("path = ?", params.Path).Count(&has)
	if has > 0 {
		return fmt.Errorf("API已存在！")
	}

	err = database.GetDB().Create(&models.SysApis{
		ParentId:    params.ParentId,
		Description: params.Description,
		Method:      params.Method,
		Path:        params.Path,
	}).Error

	return
}

// @Summary ApiList
func (self *SystemLogic) ApiList(ctx context.Context, params *dto.ApiListReq) (resp *dto.ApiListResp, err error) {
	resp = &dto.ApiListResp{}
	query := database.GetDB().Model(&models.SysApis{})
	if params.Description != "" {
		query.Where("description like ?", params.Description+"%")
	}
	if params.Path != "" {
		query.Where("path = ?", params.Path)
	}
	if params.OnlyParent {
		query.Where("parent_id = ?", 0)
	}

	var list []*models.SysApis
	if err = query.Find(&list).Error; err != nil {
		return
	}

	items := make([]*dto.ApiInfoResp, 0)
	for _, info := range list {
		item := new(dto.ApiInfoResp)
		_ = copier.Copy(item, info)
		item.CreatedAt = info.CreatedAt.UnixMilli()
		items = append(items, item)
	}
	if params.OnlyParent {
		items = append(items, &dto.ApiInfoResp{
			Id:          0,
			CreatedAt:   0,
			ParentId:    0,
			Description: "根API",
			Method:      "",
			Path:        "",
		})
	}
	resp.Items = items
	resp.Total = int64(len(items))

	return
}

// @Summary ApiUpdate
func (self *SystemLogic) ApiUpdate(ctx context.Context, id int64, params *dto.UpsertApiReq) (err error) {
	var has int64
	database.GetDB().Model(&models.SysApis{}).Where("parent_id != ?", 0).Where("path = ?", params.Path).Count(&has)
	if has > 0 {
		return fmt.Errorf("API已存在！")
	}

	err = database.GetDB().Model(&models.SysApis{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"parent_id":   params.ParentId,
			"description": params.Description,
			"method":      params.Method,
			"path":        params.Path,
		}).Error

	return
}

// @Summary ApiDelete
func (self *SystemLogic) ApiDelete(ctx context.Context, id int64) (err error) {
	var count int64
	database.GetDB().Where("api_id", id).Model(&models.SysRoleApis{}).Count(&count)
	if count > 0 {
		return fmt.Errorf("请先删除角色API权限后再操作")
	}

	err = database.GetDB().Delete(&models.SysApis{}, "id = ?", id).Error

	return
}

func (self *SystemLogic) CacheApiInfo() {
	var list []*models.SysApis
	if err := database.GetDB().Model(&models.SysApis{}).Find(&list).Error; err != nil {
		log.Printf("[CacheApiInfo]获取API信息失败！ error: %v", err)
		return
	}

	for _, item := range list {
		apiInfoMap[item.Path] = item
		recordUseApiInfoMap[fmt.Sprintf("%s_%s", item.Path, item.Method)] = item
	}
}

func (self *SystemLogic) GetRecordDescription(path, method string) string {
	info, ok := recordUseApiInfoMap[fmt.Sprintf("%s_%s", path, method)]
	if !ok {
		return ""
	}

	return info.Description
}

// @Summary DictAdd
func (self *SystemLogic) DictAdd(ctx context.Context, params *dto.UpsertDictReq) (err error) {
	data := new(models.SysDicts)
	_ = copier.Copy(data, params)
	err = database.GetDB().Model(&models.SysDicts{}).Create(data).Error

	return
}

// @Summary DictList
func (self *SystemLogic) DictList(ctx context.Context, params *dto.DictListReq) (resp *dto.DictListResp, err error) {
	resp = &dto.DictListResp{}

	query := database.GetDB().Model(&models.SysDicts{}).Order("id desc")
	if params.DictName != "" {
		query.Where("dict_name like ?", params.DictName+"%")
	}
	if params.DictType != "" {
		query.Where("dict_type like ?", params.DictType+"%")
	}
	if params.Status > 0 {
		query.Where("status = ?", params.Status)
	}
	if err = query.Count(&resp.Total).Error; err != nil {
		return
	}

	var list []*models.SysDicts
	if err = query.Scopes(utils.GormPaginate(params.Page, params.PageSize)).Find(&list).Error; err != nil {
		return
	}

	items := make([]*dto.DictInfoResp, 0)
	for _, info := range list {
		item := new(dto.DictInfoResp)
		_ = copier.Copy(item, info)
		item.CreatedAt = info.CreatedAt.UnixMilli()
		items = append(items, item)
	}
	resp.Items = items

	return
}

// @Summary DictUpdate
func (self *SystemLogic) DictUpdate(ctx context.Context, id int64, params *dto.UpsertDictReq) (err error) {
	err = database.GetDB().Model(&models.SysDicts{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"dict_name":  params.DictName,
			"dict_type":  params.DictType,
			"status":     params.Status,
			"item_key":   params.ItemKey,
			"item_value": params.ItemValue,
			"sort":       params.Sort,
			"remark":     params.Remark,
		}).Error

	return
}

// @Summary DictDelete
func (self *SystemLogic) DictDelete(ctx context.Context, id int64) (err error) {
	err = database.GetDB().Delete(&models.SysDicts{}, "id = ?", id).Error

	return
}

// @Summary MenuRouter
func (self *SystemLogic) MenuRouter(ctx context.Context) (resp *dto.MenuResp, err error) {
	var authIds []int64
	database.GetDB().Table("sys_users").
		Select("sys_role_auths.auth_id").
		Joins("JOIN sys_role_auths ON sys_users.role_id = sys_role_auths.role_id").
		Where("sys_users.id = ?", ctx.Value("userId")).
		Pluck("sys_role_auths.auth_id", &authIds)
	if len(authIds) == 0 {
		return
	}

	var list []*models.SysMenus
	if err = database.GetDB().Model(&models.SysMenus{}).
		Where("status = ?", 1).
		Where("id IN ?", authIds).
		Where("type != ?", "BUTTON").
		Order("id asc").
		Find(&list).Error; err != nil {
		return
	}

	items := self.getMenuRouter(list, 0)
	resp = &dto.MenuResp{}
	resp.Items = items
	resp.Total = int64(len(items))

	return
}

// 根据指定的pid获取菜单路由树
func (self *SystemLogic) getMenuRouter(menuList []*models.SysMenus, pid int64) (treeList []dto.RouterResp) {
	for _, v := range menuList {
		if v.ParentId == pid {
			child := self.getMenuRouter(menuList, v.Id)
			node := dto.RouterResp{
				Path:      v.Path,
				Component: v.Component,
				Name:      v.RouteName,
				Meta: dto.RouterMetaResp{
					Icon:     v.Icon,
					Sort:     v.Sort,
					Title:    v.Name,
					AffixTab: v.AffixTab == 1,
				},
			}
			node.Children = child
			treeList = append(treeList, node)
		}
	}

	return treeList
}

// @Summary MenuTree
func (self *SystemLogic) MenuTree(ctx context.Context, params *dto.MenuTreeReq) (resp *dto.MenuResp, err error) {
	query := database.GetDB().Model(&models.SysMenus{})
	if params.Name != "" {
		query.Where("name like ?", params.Name+"%")
	}
	if params.Status > 0 {
		query.Where("status = ?", params.Status)
	}
	var list []*models.SysMenus
	if err = query.Find(&list).Error; err != nil {
		return
	}

	var items []*dto.MenuInfo
	for _, info := range list {
		item := new(dto.MenuInfo)
		_ = copier.Copy(item, info)
		_ = copier.Copy(&item.Meta, info)
		items = append(items, item)
	}

	resp = &dto.MenuResp{}
	resp.Items = items
	resp.Total = int64(len(items))

	return
}

// @Summary MenuAdd
func (self *SystemLogic) MenuAdd(ctx context.Context, params *dto.MenuInfo) (err error) {
	data := new(models.SysMenus)
	_ = copier.Copy(data, params)
	_ = copier.Copy(data, params.Meta)
	err = database.GetDB().Model(&models.SysMenus{}).Create(data).Error

	return
}

// @Summary MenuUpdate
func (self *SystemLogic) MenuUpdate(ctx context.Context, id int64, params *dto.MenuInfo) (err error) {
	err = database.GetDB().Model(&models.SysMenus{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"parent_id":             params.ParentId,
			"path":                  params.Path,
			"status":                params.Status,
			"type":                  params.Type,
			"route_name":            params.RouteName,
			"perm":                  params.Perm,
			"component":             params.Component,
			"sort":                  params.Meta.Sort,
			"icon":                  params.Meta.Icon,
			"name":                  params.Meta.Name,
			"affix_tab":             params.Meta.AffixTab,
			"hide_children_in_menu": params.Meta.HideChildrenInMenu,
			"hide_in_breadcrumb":    params.Meta.HideInBreadcrumb,
			"hide_in_menu":          params.Meta.HideInMenu,
			"hide_in_tab":           params.Meta.HideInTab,
			"keep_alive":            params.Meta.KeepAlive,
		}).Error

	return
}

// @Summary MenuInfo
func (self *SystemLogic) MenuInfo(ctx context.Context, id int64) (resp *dto.MenuInfo, err error) {
	// TODO implement

	return
}

// @Summary MenuDelete
func (self *SystemLogic) MenuDelete(ctx context.Context, id int64) (err error) {
	var count int64
	database.GetDB().Model(&models.SysRoleAuths{}).Where("auth_id", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("请删除角色菜单权限！")
	}
	err = database.GetDB().Delete(&models.SysMenus{}, "id = ?", id).Error

	return
}

// @Summary RecordList
func (self *SystemLogic) RecordList(ctx context.Context, params *dto.RecordListReq) (resp *dto.RecordListResp, err error) {
	resp = &dto.RecordListResp{}

	query := database.GetDB().Model(&models.SysRecords{}).Order("id desc")
	if params.Username != "" {
		query.Where("username like ?", params.Username+"%")
	}
	if len(params.CreateTime) > 0 {
		query.Where("created_at between ? and ?", params.CreateTime[0], params.CreateTime[1])
	}
	if err = query.Count(&resp.Total).Error; err != nil {
		return
	}

	var list []*models.SysRecords
	if err = query.Scopes(utils.GormPaginate(params.Page, params.PageSize)).Find(&list).Error; err != nil {
		return
	}

	items := make([]*dto.RecordInfoResp, 0)
	for _, info := range list {
		item := new(dto.RecordInfoResp)
		_ = copier.Copy(item, info)
		item.CreatedAt = info.CreatedAt.UnixMilli()
		items = append(items, item)
	}
	resp.Items = items

	return
}

// @Summary RecordDelete
func (self *SystemLogic) RecordDelete(ctx context.Context, id int64) (err error) {
	err = database.GetDB().Delete(&models.SysRecords{}, "id = ?", id).Error

	return
}

// @Summary RoleAdd
func (self *SystemLogic) RoleAdd(ctx context.Context, params *dto.UpsertRoleReq) (err error) {
	var has int64
	database.GetDB().Model(&models.SysRoles{}).Where("code = ?", params.Code).Count(&has)
	if has > 0 {
		return fmt.Errorf("角色已存在！")
	}

	err = database.GetDB().Create(&models.SysRoles{
		Name:   params.Name,
		Code:   params.Code,
		Status: params.Status,
		Remark: params.Remark,
		Sort:   params.Sort,
	}).Error

	return
}

// @Summary RoleList
func (self *SystemLogic) RoleList(ctx context.Context, params *dto.RoleListReq) (resp *dto.RoleListResp, err error) {
	resp = &dto.RoleListResp{}

	query := database.GetDB().Model(&models.SysRoles{}).Order("id asc")
	if params.Name != "" {
		query.Where("name like ?", params.Name+"%")
	}
	if params.Status > 0 {
		query.Where("status = ?", params.Status)
	}
	if err = query.Count(&resp.Total).Error; err != nil {
		return
	}

	var list []*models.SysRoles
	if err = query.Scopes(utils.GormPaginate(params.Page, params.PageSize)).Find(&list).Error; err != nil {
		return
	}

	items := make([]*dto.RoleInfoResp, 0)
	for _, info := range list {
		item := new(dto.RoleInfoResp)
		_ = copier.Copy(item, info)
		item.CreatedAt = info.CreatedAt.UnixMilli()
		items = append(items, item)
	}
	resp.Items = items

	return
}

// @Summary RoleInfo
func (self *SystemLogic) RoleInfo(ctx context.Context, id int64) (resp *dto.RoleInfoResp, err error) {
	role := &models.SysRoles{}
	if err = database.GetDB().Preload("RoleAuths").Preload("RoleApis").First(&role, id).Error; err != nil {
		return
	}

	resp = &dto.RoleInfoResp{}
	_ = copier.Copy(resp, role)
	for _, auth := range role.RoleAuths {
		resp.AuthId = append(resp.AuthId, auth.AuthId)
	}
	for _, api := range role.RoleApis {
		resp.ApiId = append(resp.ApiId, api.ApiId)
	}

	return
}

// @Summary RoleUpdate
func (self *SystemLogic) RoleUpdate(ctx context.Context, id int64, params *dto.UpsertRoleReq) (err error) {
	err = database.GetDB().Model(&models.SysRoles{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":   params.Name,
			"code":   params.Code,
			"sort":   params.Sort,
			"status": params.Status,
			"remark": params.Remark,
		}).Error

	return
}

// @Summary RoleAssign
func (self *SystemLogic) RoleAssign(ctx context.Context, id int64, params *dto.AssignRoleReq) (err error) {
	type CasbinInfo struct {
		Path   string
		Method string
	}
	var casbinInfos []CasbinInfo

	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		if err = self.deleteRoleAuth(ctx, tx, id); err != nil {
			return err
		}

		if err = tx.Where("v0 = ?", id).Unscoped().Delete(&gormadapter.CasbinRule{}).Error; err != nil {
			return err
		}

		var roleAuths []models.SysRoleAuths
		for _, authId := range params.AuthId {
			roleAuths = append(roleAuths, models.SysRoleAuths{
				RoleId: id,
				AuthId: authId,
			})
		}
		if err = tx.Model(&models.SysRoleAuths{}).Create(&roleAuths).Error; err != nil {
			return err
		}

		var apis []models.SysApis
		if err = tx.Model(&models.SysApis{}).Where("id in ?", params.ApiId).Find(&apis).Error; err != nil {
			return err
		}

		var roleApis []models.SysRoleApis
		for _, api := range apis {
			roleApis = append(roleApis, models.SysRoleApis{
				RoleId: id,
				ApiId:  api.Id,
			})

			if api.Path != "" {
				casbinInfos = append(casbinInfos, CasbinInfo{
					Path:   api.Path,
					Method: api.Method,
				})
			}
		}

		return tx.Create(&roleApis).Error
	})

	if err == nil {
		// Update casbin policies
		enforcer := middleware.GetEnforcer()
		if enforcer != nil {
			// Remove old policies for this role
			enforcer.RemoveFilteredPolicy(0, fmt.Sprintf("%d", id))

			// Add new policies
			for _, info := range casbinInfos {
				_, err = enforcer.AddPolicy(fmt.Sprintf("%d", id), info.Path, info.Method)
				if err != nil {
					log.Printf("Failed to add casbin policy: %v", err)
				}
			}
		}
	}

	return
}

// @Summary RoleDelete
func (self *SystemLogic) RoleDelete(ctx context.Context, id int64) (err error) {
	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&models.SysRoles{}, "id = ?", id).Error; err != nil {
			return err
		}

		return self.deleteRoleAuth(ctx, tx, id)
	})

	return
}

func (self *SystemLogic) deleteRoleAuth(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	if err = tx.Where("role_id = ?", id).Unscoped().Delete(&models.SysRoleAuths{}).Error; err != nil {
		return
	}

	err = tx.Where("role_id = ?", id).Unscoped().Delete(&models.SysRoleApis{}).Error

	return
}

// @Summary UserInfo
func (self *SystemLogic) UserInfo(ctx context.Context) (resp *dto.UserInfoResp, err error) {
	user := models.SysUsers{}
	if err = database.GetDB().Where("id = ?", ctx.Value("userId")).First(&user).Error; err != nil {
		return
	}
	resp = new(dto.UserInfoResp)
	_ = copier.Copy(resp, user)
	resp.RoleName = user.SysRole.Name
	resp.RealName = user.Username

	// 获取角色信息
	info, err := self.RoleInfo(ctx, user.RoleId)
	if err != nil {
		return
	}

	// 获取权限
	var results []models.SysMenus
	if err = database.GetDB().Model(&models.SysMenus{}).
		Select("perm").
		Where("type = ?", "BUTTON").
		Where("id IN ?", info.AuthId).
		Find(&results).Error; err != nil {
		return
	}
	for _, item := range results {
		resp.Permissions = append(resp.Permissions, item.Perm)
	}

	return
}

func (self *SystemLogic) UserByName(username string) (user models.SysUsers, err error) {
	user = models.SysUsers{}
	err = database.GetDB().Where("username = ?", username).First(&user).Error
	// todo 头像
	// user.Avatar = utils.TransformImageUrl(user.Avatar)

	return
}

// @Summary UserAdd
func (self *SystemLogic) UserAdd(ctx context.Context, params *dto.UpsertUserReq) (err error) {
	if params.Password == "" {
		params.Password = "123456"
	}

	username := params.Username
	if u, _ := self.UserByName(username); u.Id != 0 {
		return fmt.Errorf("用户已存在！")
	}

	salt := utils.RandString(6)
	err = database.GetDB().Create(&models.SysUsers{
		Username: username,
		Salt:     salt,
		Password: utils.MakePasswd(params.Password, salt),
		Nickname: params.Nickname,
		Avatar:   "/uploads/default/logo.png",
		Status:   params.Status,
		Mobile:   params.Mobile,
		Email:    params.Email,
		Remark:   params.Remark,
		RoleId:   params.RoleId,
		CreateBy: ctx.Value("userId").(int64),
	}).Error

	return
}

// @Summary UserList
func (self *SystemLogic) UserList(ctx context.Context, params *dto.UserListReq) (resp *dto.UserListResp, err error) {
	resp = &dto.UserListResp{}

	query := database.GetDB().Model(&models.SysUsers{}).Order("id desc").Preload("SysRole")
	if params.Username != "" {
		query.Where("username like ?", params.Username+"%")
	}
	if params.Status > 0 {
		query.Where("status = ?", params.Status)
	}
	if err = query.Count(&resp.Total).Error; err != nil {
		return
	}

	var list []*models.SysUsers
	if err = query.Scopes(utils.GormPaginate(params.Page, params.PageSize)).Find(&list).Error; err != nil {
		return
	}

	items := make([]*dto.UserInfoResp, 0)
	for _, info := range list {
		item := new(dto.UserInfoResp)
		_ = copier.Copy(item, info)
		item.RoleName = info.SysRole.Name
		item.CreatedAt = info.CreatedAt.UnixMilli()
		items = append(items, item)
	}
	resp.Items = items

	return
}

// @Summary UserUpdate
func (self *SystemLogic) UserUpdate(ctx context.Context, id int64, params *dto.UpsertUserReq) (err error) {
	user := models.SysUsers{}
	if err = database.GetDB().First(&user, id).Error; err != nil {
		return err
	}

	if user.Username != params.Username {
		var count int64
		database.GetDB().Model(&models.SysUsers{}).Where("username = ?", params.Username).Count(&count)
		if count >= 1 {
			return errors.New("用户已存在")
		}
	}

	uMap := map[string]interface{}{
		"username":  params.Username,
		"nickname":  params.Nickname,
		"role_id":   params.RoleId,
		"status":    params.Status,
		"mobile":    params.Mobile,
		"email":     params.Email,
		"remark":    params.Remark,
		"update_by": ctx.Value("userId").(int64),
	}
	if len(params.Password) > 0 {
		newSalt := utils.RandString(6)
		uMap["salt"] = newSalt
		uMap["password"] = utils.MakePasswd(params.Password, newSalt)
	}

	err = database.GetDB().Model(&models.SysUsers{}).Where("id = ?", id).Updates(uMap).Error

	return
}

// @Summary UserDelete
func (self *SystemLogic) UserDelete(ctx context.Context, id int64) (err error) {
	err = database.GetDB().Delete(&models.SysUsers{}, "id = ?", id).Error

	return
}
