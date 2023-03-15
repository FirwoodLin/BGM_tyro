package model

import (
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	UserName    string `gorm:"varchar(32);not null;comment:用户名" json:"name"`
	NickName    string `gorm:"varchar(32);not null;comment:昵称" json:"nickname"`
	Password    string `gorm:"size:60;not null;comment:密码的哈希" json:"password"`
	Email       string `gorm:"varchar(256);not null;unique;comment:邮箱" json:"email"`
	Description string `gorm:"varchar(256);not null;comment:个人简介" json:"description"`
	Avatar      string `gorm:"varchar(128);not null;comment:头像url" json:"avatar"`
}
type AuthorizationCode struct {
	ClientId    string `gorm:"varchar(128);not null;comment:客户端ID" json:"clientId"`
	RedirectUri string `gorm:"varchar(128);not null;comment:重定向Uri" json:"redirectUri"`
	Scope       int    `gorm:"int;not null;comment:权限元组" json:"scope"`
	Code        string `gorm:"varchar(128);not null;comment:授权码" json:"code"`
}

func InitDB() {
	// 初始化数据库连接
	// 从配置文件中读取
	username := setting.DatabaseSettings.UserName
	password := setting.DatabaseSettings.Password
	host := setting.DatabaseSettings.Host
	port := setting.DatabaseSettings.Port
	//database := setting.DatabaseSettings.DBName
	// TODO:fix dbname readin
	database := "bgm"
	charset := setting.DatabaseSettings.Charset
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	// 建立连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	// 自动建表
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&AuthorizationCode{})
	fmt.Println("migrate success")
}
func CreateUser(u *User) {
	DB.Create(&u)
}
func CheckEmail(email string) error {
	// 检测邮箱的唯一性
	err := DB.Where("email = ?", email).First(&User{}).Error
	return err
}
func CheckId(id string, u *User) {
	regMail, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	var retUser User // return User 检索到的用户
	// 检索用户
	if regMail.MatchString(id) {
		// id 是邮箱
		DB.Where("email = ?", id).Find(&retUser)
	} else {
		// id 是用户名
		DB.Where("user_name = ?", id).Find(&retUser)
	}
}
func UpdateInfo(u *User, username, nickname, description, password string) {
	DB.Model(&u).Select("UserName", "NickName", "Description", "Password").Updates(User{UserName: username, NickName: nickname, Description: description, Password: password})
}

// TODO:检验客户端和权限范围是否匹配

// CheckClient 检验客户端和权限范围是否匹配
func CheckClient(clientId, scope string) error {
	return nil
}
