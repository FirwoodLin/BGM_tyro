package main

import (
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

func main() {
	//model.DsnConfigRead()
	//model.TestInit()
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
		avatar := c.PostForm("avator")
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
		//db.Take()
		isExist := db.Where("email = ?", email).Find(&model.User{})
		if isExist != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "邮箱已经注册",
			})
			return
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
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "注册成功"})
	})
	// 登录
	//r.POST("signin",func)
	r.Run(":8080")
}
