package controller

import (
	"crypto/tls"
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"gopkg.in/gomail.v2"
)

// SendEmail 向用户发送验证码
func SendEmail(usermail, username, vericode string) error {
	// 定义常量
	message := `
	 <p> 亲爱的BGM用户 %s,</p>
	 
		 <p style="text-indent:2em">你正在使用 BGM 的注册服务，验证码是：</p> 

		 <p style="text-indent:2em">%s</p>
  
		 <p style="text-indent:2em">如果不是本人操作，请忽略本信息。</P>
	 `
	host := setting.MailSettings.Host
	port := setting.MailSettings.Port
	mailUserName := setting.MailSettings.Username
	mailSecret := setting.MailSettings.Secret
	// 新建实例
	m := gomail.NewMessage()
	// 设置发件人、别名、收件人
	m.SetHeader("From", mailUserName)
	m.SetHeader("From", "BGM 运营团队"+"<"+mailUserName+">")
	m.SetHeader("To", usermail)
	m.SetHeader("这是您的BGM验证码") // 邮件主题
	// 设置正文
	m.SetBody("text/html", fmt.Sprintf(message, username, vericode))
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
