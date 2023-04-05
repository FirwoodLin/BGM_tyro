# HTTP 状态码的使用

| code | message                                   |
| ---- | ----------------------------------------- |
| 200  | OK：成功响应                              |
| 302  | Found：OAuth 中用于重定向                 |
| 400  | BadRequest：客户端的请求出错              |
| 500  | ServerInternalError：服务器处理过程中出错 |

# 自定义返回值



# 用户增改查

## 注册 - POST

路由地址：`/signup`

参数位置：Body

参数格式：form-data

```json
name:			// 用户名
email:			// 邮箱
nickname:		// 昵称
password:		// 密码：8 位以上
description:	// 简介：256 字以内
avatar:			// url 链接
```

返回值：JSON

```json
{
    "code": 200,
    "message": "注册成功",
    "token": "示例 token"
}
```

## 登录 - POST

路由地址：`/signin`

参数位置：Body

参数格式：form-data

```json
name:			// 用户名
email:			// 邮箱 与用户名二选一
password:		// 密码：8 位以上
```

返回值：JSON

```json
{
    "code": 200,
    "message": "登录成功"
    "token": "示例 token"
}
```

## 修改

Hader 中携带 Auth 信息

路由地址：`/user`

参数位置：Body

参数格式：form-data，除 ID 外均为可选参数

```json
id:
name:			
nickname:			
password:	
```

# OAuth2.0

## 请求 Authoritarian Code - Get

路由地址：`/oauth/authorization`

参数位置：Header

参数格式：Query

| key           | note                 |
| ------------- | -------------------- |
| response_type | 常量`code`           |
| client_id     | 客户端 id            |
| redirect_uri  | 重定向 uri           |
| scope         | 以 , 分隔            |
| state         | 客户端生成的随机序列 |

响应：

成功 - 重定向到如下链接

```json
redirectUri+"?code="+authCode+"&state="+state
```

## 使用 code 兑换 access token 和 refresh token  - POST

路由地址：`/oauth/accesstoken`

参数位置和类型：

Header Auth - 字符串（Basic认证）

```
Basic (urlencode 编码后的客户端id:密码)
```

Body - x-www-form-urlencoded

| 参数名       | 内容                          |
| ------------ | ----------------------------- |
| grant_type   | `authorization_code`          |
| code         | 先前拿到的 Authorization Code |
| redirect_uri | 注册的重定向 uri              |

响应：JSON

| key             | val         |
| --------------- | ----------- |
| "access_token"  |             |
| "token_type"    | "bearer"    |
| "expires_in"    | 2h          |
| "refresh_token" |             |
| "scope"         | 请求的scope |

## OIDC-请求 ID-Token
