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
	config, err := cfg.InitEnv(cfg.Email)
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
	attachments := make([]Attachments, len(filePaths))

	type result struct {
		index      int
		attachment Attachments
		err        error
	}

	results := make(chan result, len(filePaths))
	var wg sync.WaitGroup

	for i, filePath := range filePaths {
		wg.Add(1)
		go func(i int, filePath string) {
			defer wg.Done()
			fileContent, err := ioutil.ReadFile(filePath)
			if err != nil {
				results <- result{i, Attachments{}, err}
				return
			}

			fileName := filepath.Base(filePath)
			contentType := mime.TypeByExtension(filepath.Ext(fileName))

			attachment := Attachments{
				FileName:    fileName,
				Content:     fileContent,
				Encoding:    "base64", // Assuming base64 encoding
				ContentType: contentType,
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

	wg.Wait()

	mail := Mail{
		To:          mailWithoutAttachments.To,
		Subject:     mailWithoutAttachments.Subject,
		Message:     mailWithoutAttachments.Message,
		Attachments: attachments,
	}

	success, err := g.SendEmail(ctx, mail)
	if err != nil {
		log.Printf("Failed to send email with filePaths: %v\n", err)
		return nil, err
	}

	log.Printf("sendNotif took %v", time.Since(start))
	log.Println("Email Sent Successfully!")
	return success, err
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
		mail.Message +
		"\r\n"

	type result struct {
		encodedAttachment string
		err               error
	}

	results := make(chan result, len(mail.Attachments))
	var wg sync.WaitGroup

	for _, attachment := range mail.Attachments {
		wg.Add(1)
		go func(attachment Attachments) {
			defer wg.Done()
			encodedAttachment := "--MULTIPART_BOUNDARY\r\n" +
				`Content-Type: application/octet-stream` + "\r\n" +
				`Content-Transfer-Encoding: base64` + "\r\n" +
				`Content-Disposition: attachment; filename="` + attachment.FileName + `"` + "\r\n" +
				"\r\n" +
				base64.StdEncoding.EncodeToString(attachment.Content) +
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
