package controller

import (
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TODO: client 注冊
// func OauthSignup(c *gin.Context) {
//
// }

// OauthGrant 生成授权码
func OauthGrant(c *gin.Context) {
	// 解析 Header 中参数
	responseType := c.Query("response_type")
	// 非授权码模式，报错
	if responseType != "code" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    5555,
			"message": "目前不支持的授权方式",
		})
	}
	clientId := c.Query("client_id")
	redirectUri := c.Query("redirect_uri")
	scope := c.Query("scope")
	state := c.Query("state")
	// 如果没有 token/token 过期/无效，重定向进行登录，而后重定向到 auth
	// TODO:重定向回到 auth
	if err := auth.JWTTokenCheck(c); err != nil {
		return
	}
	// 颁发授权码
	// 检查 clientId 和 scope 是否匹配
	if err := model.CheckClient(clientId, scope); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    5555,
			"message": err,
		})
	}
	// 生成授权码
	code, err := auth.GenToken("clientId")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    5555,
			"message": "生成 code 错误",
		})
	}
	// 返回授权码
	c.Redirect(http.StatusTemporaryRedirect, redirectUri+"?code="+code+"?state="+state)
	// TODO:如何返回授权码

}
