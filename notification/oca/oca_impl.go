package oca

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type OCAHandler struct {
	OCAService OCAClient
}

type gateway struct {
	OCAWABASEURL string
	OCAWAToken   string
}

func NewOCAHandler() OCAClient {
	config := &config.NotifConfig{}
	config.InitEnv()
	g := &gateway{
		OCAWABASEURL: config.OCAWABASEURL,
		OCAWAToken:   config.OCAWAToken,
	}
	return g
}

func (g gateway) SendWhatsapp(ctx context.Context, body OCA) (data interface{}, err error) {

	start := time.Now()
	for _, phoneNumber := range body.PhoneNumber {
		checkPhoneNumber := phoneNumber[:2]
		if checkPhoneNumber == "08" {
			phoneNumber = "62" + phoneNumber[1:]
		} else if checkPhoneNumber == "+6" {
			phoneNumber = phoneNumber[1:]
		} else if checkPhoneNumber != "62" {
			return nil, errors.New("invalid phone number")
		}

		messageData := MessageData{
			PhoneNumber: phoneNumber,
			Message:     body.MessageData,
		}

		templateCode := messageData.Message.Template.TemplateCodeID

		templateCodePattern := `^[a-f0-9]{8}_[a-f0-9]{4}_[a-f0-9]{4}_[a-f0-9]{4}_[a-f0-9]{12}:[a-z0-9]+$`
		templateCodeRegex := regexp.MustCompile(templateCodePattern)

		if templateCode == "" {
			return nil, errors.New("template code is required")
		}
		if !templateCodeRegex.MatchString(templateCode) {
			return nil, errors.New("invalid template code")
		}

		messageDataJSON, err := json.Marshal(messageData)
		if err != nil {
			return nil, err
		}

		url := g.OCAWABASEURL + "/api/v2/push/message"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageDataJSON))
		if err != nil {
			fmt.Println("Error creating request: ", err)
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+g.OCAWAToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("Failed to send notification")
		}

	}

	response := map[string]interface{}{
		"message": "Whatsapp sent successfully",
		"status":  "success",
	}

	log.Printf("sendBell took %v", time.Since(start))

	return response, nil
}
