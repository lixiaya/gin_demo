package util

import (
	"fmt"
	"gin_demo/global"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"net/smtp"
)

func SendEmail(to, body string) error {
	from := viper.GetString("email.address")
	password := viper.GetString("email.password")
	e := email.NewEmail()
	e.From = from                                //发送方
	e.To = []string{to}                          //接收方
	e.Subject = viper.GetString("email.subject") //邮件主题
	e.Text = []byte(body)                        //发送内容
	err := e.Send(fmt.Sprintf("%s:%d", viper.GetString("email.addrhost"), 25), smtp.PlainAuth("", from, password, viper.GetString("email.addrhost")))
	if err != nil {
		global.Logger.Error("邮件发送失败：", err)
		return err
	}
	return nil
}
