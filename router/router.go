package router

import (
	"github.com/FirwoodLin/BGM_tyro/auth"
	"github.com/FirwoodLin/BGM_tyro/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) {
	// 注册
	r.POST("/signup", controller.SignUp)
	// 登录
	r.POST("/signin", controller.SignIn)
	// 修改信息
	r.PUT("/user", auth.JWTAuthMiddleWare(), controller.Update)
	// 邮箱激活链接 - 返回接口
	r.GET("/verify", controller.VerifyEmail)
	// OAuth2.0 接口
	oauth := r.Group("/oauth")
	{
		oauth.POST("/client", controller.OauthSignup)
		//oauth.POST("/signup", controller.OauthSignup)
		oauth.GET("/authorization", controller.OauthAuthCode)
		oauth.POST("/token", controller.OauthToken)

	}
	//oidc := r.Group("/oidc")
	//{
	//	oidc.GET()
	//}
	r.Run(":8080")
}
