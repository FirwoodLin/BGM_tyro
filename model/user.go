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
	fmt.Printf("## in CheckId,id:%v\n", id)
	if regMail.MatchString(id) {
		// id 是邮箱
		fmt.Println("## checking email")
		DB.Where("email = ?", id).Find(u)
		fmt.Printf("## after check mail %v\n", u)
	} else {
		// id 是用户名
		fmt.Println("## checking username")
		DB.Where("user_name = ?", id).Find(u)
		fmt.Printf("## after check username %v\n", u)

		//fmt.Println(retUser)
	}
}

// UpdateInfo 更新用户信息
func UpdateInfo(u *User, username, nickname, description, password string) {
	DB.Model(&u).Select("UserName", "NickName", "Description", "Password").Updates(User{UserName: username, NickName: nickname, Description: description, Password: password})
}

// CheckVeri 检测激活链接是否有效- 有效则成功验证；过期删除；无效报错
func CheckVeri(user *User) error {
	var realUser User
	realUser.ID = user.ID
	//err := DB.Find(&realUser).Error
	err := DB.Where("id = ?", user.ID).First(&realUser).Error
	fmt.Printf("in CheckVeri user %v;err:%v\n", user, err)

	fmt.Printf("in CheckVeri real %v;err:%v\n", realUser, err)
	if err != nil {
		return err
	}
	if realUser.VeriToken != user.VeriToken {
		fmt.Printf("real vs 传入：%v\nvs%v\n", realUser.VeriToken, user.VeriToken)
		return errors.New("token错误")
	}
	if realUser.VeriTokenExpireAt < time.Now().Unix() {
		DB.Delete(&user)
		return errors.New("token过期，请重新注册")
	}
	// 检查通过，用户验证成功；设置账户激活状态
	fmt.Printf("CheckVeri %v success\n", user)
	DB.Model(&user).Update("is_verified", 1)
	return nil
}

// CheckUserVerified 检查用户是否激活
func CheckUserVerified(u *User) error {
	if err := DB.Find(u).Error; err != nil {
		return err
	}
	if u.IsVerified == 0 {
		return errors.New("未激活")
	}
	return nil
}
