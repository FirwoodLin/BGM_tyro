## 注册 - POST

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

参数位置：Body

参数格式：form-data，除 ID 外均为可选参数

```json
id:
name:			
nickname:			
password:	
```

