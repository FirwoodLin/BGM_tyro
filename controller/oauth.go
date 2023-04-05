package controller

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

// OauthSignup TODO: client 注冊
// OauthSignup 客户端注册
func OauthSignup(c *gin.Context) {

}

// OauthAuthCode 生成 authorization code
func OauthAuthCode(c *gin.Context) {
	// 解析 Header 中参数
	responseType := c.Query("response_type")
	// 非授权码模式，报错
	if responseType != "code" {
		log.Printf("controller-OauthAuthCode:responseTypeErr:%v", responseType)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "目前不支持的授权方式",
		})
		c.Abort()
		return
	}
	// 解析其他参数
	clientId := c.Query("client_id")
	redirectUri := c.Query("redirect_uri")
	scope := c.Query("scope")
	state := c.Query("state")
	// 如果没有 token/token 过期/无效，重定向进行登录
	// TODO:改为中间件验证token，并使用 c.set 设置 userid
	if err := auth.JWTTokenCheck(c); err != nil {
		log.Printf("controller-OauthAuthCode重定向到登录页:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "请登陆后进行操作，重定向...",
		})
		// TODO:重定向至“登录页”-登录页在哪
		c.Redirect(http.StatusFound, "/signin")
		return
	}
	// 检查 clientId 和先前注册的 scope、redirect_url 是否匹配
	if err := model.CheckScope(clientId, scope, redirectUri); err != nil {
		log.Printf("controller-OauthAuthCode:%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "scope-redirect_url不匹配",
		})
		c.Abort()
		return
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
	// 应该在 JWT 鉴权后获取 user_id
	userId, _ := c.Get("user_id")
	authCode := base64.URLEncoding.EncodeToString(byteAuthCode)
	// 存储授权码
	var authCodeStruct = model.AuthorizationCode{
		UserId:   userId.(uint),
		ClientId: clientId,
		//RedirectUri: redirectUri,
		Scope: scope,
		Code:  authCode,
		// TODO:将有效期统一成全局变量
		ExpireAt: time.Now().Add(time.Minute * 10).Unix(),
	}
	log.Printf("controller-OauthAuthCode:authCodeStruct-%v", authCodeStruct)
	if err := model.UpdateAuthCode(&authCodeStruct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "存储授权码失败",
		})
		c.Abort()
		return
	}
	// 返回授权码
	log.Printf("controller-OauthAuthCode:重定向：%s", fmt.Sprintf(redirectUri+"?code="+authCode+"&state="+state))
	c.Redirect(http.StatusFound, redirectUri+"?code="+authCode+"&state="+state)
}

// OauthToken 根据请求类型code 或者 refresh，返回 access token 和 refresh token；或者刷新token
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
		log.Printf("#ERR#controller-OauthToken:授权类型无效:%v", grantType)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": "grant_type 错误",
		})
		c.Abort()
	}
}

// OauthTokenAuthCode 授权码模式,根据 code 颁发 token
func OauthTokenAuthCode(c *gin.Context) {
	code := c.PostForm("code")
	clientId := c.PostForm("client_id")
	// 检查 code 是否存在、未使用过、未过期
	var authCodeStruct = model.AuthorizationCode{
		Code:     code,
		ClientId: clientId,
	}
	if err := model.CheckCode(&authCodeStruct); err != nil {
		log.Printf("#ERR#controller-OauthTokenAuthCode:code错误 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	// 检查 client secret 是否正确,检查 redirect_uri
	clientSecret := c.PostForm("client_secret")
	redirectUri := c.PostForm("redirect_uri")
	if err := model.CheckSecret(clientId, clientSecret, redirectUri); err != nil {
		log.Printf("#ERR#controller-OauthTokenAuthCode:客户端检验错误 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	// 生成 refresh token，并编码
	byteRefreshToken := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, byteRefreshToken); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    5555,
			"message": "refresh 生成错误",
		})
		c.Abort()
		return
	}
	refreshToken := base64.URLEncoding.EncodeToString(byteRefreshToken)
	// 生成 access token，并编码
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
		UserId:       authCodeStruct.UserId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClientId:     clientId,
		//RedirectUri:     redirectUri,
		// TODO:有效期存储到配置文件中
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
	scope, err := model.GetScope(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": "查询 scope 出错",
		})
		return
	}
	// OIDC : 如果请求 id token，进行生成，并返回
	//TODO
	//if strings.Contains(scope, "openid") {
	//	idToken, err := OidcGenIdToken(scope, clientId)
	//	if err != nil {
	//
	//	} else {
	//
	//	}
	//}
	// 返回普通请求
	c.JSON(http.StatusOK, gin.H{
		"token_type":    "bearer",
		"scope":         scope,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		//"id_token":idToken,
		"expires_at": time.Now().Add(time.Hour).Unix(),
	})
	return
}

// OauthTokenRefresh 根据 refresh token 颁发 access token
func OauthTokenRefresh(c *gin.Context) {
	refresh := c.PostForm("refresh_token")
	clientId := c.PostForm("client_id")
	secret := c.PostForm("client_secret")
	var accessTokenStruct = model.AccessToken{
		RefreshToken: refresh,
		ClientId:     clientId,
	}
	// 检查 refresh 是否存在、未过期
	if err := model.CheckRefresh(&accessTokenStruct, secret); err != nil {
		log.Printf("#ERR#controller-OauthTokenRefresh:refresh token错误 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	// 更新 refresh 和 access token
	// 生成 refresh token，并编码 ### 决定不更新 Refresh Token
	//byteRefreshToken := make([]byte, 32)
	//if _, err := io.ReadFull(rand.Reader, byteRefreshToken); err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    5555,
	//		"message": "refresh 生成错误",
	//	})
	//	c.Abort()
	//	return
	//}
	//refreshToken := base64.URLEncoding.EncodeToString(byteRefreshToken)
	// 生成 access token，并编码
	byteAccessToken := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, byteAccessToken); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    5555,
			"message": "access 生成错误",
		})
		return
	}
	accessToken := base64.URLEncoding.EncodeToString(byteAccessToken)
	// 持久化更新
	accessExpireAt := time.Now().Add(time.Hour * 2).Unix()
	var tokenStruct = model.AccessToken{
		UserId:      accessTokenStruct.UserId,
		ClientId:    clientId,
		AccessToken: accessToken,
		//RefreshToken:   refreshToken,
		AccessExpireAt: accessExpireAt,
	}
	if err := model.UpdateToken(&tokenStruct); err != nil {
		log.Printf("controller-OauthTokenRefresh:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5555,
			"message": fmt.Sprintf("更新access token出错：%v", err),
		})
		c.Abort()
		return
	}
	// 成功生成 返回
	c.JSON(http.StatusOK, gin.H{
		"code":        5555,
		"message":     "成功生成",
		"accessToken": accessToken,
		"expireAt":    accessExpireAt,
	})
}
