package oca

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gateway struct {
	OCAWABASEURL string
	OCAWAToken   string
}

func NewOCAHandler() (OCAClient, error) {
	config, err := cfg.InitEnv(cfg.OCA)
	if err != nil {
		log.Fatalf("Failed to initialize environment: %v", err)
		return nil, err
	}
	g := &gateway{
		OCAWABASEURL: config.OCAConfig.OCAWABASEURL,
		OCAWAToken:   config.OCAConfig.OCAWAToken,
	}
	return g, nil
}

func (g gateway) SendWhatsapp(ctx context.Context, body OCA) (data interface{}, err error) {
	start := time.Now()
	var wg sync.WaitGroup
	results := make(chan error, len(body.PhoneNumber))

	for _, phoneNumber := range body.PhoneNumber {
		wg.Add(1)
		go func(phoneNumber string) {
			defer wg.Done()

			checkPhoneNumber := phoneNumber[:2]
			if checkPhoneNumber == "08" {
				phoneNumber = "62" + phoneNumber[1:]
			} else if checkPhoneNumber == "+6" {
				phoneNumber = phoneNumber[1:]
			} else if checkPhoneNumber != "62" {
				results <- errors.New("invalid phone number")
				return
			}

			messageData := MessageData{
				PhoneNumber: phoneNumber,
				Message:     body.MessageData,
			}

			templateCode := messageData.Message.Template.TemplateCodeID
			templateCodePattern := `^[a-f0-9]{8}_[a-f0-9]{4}_[a-f0-9]{4}_[a-f0-9]{4}_[a-f0-9]{12}:[a-z0-9]+$`
			templateCodeRegex := regexp.MustCompile(templateCodePattern)

			if templateCode == "" {
				results <- errors.New("template code is required")
				return
			}
			if !templateCodeRegex.MatchString(templateCode) {
				results <- errors.New("invalid template code")
				return
			}

			messageDataJSON, err := json.Marshal(messageData)
			if err != nil {
				results <- err
				return
			}

			url := g.OCAWABASEURL + "/api/v2/push/message"
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageDataJSON))
			if err != nil {
				results <- err
				return
			}

			req.Header.Set("Authorization", "Bearer "+g.OCAWAToken)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- errors.New("failed to send notification")
				return
			}

			results <- nil
		}(phoneNumber)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res != nil {
			return nil, res
		}
	}

	response := map[string]interface{}{
		"message": "Whatsapp sent successfully",
		"status":  "success",
	}

	log.Printf("sendNotif took %v", time.Since(start))

	return response, nil
}
