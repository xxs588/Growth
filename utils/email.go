package utils

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// gomail 会自动处理 SSL/TLS 连接
	return d.DialAndSend(m)
}

// GenerateCode 生成 6 位随机验证码
func GenerateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	code := n.Int64() + 100000 // 范围 100000 - 999999
	return strconv.Itoa(int(code))
}
