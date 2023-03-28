package model

import (
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	UserName          string `gorm:"varchar(32);not null;comment:用户名" json:"name" validate:"required,max=32"`
	NickName          string `gorm:"varchar(32);not null;comment:昵称" json:"nickname" validate:"required,max=32"`
	Password          string `gorm:"size:60;not null;comment:密码的哈希" json:"password"`
	Email             string `gorm:"varchar(256);not null;unique;comment:邮箱" json:"email" validate:"required,max=256,email"`
	Description       string `gorm:"varchar(256);not null;comment:个人简介" json:"description" validate:"required,max=256"`
	Avatar            string `gorm:"varchar(128);not null;comment:头像url" json:"avatar" validate:"required,max=128,uri"`
	VeriToken         string `gorm:"varchar(128);not null;comment:激活token" json:"veriToken"`
	VeriTokenExpireAt int64  `gorm:"BIGINT;not null;comment:token过期时间" json:"veriTokenExpireAt"`
	IsVerified        int    `gorm:"tinyint;default:0;comment:账户是否激活" json:"isVerified"`
}
type Client struct {
	gorm.Model
	ClientId     string `gorm:"varchar(128);not null;comment:客户端ID" json:"clientId"`
	ClientSecret string `gorm:"varchar(128);not null;comment:客户端密码" json:"clientSecret"`
	RedirectUri  string `gorm:"varchar(128);not null;comment:重定向Uri" json:"redirectUri"`
	Scope        string `gorm:"varchar(128);not null;comment:权限元组" json:"scope"`
}
type AuthorizationCode struct {
	gorm.Model
	ClientId string `gorm:"varchar(128);not null;comment:客户端ID" json:"clientId"`
	Code     string `gorm:"varchar(32);not null;comment:授权码" json:"code"`
	Scope    string `gorm:"varchar(128);not null;comment:用户同意的权限元组" json:"scope"`
	ExpireAt int64  `gorm:"int;not null;comment:code过期时间" json:"expireAt"`
	IsUsed   int    `gorm:"int;default:0;comment:code是否已经使用过" json:"isUsed"`
	//AccessToken  string    `gorm:"varchar(128);not null;comment:授权码" json:"AccessToken"`
	//RefreshToken string    `gorm:"varchar(128);not null;comment:授权码" json:"refreshToken"`
}
type AccessToken struct {
	gorm.Model
	ClientId     string `gorm:"varchar(128);not null;comment:客户端ID" json:"clientId"`
	AccessToken  string `gorm:"varchar(128);not null" json:"accessToken"`
	RefreshToken string `gorm:"varchar(128);not null" json:"refreshToken"`
	//RedirectUri     string `gorm:"varchar(128);not null;comment:重定向Uri" json:"redirectUri"`
	AccessExpireAt  int64 `gorm:"int;not null;comment:token过期时间" json:"expireAt"`
	RefreshExpireAt int64 `gorm:"int;not null;comment:token过期时间" json:"refreshExpireAt"`
}

// InitDB 初始化连接并自动迁移
func InitDB() {
	// 从配置文件中读取
	username := setting.DatabaseSettings.UserName
	password := setting.DatabaseSettings.Password
	host := setting.DatabaseSettings.Host
	port := setting.DatabaseSettings.Port
	database := setting.DatabaseSettings.DBName
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
	DB.AutoMigrate(&AccessToken{})
	DB.AutoMigrate(&Client{})
}
