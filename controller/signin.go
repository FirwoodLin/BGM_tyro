package controller

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// SignUp 注册用户
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "用户名不为空",
		})
		return
	}
	if len(name) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "用户名长度不合法",
		})
		return
	}

	// 简介校验 - 长度
	if len(description) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "简介长度不能长于255",
		})
		return
	}
	// 密码校验 - 长度
	if len(password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "密码长度不小于8",
		})
		return
	}
	// 邮箱校验（唯一，合法，有效）
	reg, _ := regexp.Compile("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	// 邮箱校验 - 合法
	isMatch := reg.MatchString(email)
	if isMatch == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "邮箱地址不合法",
		})
		return
	}
	// 邮箱校验 - 唯一性
	if err := model.CheckEmail(email); err != gorm.ErrRecordNotFound {
		//fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "邮箱已经注册",
		})
		c.Abort()
		return
	}
	// 密码 - 加密
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    422,
			"message": "密码加密错误",
		})
		return
	}
	// 邮箱校验 - 有效 - 发送激活链接
	byteVeriSecret := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, byteVeriSecret); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    5555,
			"message": "激活 token 生成错误",
		})
	}
	veriSecret := base64.URLEncoding.EncodeToString(byteVeriSecret)
	// 存储到数据库
	newUser := model.User{
		UserName:          name,
		NickName:          nickname,
		Password:          string(encryptedPassword),
		Email:             email,
		Description:       description,
		Avatar:            avatar,
		VeriToken:         veriSecret,
		IsVerified:        0,
		VeriTokenExpireAt: time.Now().Add(time.Minute * 15).Unix(),
	}
	id := model.CreateUser(&newUser)
	// 发送激活邮件
	err = SendEmail(newUser, veriSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "邮件发送错误",
		})
	}
	// 生成 token
	token, err := auth.GenToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "token生成失败",
			//"token":""
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"id":      newUser.ID,
		"token":   token,
	})
	c.Abort()
}

// SignIn 登录
func SignIn(c *gin.Context) {
	//fmt.Println("in signin")
	rawId := c.PostForm("id")
	password := c.PostForm("password")
	// 检索用户用户名/密码
	var retUser model.User // 检索到的用户
	//fmt.Println(retUser)
	//fmt.Printf("before search:%v\n", retUser)

	model.CheckId(rawId, &retUser)
	//fmt.Printf("after search:%v\n", retUser)
	// 检索用户用户名/密码 - 没有检索到
	if retUser.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "用户名或密码无效",
		})
		return
	}
	// 检验密码正确性
	isValid := bcrypt.CompareHashAndPassword([]byte(retUser.Password), []byte(password))
	if isValid != nil {
		// 密码-用户名 检验失败
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "用户名或密码无效",
		})
		return
	}
	// 检验是否激活
	if model.CheckUserVerified(&retUser) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "用户尚未激活，请前往邮箱激活",
		})
		return
	}
	// 生成 token
	token, err := auth.GenToken(retUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "token生成失败",
		})
		return
	}
	// 返回token和ID；更新 token

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"id":      retUser.ID,
		"token":   token,
	})
	return
}

// Update 更新用户数据
func Update(c *gin.Context) {
	username := c.PostForm("username")
	if tokenUserName, ok := c.Get("username"); !ok || username != tokenUserName {
		fmt.Printf("two username:%v-%v\n", tokenUserName, username)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "token 用户名与请求用户名不一致",
		})
		c.Abort()
		return
	}

	rawid := c.PostForm("id")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")
	description := c.PostForm("description")
	var user model.User
	id, _ := strconv.Atoi(rawid)
	user.ID = uint(id)
	// 密码转换、加密
	if len(password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    422,
			"message": "密码长度不小于8",
		})
		return
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    422,
			"message": "密码加密错误",
		})
		return
	}
	model.UpdateInfo(&user, username, nickname, description, string(encryptedPassword))
	//if err != nil {
	c.JSON(http.StatusOK, gin.H{
		"code":    5555,
		"message": "更新信息成功",
	})
}
