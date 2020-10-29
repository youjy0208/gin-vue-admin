package datas

import (
	"encoding/hex"
	"gitee.com/youjy0208/go-common/mhash"
	"gitee.com/youjy0208/go-common/mrand"
	"time"

	"gin-vue-admin/model"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var Users = []model.SysUser{
	{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, UUID: uuid.NewV4(), Username: "admin", Password: "123456", NickName: "超级管理员", HeaderImg: "http://qmplusimg.henrongyi.top/1571627762timg.jpg", AuthorityId: "888",GroupId:0},
	{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, UUID: uuid.NewV4(), Username: "test", Password: "123456", NickName: "测试", HeaderImg: "http://qmplusimg.henrongyi.top/1572075907logo.png", AuthorityId: "9528",GroupId:0},
}

func InitSysUser(db *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, user := range Users{
			user.Sale = mrand.StringAll(8);
			user.Password = hex.EncodeToString(mhash.Md5Byte([]byte((user.Sale + user.Password))))
			if tx.Create(&user).Error != nil { // 遇到错误时回滚事务
				return err
			}
		}
		return nil
	})
}
