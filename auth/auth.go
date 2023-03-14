package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var MySecret = []byte("测试秘钥1")

const TokenExpireDuration = time.Hour * 2

// 造轮子 - 生成 token
func GenToken(username, password string) (string, error) {
	// 拼装 claim
	claim := MyClaims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bgm",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenClaims.SignedString(MySecret)
	return token, err
}

// 造轮子 - 解析 token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
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
		c.Set("username", token.Username)
		c.Next()
	}
}
