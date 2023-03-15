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
	UserName    string `gorm:"varchar(32);not null;comment:用户名" json:"name"`
	NickName    string `gorm:"varchar(32);not null;comment:昵称" json:"nickname"`
	Password    string `gorm:"size:60;not null;comment:密码的哈希" json:"password"`
	Email       string `gorm:"varchar(256);not null;unique;comment:邮箱" json:"email"`
	Description string `gorm:"varchar(256);not null;comment:个人简介" json:"description"`
	Avatar      string `gorm:"varchar(128);not null;comment:头像url" json:"avatar"`
}
type DsnConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
}

//	func DsnConfigRead() DsnConfig {
//		// 读取数据库连接配置信息
//		// 打开配置文件，延时关闭
//		file, _ := os.Open("./model/config.json")
//		fmt.Println(file)
//
//		defer file.Close()
//		// 创建解码器
//		//decoder := json.NewDecoder(file)
//		//dsnconf := DsnConfig{}
//		//Decode从输入流读取下一个json编码值并保存在v指向的值里
//		//err := decoder.Decode(&dsnconf)
//		//err = json.Unmarshal(dsnconf, &dsnconf)
//		if err != nil {
//			fmt.Println("Error:", err)
//		}
//		fmt.Println(dsnconf)
//		return dsnconf
//	}
func InitDB() *gorm.DB {
	// 初始化数据库连接
	// 从配置文件中读取
	//DsnElement := DsnConfigRead()
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
	fmt.Println("migrate success")
	return DB
	//DB = Db
}
func TestInit() {
	// 测试 InitDB 函数
	fmt.Println("Start TestInit")
	db := InitDB()
	s1 := &User{
		UserName:    "itsaiddddd",
		NickName:    "nicknamefadfdasf",
		Email:       "tan@163.com.hkj",
		Description: "desc test",
		Password:    "fadsfas",
	}
	db.Create(&s1)
	fmt.Println(s1)
}
func Create(db *gorm.DB) {

}
