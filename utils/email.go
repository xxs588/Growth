package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// 组装认证信息
	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	// 在正式的SMTP协议规范中，使用 \r\n 是必须的
	// 连接 SMTP 服务器发送邮件
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	return err
}

// GenerateCode 生成 6 位随机验证码
func GenerateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	code := n.Int64() + 100000 // 范围 100000 - 999999
	return strconv.Itoa(int(code))
}
