package model

import (
	"errors"
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
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
	gorm.Model
	ClientId    string    `gorm:"varchar(128);not null;comment:客户端ID" json:"clientId"`
	RedirectUri string    `gorm:"varchar(128);not null;comment:重定向Uri" json:"redirectUri"`
	Scope       string    `gorm:"varchar(128);not null;comment:权限元组" json:"scope"`
	Code        string    `gorm:"varchar(128);not null;comment:授权码" json:"code"`
	ExpireAt    time.Time `gorm:"datetime(3);not null;comment:过期时间" json:"expireAt"`
}

// InitDB 初始化连接并自动迁移
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

// CheckClient 检验 clientId 和 scope 范围是否匹配
func CheckClient(clientId, scope string) error {
	var client AuthorizationCode
	// 检查 client 是否注册
	if err := DB.Where("client_id = ?", clientId).Find(&client).Error; err != nil {
		return errors.New("client not found")
	}
	// 检查 scope 是否已经授权
	var authedScopeMap map[string]int
	for _, v := range strings.Split(client.Scope, ",") {
		authedScopeMap[v] = 1
	}
	for _, v := range strings.Split(scope, ",") {
		if authedScopeMap[v] != 1 {
			return errors.New("unauthed scope")
		}
	}
	return nil
}
