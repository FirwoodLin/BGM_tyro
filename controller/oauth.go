package controller

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

// TODO: client 注冊
// func OauthSignup(c *gin.Context) {
//
// }

// OauthAuthCode 生成 authorization code
func OauthAuthCode(c *gin.Context) {
	// 解析 Header 中参数
	responseType := c.Query("response_type")
	// 非授权码模式，报错
	if responseType != "code" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    5555,
			"message": "目前不支持的授权方式",
		})
	}
	// 解析其他参数
	clientId := c.Query("client_id")
	redirectUri := c.Query("redirect_uri")
	scope := c.Query("scope")
	state := c.Query("state")
	// 如果没有 token/token 过期/无效，重定向进行登录
	if err := auth.JWTTokenCheck(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    5555,
			"message": "请登陆后进行操作，重定向...",
		})
		c.Redirect(http.StatusFound, "/signin")
		return
	}
	// ## 颁发授权码
	// 检查 clientId 和 scope 是否匹配
	if err := model.CheckClient(clientId, scope); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    5555,
			"message": err,
		})
	}
	// 生成授权码
	byteAuthCode := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, byteAuthCode); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    5555,
			"message": "授权码生成错误",
		})
	}
	authCode := base64.URLEncoding.EncodeToString(byteAuthCode)
	// 生成 refresh token
	//byteRefreshToken := make([]byte, 32)
	//if _, err := io.ReadFull(rand.Reader, byteRefreshToken); err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    5555,
	//		"message": "授权码生成错误",
	//	})
	//}
	//refreshCode := base64.URLEncoding.EncodeToString(byteRefreshToken)
	// 存储授权码
	var authCodeStruct = model.AuthorizationCode{
		ClientId:    clientId,
		RedirectUri: redirectUri,
		Scope:       scope,
		Code:        authCode,
		ExpireAt:    time.Now().Add(time.Minute * 10).Unix(),
	}
	if err := model.UpdateAuthCode(&authCodeStruct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "存储授权码失败",
		})
	}
	// 返回授权码
	c.Redirect(http.StatusFound, redirectUri+"?code="+authCode+"?state="+state)
}
