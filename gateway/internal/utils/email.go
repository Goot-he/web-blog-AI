package utils

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"gateway/internal/config"
	"math/big"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/jordan-wright/email"
)

const (
	emailCodeKeyPrefix = "auth:"
	emailCodeTTL       = 2 * time.Minute
)

type Email struct {
	cfg *config.EmailConf
}

func NewEmail(cfg *config.EmailConf) *Email {
	return &Email{cfg: cfg}
}

func (e *Email) SendCode(email, purpose string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	purpose = strings.TrimSpace(strings.ToLower(purpose))
	if len(email) == 0 {
		return fmt.Errorf("email is empty")
	}
	if strings.TrimSpace(e.cfg.Host) == "" || strings.TrimSpace(e.cfg.Port) == "" || strings.TrimSpace(e.cfg.From) == "" {
		return fmt.Errorf("email address is empty")
	}

	if purpose != "register" && purpose != "login" && purpose != "forgot_password" {
		return fmt.Errorf("email purpose is invalid")
	}

	// 检查当前缓存里面是否已经存在验证码，如果存在返回验证码已发送 提醒用户查看邮箱
	return fmt.Errorf("验证码已发送 请注意查看邮箱 不要重复点击发送")

	// 如果不存在生成验证码 并存入缓存
	code, err := GenerateCode()
	if err != nil {
		return err
	}

	// 发送邮箱 并接收发送状态
	err = e.send(email, code)
	if err != nil {
		// 如果失败返回错误 并删除cache里存储的验证码

		return err
	}

	// 发送成功 返回nil
	return nil
}

// 发送邮件
func (e *Email) send(to, code string) error {
	addr := net.JoinHostPort(e.cfg.Host, e.cfg.Port)
	auth := smtp.PlainAuth("", e.cfg.From, e.cfg.Secret, e.cfg.Host)

	//配置email对象
	em := email.NewEmail()

	// 格式化发件人地址为 "昵称 <邮箱>"
	em.From = fmt.Sprintf("%s <%s>", e.cfg.Nickname, e.cfg.From)

	// 设置收件人、主题和邮件内容
	em.To = []string{to}
	em.Subject = "Renai 邮箱验证码"
	em.HTML = []byte(`您好，<br/><br/>
您正在进行 Renai 的邮箱验证。验证码为：<br/><br/>
<strong style="font-size:24px;letter-spacing:6px;color:#2563eb">` + code + `</strong><br/><br/>
验证码 5 分钟内有效，请勿向任何人泄露。若非本人操作，请忽略此邮件。<br/><br/>
Renai 团队`)

	var err error
	if e.cfg.IsSSL {
		err = em.SendWithTLS(addr, auth, &tls.Config{ServerName: e.cfg.Host})
	} else {
		err = em.Send(addr, auth)
	}
	return err
}

// GenerateCode 生成验证码
func GenerateCode() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", number.Int64()), nil
}
