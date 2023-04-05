# 03-11 Day0

- [x] 选题：OIDC

- [x] 配置 Git

- [x] 连接数据库


## 学习记录

-   学习了 JSON 配置文件的读取，避免了 push 敏感信息上网。
-   commit 的 Angular 规范

# 03-12 Day1

-   [ ] 阶段一基础任务
-   [ ] 阶段一进阶任务

## 学习记录

-   使用 APIFox 进行本地调试
-   密码加密的两种方式：
    -   CSPRNG 生成盐值，SHA256 加盐到密码前段
    -   bcrypt 加密，无需单独存储盐值（采用此方案）
-   正则表达式在字段校验中的使用 

## 项目进度

-   设计用户表
-   完成注册功能

# 03-13 Day2

## 学习记录

-   JWT 的原理，JWT 配合 gin 的使用

## 项目进度

-   完成登录，增加 token

# 03-14 Day3

## 学习记录

-   使用 viper 读取配置文件

## 项目进度

# 0315 - Day4

## 学习记录

-   初步学习 Go 工程项目架构
-   OAuth 2.0 截至获取授权码（阮一峰）
-   OAuth 2.0 实现细节：[授权码 - 原理和方法](https://www.cnblogs.com/blowing00/p/4524412.html)，[state 作用](https://www.cnblogs.com/blowing00/p/14872312.html)，[CSRF攻击](https://www.cnblogs.com/hyddd/archive/2009/04/09/1432744.html)。[Client 的注册和登录](https://blog.yorkxin.org/posts/oauth2-2-cilent-registration/)，[授权码模式 - 细节](https://blog.yorkxin.org/posts/oauth2-4-1-auth-code-grant-flow/)，[JWT、JWE、JWS 、JWK](https://www.51cto.com/article/630971.html)

## 项目进度

-   完成登录、修改功能
-   将 controller, router, db 操作分离

# 0316 - Day5

## 学习记录

-   [OAuth2.0报文](https://learnku.com/articles/20082)
-   对称加密-AES，填充算法。由于 JWT 签名结果较长，考虑 code 仅在服务器使用，~~所以采用 AES 加密 ClientId~~
-   使用`crypto/rand`生成随机字符串作为授权码
-   [Github OAuth2文档](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps)参考设计 API

## 项目进度

-   完成 OAuth2.0 的授权码部分

## Debug 记录

-   DB 中创建了 retUser，形参定义为 u；signin 中创建了 retUser。u 未使用，导致 signin 中的 retUser 始终为空。

# 0317 - Day6

## 学习记录

-   `validator`的使用。`min/max`限制字符串的长度；`email/uri`限制类别

-   `gomail`的使用；

-   [access token的格式](https://stackoverflow.com/questions/50031993/what-characters-are-allowed-in-an-oauth2-access-token)：b64编码，适合在 header 中传输

-   嵌套结构体的字段初始化

-   rand 生成的随机数据，使用base64URL编码为字符串

## 项目进度

-   [x] 阶段一  - 邮箱认证功能
-   [x] 阶段二 - auth code 的生成与返回

## Debug 记录

# 0318 - Day7

验收日

## 学习记录

-   API 中 HTTP 状态码和自定义状态码的设计
-   重读[授权码模式-全流程解析](https://www.cnblogs.com/blowing00/p/4524412.html)，[授权码-实现细节](https://blog.yorkxin.org/posts/oauth2-4-1-auth-code-grant-flow/)

## 开发记录

-   AccessToken 和 Refresh Token 的颁发

# 0325 - Day n+1

整理一下 gORM 的查找

## Debug 记录

-   ```go
    if realUser.VeriTokenExpireAt < time.Now().Unix() {
    		DB.Delete(&user)
    		return errors.New("token过期，请重新注册")
    	}
    ```

    应为小于；本为大于

-   检测激活链接的有效性后，没有设置数据库中的“激活状态”一列

## 进度

-   邮箱激活 fix
-   发现 JWT token 应该用 id 作为统一的身份标识，而不是用户名/邮箱——这会造成验证困难

# 0327 - Day n+2

## 开发进度

-   处理使用 map 前先 make 的问题

-   Oauth2.0 颁发 code


## 学习记录

## TODO

-   重定向至登录页（重定向后默认是GET，但注册为POST）
-   关于重定向
    -   返回到 client 网站（本地监听一个端口）。服务器处理 code 换 access 的请求

-   一些对过程的梳理
    -   ClientId 的校验：去client表查找；
    -   code 的颁发：依据 clientID 和登录产生的 token，存储到授权码这个 table；
    -   access 的颁发：
        -   依据 code 和 clientID 存储到 access token 表
        -   依据 refresh token 和 id 更新到 access token

# 0328 Day n+3

将Code、Token、Client的表重新设计。

重新设计颁发 Access Token 的流程。

# 0329 Dayn+4

使用 Refresh Token 进行换发 Access Token，在此过程中不更新 Refresh Token 及其有效期。

## TODO

-   Authorization Code 被二度利用时，撤销先前的授权
-   Client 的注册
-   第四阶段 - OIDC 服务逻辑梳理。完成 提供ID Token 的功能。
-   了解OIDC和OAuth2.0区别

# 0331 Dayn+6

突然发现可以使用 gin 提供的 JWT 中间件

```
范围参数必须以 openid 值开头，然后包含 profile 值和/或 email 值。
如果存在 profile 范围值，则 ID 令牌可能包含（但不保证）包含用户的默认 profile 声明。
如果存在 email 范围值，则 ID 令牌包含 email 和 email_verified 声明。

一起在 ID-token 中返回
```

-   idtoken JWT 使用 RS256 加密，公钥通过接口暴露；对比：登录token 使用 HS256 进行对称加密

# 0401 Day n+7

-   ~~每个 user 的 client ID 是不同的~~

# 0404

1.25h

-   其实是根据 header 中的 token 确定用户身份
-   jwt 解析时，需要验证 jwt 的 user 和 body 中的 user 相同
-   Oauth 中 根据 header 中的 token 确定用户身份

## 学习的资料

细节在于，“同意授权”的过程中，会将 request 发向服务器，其中会包含用户的信息

>   在 (B) 裡面， Resource Owner 若同意授權，這個**「同意授權」的 request 會往 Authorization Endpoint 發送**，接著會收到 302 的轉址 response ，裡面帶有「前往 Client 的 Redirection Endpoint 的 URL」的轉址 (Location header)，從而產生「向 Redirection URI 發送 GET Request」的操作。

微博对于Oauth2.0过程中产生的错误的定义（部分）[授权机制说明](https://open.weibo.com/wiki/授权机制说明)

| 错误码(error)  | 错误编号(error_code) | 错误描述(error_description)      |
| :------------- | :------------------- | :------------------------------- |
| invalid_client | 21324                | client_id或client_secret参数无效 |
| expired_token  | 21327                | token过期                        |

-   QQ 的文档[QQ互联](https://wiki.connect.qq.com/使用authorization_code获取access_token)：客户端注意需要将url进行URLEncode；refresh_token仅一次有效；

例子：花瓣网点击微博登录，跳转到

```
(已经 url 解码)
https://api.weibo.com/oauth2/authorize
?client_id=2499394483
&response_type=code
&redirect_uri=https://huaban.com/oauth/callback/&display=default
// 没有登录微博时
https://api.weibo.com/oauth2/authorize
?client_id=2499394483
&response_type=code
&redirect_uri=https%3A%2F%2Fhuaban.com%2Foauth%2Fcallback%2F
&display=default###
```

# 0405

回看 auth code 的请求，是要在 header 中带有 bearer token 的信息的

不管是前端是怎么获取的吧，默认会提供 auth 信息。存储 auth code 和 access token 的时候，将 user_id 作为主键

## 计划

重定向链接需要进行 url 编码

## 记录

-   [debug]JWT验证后`c.Set("user_id")`为字符串类型，但表单中是 uint 类型，二者无法比较相等
-   [feat]

# 开发计划

4月2日中午前完成最终 commit

-   Refresh Token 换 Access Token
-   实现 OIDC 服务
-   收藏番剧
-   好友功能
-   绑定 bangumi