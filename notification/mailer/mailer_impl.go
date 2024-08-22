package mailer

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gateway struct {
	BaseURL  string
	Host     string
	Port     string
	Username string
	Password string
}

func NewMailerHandler() (SmtpClient, error) {
	config := &cfg.EmailConfig{}
	err := cfg.InitEnv(config)
	if err != nil {
		return nil, err
	}
	g := &gateway{
		Host:     config.EmailHost,
		Port:     config.EmailPort,
		Username: config.EmailUserName,
		Password: config.EmailPassword,
	}
	return g, err
}

func (g *gateway) SendEmailWithFilePaths(ctx context.Context, mailWithoutAttachments MailWithoutAttachments, filePaths []string) (data interface{}, err error) {
	start := time.Now()
	attachments := make([]Attachments, 0)

	for _, filePath := range filePaths {
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		fileName := filepath.Base(filePath)
		contentType := mime.TypeByExtension(filepath.Ext(fileName))

		attachment := Attachments{
			FileName:    fileName,
			Content:     fileContent,
			Encoding:    "base64", // Assuming base64 encoding
			ContentType: contentType,
		}

		attachments = append(attachments, attachment)
	}

	mail := Mail{
		To:          mailWithoutAttachments.To,
		Subject:     mailWithoutAttachments.Subject,
		Message:     mailWithoutAttachments.Message,
		Attachments: attachments,
	}

	return g.SendEmail(ctx, mail, start)
}

func (g *gateway) SendEmail(ctx context.Context, mail Mail, timer time.Time) (data interface{}, err error) {

	if timer == (time.Time{}) {
		timer = time.Now()
	}

	from := g.Username
	password := g.Password
	smtpHost := g.Host
	smtpPort := g.Port

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(mail.To, ",")
	headers["Subject"] = mail.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `multipart/mixed; boundary="MULTIPART_BOUNDARY"`

	header := ""
	for k, v := range headers {
		header += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	bodyHeader := "--MULTIPART_BOUNDARY\r\n" +
		`Content-Type: text/html; charset="UTF-8"` + "\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n" +
		"\r\n" +
		mail.Message +
		"\r\n"

	newAttachments := ""
	for _, attachment := range mail.Attachments {
		newAttachments += "--MULTIPART_BOUNDARY\r\n" +
			`Content-Type: application/octet-stream` + "\r\n" +
			`Content-Transfer-Encoding: base64` + "\r\n" +
			`Content-Disposition: attachment; filename="` + attachment.FileName + `"` + "\r\n" +
			"\r\n" +
			base64.StdEncoding.EncodeToString(attachment.Content) +
			"\r\n"
	}

	newMessage := []byte(header + "\r\n" + bodyHeader + newAttachments + "--MULTIPART_BOUNDARY--")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, mail.To, newMessage)

	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	log.Printf("sendBell took %v", time.Since(timer))

	fmt.Println("Email Sent Successfully!")
	return "OKAY", nil
}

func (g *gateway) NewSmtpClient() SmtpClient {
	return &gateway{
		BaseURL:  g.BaseURL,
		Host:     g.Host,
		Port:     g.Port,
		Username: g.Username,
		Password: g.Password,
	}
}
