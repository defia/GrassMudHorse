package main

import (
	"net/smtp"
	"strings"
	"time"
)

func sendMail(user, password, host, to, subject, body, mailtype string) error {

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")

	err := smtp.SendMail(host, auth, user, send_to, msg)

	return err
}

func SendMail(user, password, host, to string, typ bool) error {
	var subject string
	if typ {
		subject = "所有服务器丢包率超过阕值！"
	} else {
		subject = "有服务器恢复正常"
	}
	go func() {
		for i := 0; i < 5; i++ {
			err := sendMail(user, password, host, to, subject, "建议您更改设置", "html")
			if err != nil {
				time.Sleep(time.Second * 3)
				continue

			}
			break
		}
	}()
	return nil
}
