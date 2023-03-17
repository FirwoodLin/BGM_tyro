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

-   [ ] 阶段一  - 邮箱认证功能
-   [ ] 阶段三  - refresh
-   [ ] 阶段三 - 合规

## Debug 记录