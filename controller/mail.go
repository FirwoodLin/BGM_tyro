package controller

import (
	"crypto/tls"
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// SendEmail 向用户发送验证码
func SendEmail(user model.User, veriSecret string) error {
	// 定义常量
	message := `
	 <p> 亲爱的BGM用户 %s,</p>
	 
		 <p style="text-indent:2em">你正在使用 BGM 的注册服务，你的激活链接是：</p> 

		 <p style="text-indent:2em"><a href="%s">链接</a> 点击链接完成账号激活 </p>
  
		 <p style="text-indent:2em">请在十分钟内完成激活。如果不是本人操作，请忽略本信息。</P>
	 `
	host := setting.MailSettings.Host
	port := setting.MailSettings.Port
	mailUserName := setting.MailSettings.Username
	mailSecret := setting.MailSettings.Secret
	// 新建实例
	m := gomail.NewMessage()
	// 设置发件人、别名、收件人
	m.SetHeader("From", mailUserName)
	m.SetHeader("From", "BGM Team"+"<"+mailUserName+">")
	m.SetHeader("To", user.Email)
	m.SetHeader("Here's your verify link of BGM") // 邮件主题
	// 设置正文
	link := "http://localhost:8080/verify?token=%s&id=%d"
	veriLink := fmt.Sprintf(link, veriSecret, user.ID)
	m.SetBody("text/html", fmt.Sprintf(message, user.UserName, veriLink))
	// 设置发件邮箱
	d := gomail.NewDialer(
		host,
		port,
		mailUserName,
		mailSecret,
	)
	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 发送
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// VerifyEmail 验证用户激活链接的有效性
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	idInt, _ := strconv.Atoi(c.Query("id"))
	fmt.Printf("in veri:%v %v\n", token, idInt)
	var user = model.User{VeriToken: token, Model: gorm.Model{ID: uint(idInt)}}
	fmt.Printf("before CheckVeri:%v\n", user)
	err := model.CheckVeri(&user)
	fmt.Printf("after CheckVeri err:%v\n", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5555,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    5555,
		"message": "激活成功",
	})
}
