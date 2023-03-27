package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

type MyClaims struct {
	UserId uint `json:"id"`
	//Username string `json:"username"`
	//Password string `json:"password"`
	jwt.StandardClaims
	//Audience  string `json:"aud,omitempty"`
	//ExpiresAt int64  `json:"exp,omitempty"`
	//Id        string `json:"jti,omitempty"`
	//IssuedAt  int64  `json:"iat,omitempty"`
	//Issuer    string `json:"iss,omitempty"`
	//NotBefore int64  `json:"nbf,omitempty"`
	//Subject   string `json:"sub,omitempty"`
}

var MySecret = []byte("测试秘钥1")

//var MySecret = []byte(setting.JWTSettings.Secret)

const TokenExpireDuration = time.Hour * 2

// GenToken 生成 token
func GenToken(id uint) (string, error) { //func GenToken(username, password string) (string, error) {
	// 拼装 claim
	claim := MyClaims{
		id,
		//password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bgm",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenClaims.SignedString(MySecret)
	return token, err
}

// ParseToken 验证 token 的有效性
func ParseToken(tokenString string) (*MyClaims, error) {
	log.Printf("model-ParseToken:%v", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	log.Printf("model-ParseToken:JWT token after parse:%v\n", token)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// func CheckToken()

// JWTAuthMiddleWare 在 update 中使用；解析 context 中的 token 字段，验证有效性
func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// API 规定：token 在请求头中
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 如果 Header 中不包含 Token
			c.JSON(http.StatusOK, gin.H{
				"code":    5555,
				"message": "请求头中 auth 为空",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code":    5555,
				"message": "请求头中 auth 有误",
			})
			c.Abort()
			return
		}
		token, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    5555,
				"message": "Token 验证失败",
			})
			c.Abort()
			return
		}
		c.Set("user_id", token.UserId)
		c.Next()
	}
}

// JWTTokenCheck 检查登录产生的 token;在Oauth中使用；验证不成功则重定向；返回值为验证结果
func JWTTokenCheck(c *gin.Context) error {
	// API 规定：token 在请求头中
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.Set("redirect_url", "/authorization")
		c.Redirect(http.StatusFound, "/signin")
		return errors.New("请求头中 auth 为空")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.Set("redirect_url", "/authorization")
		c.Redirect(http.StatusFound, "/signin")
		return errors.New("请求头中 auth 有误")
	}
	token, err := ParseToken(parts[1])
	//log.Printf("model-JWTTokenCheck:parts[1]:%v", parts[1])
	if err != nil {
		log.Printf("model-JWTTokenCheck:err:%v", err)
		c.Set("redirect_url", "/authorization")
		c.Redirect(http.StatusFound, "/signin")
		return errors.New("token 验证失败")
	}
	c.Set("user_id", token.UserId)
	return nil
}
