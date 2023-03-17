package model

import (
	"errors"
	"strings"
)

// CheckClient 检验 clientId 和 scope 范围是否匹配
func CheckClient(clientId, scope string) error {
	var client AuthorizationCode
	// 检查 client 是否注册
	if err := DB.Where("client_id = ?", clientId).Find(&client).Error; err != nil {
		return errors.New("client not found")
	}
	// 检查 scope 是否已经授权
	var authedScopeMap map[string]int
	for _, v := range strings.Split(client.Scope, ",") {
		authedScopeMap[v] = 1
	}
	for _, v := range strings.Split(scope, ",") {
		if authedScopeMap[v] != 1 {
			return errors.New("unauthed scope")
		}
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

// UpdateClient 更新 client 的过期时间、code
func UpdateClient(clientId, scope string) error {
	// TODO:更新 client 的过期时间
	return nil
}
