package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// CreateUser 插入用户
func CreateUser(u *User) {
	DB.Create(&u)
}

// CheckEmail 检查邮箱唯一性
func CheckEmail(email string) error {
	err := DB.Where("email = ?", email).First(&User{}).Error
	return err
}

// CheckId 检查登陆时用户名/邮箱和密码是否匹配
func CheckId(id string, u *User) {
	regMail, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	// 检索用户
	if regMail.MatchString(id) {
		// id 是邮箱
		fmt.Println("checking email")
		DB.Where("email = ?", id).Find(&u)
	} else {
		// id 是用户名
		fmt.Println("checking username")
		DB.Where("user_name = ?", id).Find(&u)
		//fmt.Println(retUser)
	}
}

// UpdateInfo 更新用户信息
func UpdateInfo(u *User, username, nickname, description, password string) {
	DB.Model(&u).Select("UserName", "NickName", "Description", "Password").Updates(User{UserName: username, NickName: nickname, Description: description, Password: password})
}

// CheckVeri 检测激活链接是否有效；过期删除；无效报错
func CheckVeri(user *User) error {
	var realUser User
	realUser.ID = user.ID
	err := DB.Find(&realUser).Error
	if err != nil {
		return errors.New("找不到此用户")
	}
	if realUser.VeriToken != user.VeriToken {
		return errors.New("token错误")
	}
	if realUser.VeriTokenExpireAt > time.Now().Unix() {
		DB.Delete(&realUser)
		return errors.New("token过期，请重新注册")
	}
	// 检查通过，用户验证成功
	DB.Model(&realUser).Update("is_verified", 1)
	return nil
}
