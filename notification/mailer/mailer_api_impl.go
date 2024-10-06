package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gatewayApi struct {
	FabdBaseUrl string
	ApiKey      string
}

type ApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func NewMailerApiHandler() (SmtpClient, error) {
	config, err := cfg.InitEnv(cfg.API)
	if err != nil {
		return nil, err
	}
	g := &gatewayApi{
		FabdBaseUrl: config.ApiConfig.FabdBaseUrl,
		ApiKey:      config.ApiConfig.ApiKey,
	}
	return g, err
}

func (g *gatewayApi) SendEmailWithFilePaths(ctx context.Context, mailWithoutAttachments MailWithoutAttachments, filePaths []string) (data interface{}, err error) {
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

func (g *gatewayApi) SendEmail(ctx context.Context, payload Mail) (data interface{}, err error) {
	url := g.FabdBaseUrl + "/v4/webhooks/email-notifications"
	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	for _, recipient := range payload.To {
		_ = writer.WriteField("to", recipient)
	}
	for _, cc := range payload.CC {
		_ = writer.WriteField("cc", cc)
	}
	for _, bcc := range payload.BCC {
		_ = writer.WriteField("bcc", bcc)
	}
	_ = writer.WriteField("subject", payload.Subject)
	_ = writer.WriteField("template_code", payload.TemplateCode)
	dataJson, _ := json.Marshal(payload.Data)
	_ = writer.WriteField("data", string(dataJson))

	if len(payload.Attachments) > 0 {
		for _, attachment := range payload.Attachments {
			file, err := os.Open(attachment.Path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			part, err := writer.CreateFormFile("attachments", attachment.FileName)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, err
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, form)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", g.ApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	log.Println("Response from external endpoint:", resp.Status)
	return apiResponse, nil
}
