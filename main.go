package main

import (
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
	//model.DsnConfigRead()
	//model.TestInit()
	err := initSettings()
	if err != nil {
		//fmt.Println("err in initSettings")
		fmt.Println(err)
		// TODO:错误处理
		//return err
	}
	db := model.InitDB()
	//model.DB
	r := gin.Default()
	// 注册
	r.POST("/signup", func(c *gin.Context) {
		// 数据解析
		name := c.PostForm("name")
		email := c.PostForm("email")
		nickname := c.PostForm("nickname")
		password := c.PostForm("password")
		description := c.PostForm("description")
		avatar := c.PostForm("avatar")
		// 数据检验
		// 姓名校验
		if len(name) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不为空",
			})
			return
		}
		if len(name) > 32 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名长度不合法",
			})
			return
		}
		// 邮箱校验（唯一，合法）
		// TODO：检验邮箱唯一
		reg, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
		isMatch := reg.MatchString(email)
		if isMatch == false {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "邮箱地址不合法",
			})
			return
		}
		//var emailUser model.User
		if err := db.Where("email = ?", email).Find(&model.User{}).Error; err != nil {
			//fmt.Println(err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "邮箱已经注册",
			})
			return
		} else {
		}
		// 简介校验
		if len(description) > 255 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "简介长度不能长于255",
			})
			return
		}
		// 密码校验（合法）
		if len(password) < 8 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码长度不小于8",
			})
			return
		}
		// 存储到数据库
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码加密错误",
			})
			return
		}
		newUser := model.User{
			UserName:    name,
			NickName:    nickname,
			Password:    string(encryptedPassword),
			Email:       email,
			Description: description,
			Avatar:      avatar,
		}
		db.Create(&newUser)
		token, err := auth.GenToken(name, password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    5555,
				"message": "token生成失败",
				//"token":""
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "注册成功",
			"token":   token,
			//"token":""
		})
	})
	// 登录
	r.POST("/signin", func(c *gin.Context) {
		//fmt.Println("in signin")
		id := c.PostForm("id")
		password := c.PostForm("password")
		regMail, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
		var retUser model.User // return User 检索到的用户
		// 检索用户
		if regMail.MatchString(id) {
			// id 是邮箱
			//fmt.Println("mail signin")
			db.Where("email = ?", id).Find(&retUser)
		} else {
			// id 是用户名
			//fmt.Println("name signin")
			//fmt.Println(id)
			db.Where("user_name = ?", id).Find(&retUser)
			//fmt.Println(retUser)
		}
		// 没有检索到
		if retUser == (model.User{}) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名或密码无效",
			})
			return
		}
		isValid := bcrypt.CompareHashAndPassword([]byte(retUser.Password), []byte(password))
		if isValid != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名或密码无效",
			})
			return
		} else {
			token, err := auth.GenToken(id, password)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    5555,
					"message": "token生成失败",
					//"token":""
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "登录成功",
				"token":   token,
				//"token":
			})
			return
		}

	})
	//修改信息
	r.PUT("/user", auth.JWTAuthMiddleWare(), func(c *gin.Context) {
		//c.JSON(http.StatusOK, time.Now().Unix())
		rawid := c.PostForm("id")
		username := c.PostForm("username")
		nickname := c.PostForm("nickname")
		description := c.PostForm("description")
		var user model.User
		id, _ := strconv.Atoi(rawid)
		user.ID = uint(id)
		err := db.Model(&user).Select("UserName", "NickName", "Description").Updates(model.User{UserName: username, NickName: nickname, Description: description})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    "5555",
				"message": "更新信息成功",
			})
		}
	})

	//r.PUT("/user")
	r.Run(":8080")
}
func signinHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    gin.H{"username": username},
	})
}
func initSettings() error {
	//vp, err := setting.NewSetting()
	vp := viper.New()
	vp.SetConfigFile("config.yaml") /// ./config/
	if err := vp.ReadInConfig(); err != nil {
		return err
		//fmt.Println(err)
	}
	if err := vp.UnmarshalKey("mysql", &setting.DatabaseSettings); err != nil {
		// if err := vp.UnmarshalKey("mysql.username", &s); err != nil {
		//fmt.Println(err)
		return err
	}
	if err := vp.UnmarshalKey("JWT", &setting.JWTSettings); err != nil {
		// if err := vp.UnmarshalKey("mysql.username", &s); err != nil {
		//fmt.Println(err)
		return err
	}

	//if err != nil {
	//	return err
	//}

	//err = setting.ReadSection("JWT", &config.JWT)
	//err = vp.UnmarshalKey("mysql", setting.DatabaseSettings)
	//if err != nil {
	//	return err
	//}
	//err = vp.UnmarshalKey("JWT", setting.JWTSettings)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(setting.JWTSettings)
	fmt.Println(setting.DatabaseSettings)
	//fmt.Printf("main-init:%T %v\n", setting.JWTSettings.Secret, setting.JWTSettings.Secret)
	return nil
	//return err
}
