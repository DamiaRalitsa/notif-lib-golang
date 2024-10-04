package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
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

	mail := Mail{
		To:          mailWithoutAttachments.To,
		Subject:     mailWithoutAttachments.Subject,
		Message:     mailWithoutAttachments.Message,
		Attachments: attachments,
	}

	g.SendEmail(ctx, mail)

	return "OKAY", nil
}

func (g *gatewayApi) SendEmail(ctx context.Context, payload Mail) (data interface{}, err error) {
	url := g.FabdBaseUrl + "/v4/notification-service/notifications/mailer"
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	log.Printf("Payload: %s", string(jsonData))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", g.ApiKey)
	req.Header.Set("Content-Type", "application/json")

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
