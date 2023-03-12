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
		reg, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
		isMatch := reg.MatchString(email)
		if isMatch == false {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "邮箱地址不合法",
			})
			return
		}
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
	r.POST("signin", func(c *gin.Context) {
		id := c.PostForm("id")
		password := c.PostForm("password")
		regMail, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
		var retUser model.User // return User 检索到的用户
		// 检索用户
		if regMail.MatchString(id) {
			// id 是邮箱
			db.Where("email = ?", id).Find(&retUser)
		} else {
			// id 是用户名
			db.Where("user_name = ?", id).Find(&retUser)
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
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "登录成功",
			})
			return
		}

	})
	r.Run(":8080")
}
