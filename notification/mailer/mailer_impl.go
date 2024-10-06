package mailer

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"path/filepath"
	"strings"
	"sync"
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
	config, err := cfg.InitEnv(cfg.EMAIL)
	if err != nil {
		return nil, err
	}
	g := &gateway{
		Host:     config.EmailConfig.EmailHost,
		Port:     config.EmailConfig.EmailPort,
		Username: config.EmailConfig.EmailUserName,
		Password: config.EmailConfig.EmailPassword,
	}
	return g, err
}

func (g *gateway) SendEmailWithFilePaths(ctx context.Context, mailWithoutAttachments MailWithoutAttachments, filePaths []string) (data interface{}, err error) {
	start := time.Now()
	defer func() {
		log.Printf("readFiles %v", time.Since(start))
	}()

	attachments := make([]Attachment, len(filePaths))

	type result struct {
		index      int
		attachment Attachment
		err        error
	}

	results := make(chan result, len(filePaths))
	var wg sync.WaitGroup

	for i, filePath := range filePaths {
		wg.Add(1)
		go func(i int, filePath string) {
			defer wg.Done()
			fileName := filepath.Base(filePath)
			attachment := Attachment{
				FileName: fileName,
				Path:     filePath,
			}
			results <- result{i, attachment, nil}
		}(i, filePath)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		attachments[res.index] = res.attachment
	}

	mail := Mail{
		To:           mailWithoutAttachments.To,
		Subject:      mailWithoutAttachments.Subject,
		TemplateCode: mailWithoutAttachments.Message,
		Data:         map[string]interface{}{"text": mailWithoutAttachments.Text},
	}

	return g.SendEmail(ctx, mail)
}

func (g *gateway) SendEmail(ctx context.Context, mail Mail) (data interface{}, err error) {
	start := time.Now()

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
		mail.TemplateCode +
		"\r\n"

	type result struct {
		encodedAttachment string
		err               error
	}

	results := make(chan result, len(mail.Attachments))
	var wg sync.WaitGroup

	for _, attachment := range mail.Attachments {
		wg.Add(1)
		go func(attachment Attachment) {
			defer wg.Done()
			fileContent, err := ioutil.ReadFile(attachment.Path)
			if err != nil {
				results <- result{"", err}
				return
			}
			encodedAttachment := "--MULTIPART_BOUNDARY\r\n" +
				`Content-Type: application/octet-stream` + "\r\n" +
				`Content-Transfer-Encoding: base64` + "\r\n" +
				`Content-Disposition: attachment; filename="` + attachment.FileName + `"` + "\r\n" +
				"\r\n" +
				base64.StdEncoding.EncodeToString(fileContent) +
				"\r\n"
			results <- result{encodedAttachment, nil}
		}(attachment)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	newAttachments := ""
	for res := range results {
		if res.err != nil {
			return "Failed", res.err
		}
		newAttachments += res.encodedAttachment
	}

	wg.Wait()

	newMessage := []byte(header + "\r\n" + bodyHeader + newAttachments + "--MULTIPART_BOUNDARY--")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, mail.To, newMessage)

	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	log.Printf("sendNotif took %v", time.Since(start))
	log.Println("Email Sent Successfully!")
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
