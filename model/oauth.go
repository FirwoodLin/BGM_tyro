package model

import (
	"errors"
	"log"
	"strings"
	"time"
)

// CheckScope 检验 clientId 和 scope 范围是否匹配
func CheckScope(clientId, scope, redirectUri string) error {
	var client AuthorizationCode
	// 检查 client 是否注册
	if err := DB.Where("client_id = ?", clientId).Find(&client).Error; err != nil {
		log.Printf("model-CheckScope:client not found\n")
		return errors.New("client not found")
	}
	// 检查 scope 是否已经授权
	var authedScopeMap map[string]int
	authedScopeMap = make(map[string]int)
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
	return nil
}

// UpdateAuthCode 发布 authorization code 之后持久化到数据库
func UpdateAuthCode(authCode *AuthorizationCode) error {
	if err := DB.Create(&authCode).Error; err != nil {
		return err
	}
	return nil
}

// CheckCode 检验 code 是否使用过、是否过期；有效就返回 scope 和nil
func CheckCode(code string) error {
	var authCode AuthorizationCode
	if err := DB.Where(&AuthorizationCode{Code: code}).First(&authCode); err != nil {
		return errors.New("授权码无效-不存在")
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
	if err := DB.Where(&Client{ClientId: id}).First(&client).Error; err != nil {
		return err
	}
	if client.ClientSecret != secret {
		log.Printf("model-CheckSecret:秘钥不匹配DB:%v,Client:%v", client.ClientSecret, secret)
		return errors.New("client id 与 secret不匹配")
	}
	var authCode AuthorizationCode
	if err := DB.Where(&AuthorizationCode{ClientId: id}).First(&authCode).Error; err != nil {
		return errors.New("client id 未授权")
	}
	if authCode.RedirectUri != uri {
		return errors.New("RedirectUri 错误")
	}
	return nil
}

// GetScope 查询 code 对应的 scope
func GetScope(id string) (string, error) {
	var authCode AuthorizationCode
	if err := DB.Where(&AuthorizationCode{ClientId: id}).First(&authCode).Error; err != nil {
		return "", err
	}
	return authCode.Scope, nil
}

//func UpdateVeri() {
//
//}

// UpdateClient 更新 client 的过期时间、code
//func UpdateClient(clientId, scope string) error {
//	// 更新 client 的过期时间
//	return nil
//}
