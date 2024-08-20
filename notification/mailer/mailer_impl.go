package mailer

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type MailerHandler struct {
	notificationService SmtpClient
}

type gateway struct {
	BaseURL  string
	Host     string
	Port     string
	Username string
	Password string
	// HttpClient *helpers.ToolsAPI
}

// NewMailerHandler creates a new MailerHandler instance.
func NewMailerHandler() SmtpClient {
	config := &config.NotifConfig{}
	config.InitEnv()
	g := &gateway{
		BaseURL:  config.FabdCoreUrl,
		Host:     config.EmailHost,
		Port:     config.EmailPort,
		Username: config.EmailUserName,
		Password: config.EmailPassword,
		// HttpClient: helpers.NewToolsAPI(payloadEmail.FabdCoreUrl),
	}
	return g
}

func (g *gateway) SendEmailWithFilePaths(ctx context.Context, mailWithoutAttachments MailWithoutAttachments, filePaths []string) (data interface{}, err error) {
	attachments := make([]Attachments, 0)

	for _, filePath := range filePaths {
		// cwd, err := os.Getwd()
		// if err != nil {
		// 	log.Fatalf("unable to get current directory: %v", err)
		// }

		// templatePath := filepath.Join(cwd, "assets", "templates", filePath)
		// log.Printf("templatePath: %v", templatePath)

		// Read the file content
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		// Determine the file's content type
		fileName := filepath.Base(filePath)
		contentType := mime.TypeByExtension(filepath.Ext(fileName))

		// Map the file content and metadata into the model.Attachments struct
		attachment := Attachments{
			FileName:    fileName,
			Content:     fileContent,
			Encoding:    "base64", // Assuming base64 encoding
			ContentType: contentType,
		}

		attachments = append(attachments, attachment)
	}

	// Create the Mail struct
	mail := Mail{
		To:          mailWithoutAttachments.To,
		Subject:     mailWithoutAttachments.Subject,
		Message:     mailWithoutAttachments.Message,
		Attachments: attachments,
	}

	// Call the existing SendEmail method with the Mail struct
	return g.SendEmail(ctx, mail)
}

func (g *gateway) SendEmail(ctx context.Context, mail Mail) (data interface{}, err error) {
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

	fmt.Println("Email Sent Successfully!")
	return "OKAY", nil
}

func (g *gateway) NewSmtpClient() SmtpClient {
	// baseUrl := os.Getenv("FABD_API_CORE_URL")
	// host := os.Getenv("EMAIL_HOST")
	// port := os.Getenv("EMAIL_PORT")
	// username := os.Getenv("EMAIL_USERNAME")
	// password := os.Getenv("EMAIL_PASSWORD")
	// httpClient := &helpers.ToolsAPI{}
	return &gateway{
		BaseURL: g.BaseURL,
		Host:    g.Host,
		Port:    g.Port,
		// HttpClient: httpClient,
		Username: g.Username,
		Password: g.Password,
	}
}
