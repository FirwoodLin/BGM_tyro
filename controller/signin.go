package controller

import (
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
)

//var db *gorm.DB = model.DB

func SignUp(c *gin.Context) {
	// ### 数据解析 ###
	name := c.PostForm("name")

	email := c.PostForm("email")
	nickname := c.PostForm("nickname")
	password := c.PostForm("password")
	description := c.PostForm("description")
	avatar := c.PostForm("avatar")
	// ### 数据检验 ###
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
	// 邮箱校验 - 邮箱格式
	isMatch := reg.MatchString(email)
	if isMatch == false {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "邮箱地址不合法",
		})
		return
	}
	// 邮箱校验 - 唯一性
	if err := model.CheckEmail(email); err != gorm.ErrRecordNotFound {
		//fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "邮箱已经注册",
		})
		c.Abort()
		return
	} else {
	}
	// 简介校验 - 长度
	if len(description) > 255 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "简介长度不能长于255",
		})
		return
	}
	// 密码校验 - 长度
	if len(password) < 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码长度不小于8",
		})
		return
	}
	// 密码 - 加密
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码加密错误",
		})
		return
	}
	// 存储到数据库
	newUser := model.User{
		UserName:    name,
		NickName:    nickname,
		Password:    string(encryptedPassword),
		Email:       email,
		Description: description,
		Avatar:      avatar,
	}
	model.CreateUser(&newUser)
	//token, err := auth.GenToken(name, password)
	token, err := auth.GenToken(name)

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
	c.Abort()
}
func SignIn(c *gin.Context) {
	//fmt.Println("in signin")
	id := c.PostForm("id")
	password := c.PostForm("password")
	// 检索用户
	var retUser model.User // 检索到的用户
	model.CheckId(id, &retUser)
	fmt.Println(retUser)
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
		// 密码-用户名 检验失败
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户名或密码无效",
		})
		return
	} else {
		// 生成 token
		token, err := auth.GenToken(id)
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
		})
		return
	}
}
func Update(c *gin.Context) {
	rawid := c.PostForm("id")
	password := c.PostForm("password")
	username := c.PostForm("username")
	nickname := c.PostForm("nickname")
	description := c.PostForm("description")
	var user model.User
	id, _ := strconv.Atoi(rawid)
	user.ID = uint(id)
	model.UpdateInfo(&user, username, nickname, description, password)
	//if err != nil {
	c.JSON(http.StatusOK, gin.H{
		"code":    "5555",
		"message": "更新信息成功",
	})
}
