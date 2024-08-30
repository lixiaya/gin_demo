package util

import (
	"gin_demo/global"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendEmail(to, body string) error {
	from := "trhxnlove@yeah.net"
	password := "FNOAIURDUWQAIYAM"
	e := email.NewEmail()
	e.From = from         //发送方
	e.To = []string{to}   //接收方
	e.Subject = "测试验证码"   //邮件主题
	e.Text = []byte(body) //发送内容
	err := e.Send("smtp.yeah.net:25", smtp.PlainAuth("", from, password, "smtp.yeah.net"))
	if err != nil {
		global.Logger.Error("邮件发送失败：", err)
		return err
	}
	return nil
}
