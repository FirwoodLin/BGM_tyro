package model

import (
	"errors"
	"log"
	"strings"
	"time"
)

// CheckScope 检验 clientId 和 scope 范围是否匹配
func CheckScope(clientId, scope, redirectUri string) error {
	var client Client
	// 检查 client 是否注册
	if err := DB.Where("client_id = ?", clientId).Find(&client).Error; err != nil {
		log.Printf("model-CheckScope:client not found\n")
		return errors.New("client not found")
	}
	// 检查 scope 是否已经授权
	var authedScopeMap = make(map[string]int)
	for _, v := range strings.Split(client.Scope, ",") {
		authedScopeMap[v] = 1
	}
	for _, v := range strings.Split(scope, ",") {
		if authedScopeMap[v] != 1 {
			log.Printf("model-CheckScope:unauthed scope\n")
			log.Printf("scope-diff:db:%v;passin:%v", authedScopeMap, scope)
			return errors.New("unauthed scope")
		}
	}
	// 检查 redirect_uri 是否一致
	if client.RedirectUri != redirectUri {
		log.Printf("model-CheckScope:重定向 Uri 与注册不一致\n")
		log.Printf("uri-diff:db:%v;passin:%v", client.RedirectUri, redirectUri)
		return errors.New("重定向 Uri 与注册不一致")
	}
	// 检查通过
	return nil
}

// UpdateAuthCode 发布 authorization code 之后持久化到数据库
func UpdateAuthCode(authCode *AuthorizationCode) error {
	if err := DB.Create(&authCode).Error; err != nil {
		return err
	}
	return nil
}

// CheckCode 检验 code 是否存在、未使用过、未过期
func CheckCode(authCode *AuthorizationCode) error {
	//var authCode AuthorizationCode
	err := DB.Where("client_id = ?", authCode.ClientId).Find(&authCode).Error
	if err != nil {
		return errors.New("授权码无效-客户不存在")
	}
	if authCode.Code != authCode.Code {
		return errors.New("授权码错误-授权码错误")
	}
	if authCode.IsUsed == 1 {
		return errors.New("授权码无效-已经使用过")
	}
	if authCode.ExpireAt < time.Now().Unix() {
		return errors.New("授权码无效-过期")
	}
	return nil
}
func CreateToken(accessToken *AccessToken) error {
	if err := DB.Create(&accessToken).Error; err != nil {
		return err
	}
	return nil
}

// CheckSecret 检查 ClientSecret,redirect_uri 和 ClientId 是否一致
func CheckSecret(id, secret, uri string) error {
	var client Client
	if err := DB.Where("client_id = ?", id).First(&client).Error; err != nil {
		log.Printf("#ERR#model-CheckSecret:客户端检索错误，%v", err)
		return err
	}
	if client.ClientSecret != secret {
		log.Printf("model-CheckSecret:秘钥错误DB:%v,Client:%v", client.ClientSecret, secret)
		return errors.New("client id 与 secret 不匹配")
	}
	if client.RedirectUri != uri {
		log.Printf("model-CheckSecret:uri DB:%v,Client:%v", client.RedirectUri, uri)
		return errors.New("uri 不匹配")
	}
	//var authCode AuthorizationCode
	//if err := DB.Where(&AuthorizationCode{ClientId: id}).First(&authCode).Error; err != nil {
	//	return errors.New("client id 未授权")
	//}
	//if authCode.RedirectUri != uri {
	//	return errors.New("RedirectUri 错误")
	//}
	return nil
}

// GetScope 查询 id 对应 client 的 scope
func GetScope(id string) (string, error) {
	var client Client
	if err := DB.Where("client_id = ?", id).First(&client).Error; err != nil {
		return "", err
	}
	return client.Scope, nil
}

// CheckRefresh 检查 refresh 是否存在、未使用过、未过期
func CheckRefresh(accessTokenParam *AccessToken, secret string) error {
	var client Client
	var accessTokenQuery AccessToken
	if err := DB.Where("client_id = ?", accessTokenParam.ClientId).First(&client).Error; err != nil {
		return errors.New("client_id 未找到")
	}
	if client.ClientSecret != secret {
		return errors.New("secret 错误")
	}
	if err := DB.Where("client_id = ?", accessTokenParam.ClientId).First(&accessTokenQuery).Error; err != nil {
		return errors.New("accessTokenParam token 未找到")
	}
	if accessTokenQuery.ClientId != accessTokenParam.ClientId {
		return errors.New("client_id 与 token 不匹配")
	}
	if accessTokenQuery.RefreshToken != accessTokenParam.RefreshToken {
		return errors.New("refresh token 错误")
	}
	if accessTokenQuery.RefreshExpireAt < time.Now().Unix() {
		return errors.New("refresh token 已经过期")
	}
	// 将参数传递回 controller
	accessTokenParam.UserId = accessTokenQuery.UserId
	return nil
}

// UpdateToken 更新 access 和 有效期(根据用户的id和client_id)
func UpdateToken(tokenStruct *AccessToken) error {
	if err := DB.
		Model(&AccessToken{}).
		Where("user_id = ? AND client_id = ?", tokenStruct.UserId, tokenStruct.ClientId).
		Select("access_token", "access_expire_at").
		Updates(*tokenStruct).
		Error; err != nil {
		return err
	}
	return nil
}

// UpdateClient 更新 client 的过期时间、code
//func UpdateClient(clientId, scope string) error {
//	// 更新 client 的过期时间
//	return nil
//}
