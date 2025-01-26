package email

import (
	"GoForPT/pkg/cfg"
	"fmt"
	"gopkg.in/gomail.v2"
)

// Mail 结构体用于存储邮件信息
type Mail struct {
	From    string // 发件人
	To      string // 收件人
	Subject string // 主题
	Body    string // 正文
}

// SMTPConfig 结构体用于存储 SMTP 配置
type SMTPConfig struct {
	Host     string // SMTP 服务器地址
	Port     int    // SMTP 服务器端口
	Username string // SMTP 用户名
	Password string // SMTP 密码
}

// SendMail 发送邮件
func SendEmail(to, subject, body string) error {

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Cfg.SMTP.Username) // 发件人
	m.SetHeader("To", to)                      // 收件人
	m.SetHeader("Subject", subject)            // 邮件主题
	m.SetBody("text/plain", body)              // 邮件正文

	// 创建 SMTP 客户端
	d := gomail.NewDialer(cfg.Cfg.SMTP.Host, cfg.Cfg.SMTP.Port, cfg.Cfg.SMTP.Username, cfg.Cfg.SMTP.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Email sent successfully")
	return nil
}
