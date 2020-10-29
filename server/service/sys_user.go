package service

import (
	"encoding/hex"
	"errors"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	"gitee.com/youjy0208/go-common/mhash"
	"gitee.com/youjy0208/go-common/mrand"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// @title    Register
// @description   register, 用户注册
// @auth                     （2020/04/05  20:22）
// @param     u               model.SysUser
// @return    err             error
// @return    userInter       *SysUser

func Register(u model.SysUser) (err error, userInter model.SysUser) {
	var user model.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return errors.New("用户名已注册"), userInter
	}
	// 否则 附加uuid 密码md5简单加密 注册
	u.Sale = mrand.StringAll(8)
	u.Password = hex.EncodeToString(mhash.Md5Byte([]byte(u.Password)))
	u.UUID = uuid.NewV4()
	err = global.GVA_DB.Create(&u).Error
	return err, u
}

// @title    Login
// @description   login, 用户登录
// @auth                     （2020/04/05  20:22）
// @param     u               *model.SysUser
// @return    err             error
// @return    userInter       *SysUser

func Login(u *model.SysUser) (err error, userInter *model.SysUser) {
	var user model.SysUser
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authority").First(&user).Error
	if err != nil{
		return err,&user
	}
	if user.Password == hex.EncodeToString(mhash.Md5Byte([]byte(user.Sale + u.Password))){
		return nil,&user
	}
	return errors.New("密码错误!"), &user
}

// @title    ChangePassword
// @description   change the password of a certain user, 修改用户密码
// @auth                     （2020/04/05  20:22）
// @param     u               *model.SysUser
// @param     newPassword     string
// @return    err             error
// @return    userInter       *SysUser

func ChangePassword(u *model.SysUser, newPassword string) (err error, userInter *model.SysUser) {
	var user model.SysUser
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authority").First(&user).Error
	if err != nil{
		return err,&user
	}
	oldPwd := hex.EncodeToString(mhash.Md5Byte([]byte(user.Sale + u.Password)))
	newPwd := hex.EncodeToString(mhash.Md5Byte([]byte(user.Sale + newPassword)))

	err = global.GVA_DB.Where("username = ? AND password = ?", u.Username, oldPwd).First(&user).Update("password", newPwd).Error
	return err, u
}

// @title    GetInfoList
// @description   get user list by pagination, 分页获取数据
// @auth                      （2020/04/05  20:22）
// @param     info             request.PageInfo
// @return    err              error
// @return    list             interface{}
// @return    total            int

func GetUserInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&model.SysUser{})
	var userList []model.SysUser
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Preload("Authority").Find(&userList).Error
	return err, userList, total
}

// @title    SetUserAuthority
// @description   set the authority of a certain user, 设置一个用户的权限
// @auth                     （2020/04/05  20:22）
// @param     uuid            UUID
// @param     authorityId     string
// @return    err             error

func SetUserAuthority(uuid uuid.UUID, authorityId string) (err error) {
	err = global.GVA_DB.Where("uuid = ?", uuid).First(&model.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

// @title    SetUserAuthority
// @description   set the authority of a certain user, 设置一个用户的权限
// @auth                     （2020/04/05  20:22）
// @param     uuid            UUID
// @param     authorityId     string
// @return    err             error

func DeleteUser(id float64) (err error) {
	var user model.SysUser
	err = global.GVA_DB.Where("id = ?", id).Delete(&user).Error
	return err
}

// @title    SetUserInfo
// @description   set the authority of a certain user, 设置用户信息
// @auth                     （2020/04/05  20:22）
// @param     uuid            UUID
// @param     authorityId     string
// @return    err             error

func SetUserInfo(reqUser model.SysUser) (err error, user model.SysUser) {
	err = global.GVA_DB.Updates(&reqUser).Error
	return err, reqUser
}

// @title    FindUserById
// @description   Get user information by id, 通过id获取用户信息
// @auth                     （2020/04/05  20:22）
// @param     id              int
// @return    err             error
// @return    user            *model.SysUser

func FindUserById(id int) (err error, user *model.SysUser) {
	var u model.SysUser
	err = global.GVA_DB.Where("`id` = ?", id).First(&u).Error
	return err, &u
}

// @title    FindUserByUuid
// @description   Get user information by uuid, 通过uuid获取用户信息
// @auth                     （2020/04/05  20:22）
// @param     uuid            string
// @return    err             error
// @return    user            *model.SysUser

func FindUserByUuid(uuid string) (err error, user *model.SysUser) {
	var u model.SysUser
	if err = global.GVA_DB.Where("`uuid` = ?", uuid).First(&u).Error; err != nil{
		return errors.New("用户不存在"), &u
	}
	return nil, &u
}