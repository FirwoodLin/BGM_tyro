package controller

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
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
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "请登陆后进行操作，重定向...",
		})
		c.Redirect(http.StatusFound, "/signin")
		return
	}
	// 检查 clientId 和先前注册的 scope、redirect_url 是否匹配
	if err := model.CheckScope(clientId, scope, redirectUri); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
		return
	}
	authCode := base64.URLEncoding.EncodeToString(byteAuthCode)
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
		return
	}
	// 返回授权码
	c.Redirect(http.StatusFound, redirectUri+"?code="+authCode+"?state="+state)
}

// OauthToken 根据 code 返回 access token 和 refresh token；或者刷新token
func OauthToken(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	// 检查 grant_type
	switch grantType {
	case "authorization_code":
		OauthTokenAuthCode(c)
	case "refresh_token":
		OauthTokenRefresh(c)
	default:
		// 无效类型进行报错
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "grant_type 错误",
		})
	}
}

// OauthTokenAuthCode 授权码模式颁发 token
func OauthTokenAuthCode(c *gin.Context) {
	code := c.PostForm("code")
	redirectUri := c.PostForm("redirect_uri")
	// 检查 code 是否使用过、是否过期
	if err := model.CheckCode(code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": err,
		})
	}
	// 检查 client id 是否注册、与 code 是否一致;检查 redirect_uri
	authHeader := c.Request.Header.Get("Authorization")
	authList := strings.SplitN(authHeader, " ", 2)
	if strings.ToLower(authList[0]) != "bearer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "header 鉴权类型错误",
		})
		return
	}
	clientCreds, err := base64.StdEncoding.DecodeString(authList[1])
	if err != nil {
		// 解码失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "header中不包含有效鉴权信息",
		})
		return
	}
	creds := strings.Split(string(clientCreds), ":")
	clientID := creds[0]
	clientSecret := creds[1]
	err = model.CheckSecret(clientID, clientSecret, redirectUri)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "Client 信息错误",
		})
		return
	}
	// 生成 refresh token
	byteRefreshToken := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, byteRefreshToken); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    5555,
			"message": "refresh 生成错误",
		})
	}
	refreshToken := base64.URLEncoding.EncodeToString(byteRefreshToken)
	// 生成 access token
	byteAccessToken := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, byteAccessToken); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    5555,
			"message": "access 生成错误",
		})
		return
	}
	accessToken := base64.URLEncoding.EncodeToString(byteAccessToken)
	// 持久化
	var accessTokenStruct = model.AccessToken{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		ClientId:        clientID,
		RedirectUri:     redirectUri,
		AccessExpireAt:  time.Now().Add(time.Hour * 2).Unix(),
		RefreshExpireAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	if err := model.CreateToken(&accessTokenStruct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "存储 token 出错",
		})
		return
	}
	// 查询授权范围
	scope, err := model.GetScope(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "查询 scope 出错",
		})
		return
	}
	// 返回
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"token_type":    "bearer",
		"expires_at":    time.Now().Add(time.Hour).Unix(),
		"refresh_token": refreshToken,
		"scope":         scope,
	})
	return
}

// OauthTokenRefresh 根据 refresh token 颁发新 access token
func OauthTokenRefresh(c *gin.Context) {
	// TODO:完成令牌的刷新
}
