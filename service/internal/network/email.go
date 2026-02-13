package network

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/smtp"

	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/internal/configuration"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type Email struct {
	Configuration configuration.ConfigurationService
}

func (m Email) dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func (m Email) sendMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	//create smtp client
	c, err := m.dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func (m Email) SendEmailTLS(ctx context.Context, oid dao.PrimaryKey, c serviceargument.EmailContent) error {
	stmpInfo := m.Configuration.GetEmailSTMP(ctx, oid)
	addr := stmpInfo.EmailSTMPHost + ":" + stmpInfo.EmailSTMPPort
	auth := smtp.PlainAuth("", stmpInfo.EmailSTMPFrom, stmpInfo.EmailSTMPPassword, stmpInfo.EmailSTMPHost)

	for _, email := range c.ToEmails {
		message := []byte("To: " + email + "\r\n" +
			"Subject: " + c.Subject + "\r\n" +
			"From: " + c.SenderName + " <" + stmpInfo.EmailSTMPFrom + ">\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=utf-8\r\n" + // 将内容类型设置为 HTML
			"\r\n" + c.Content)

		err := m.sendMailUsingTLS(addr, auth, stmpInfo.EmailSTMPFrom, c.ToEmails, message)
		if err != nil {
			return err
		}
	}
	return nil
}
