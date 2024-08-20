package oca

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/DamiaRalitsa/notif-agent-golang/notification/config"
)

type OCAHandler struct {
	OCAService OCAClient
}

// Gateway represents the structure of the gateway.
type gateway struct {
	OCAWABASEURL string
	OCAWAToken   string
	// HttpClient *helpers.ToolsAPI
}

// NewOCAHandler creates a new OCAHandler instance.
func NewOCAHandler() OCAClient {
	config := &config.NotifConfig{}
	config.InitEnv()
	g := &gateway{
		OCAWABASEURL: config.OCAWABASEURL,
		OCAWAToken:   config.OCAWAToken,
	}
	return g
}

// SendOCA sends a OCA message to multiple phone numbers.
func (g gateway) SendWhatsapp(ctx context.Context, body OCA) (data interface{}, err error) {
	for _, phoneNumber := range body.PhoneNumber {
		// Reformat phone number
		checkPhoneNumber := phoneNumber[:2]
		if checkPhoneNumber == "08" {
			phoneNumber = "62" + phoneNumber[1:]
		} else if checkPhoneNumber == "+6" {
			phoneNumber = phoneNumber[1:]
		} else if checkPhoneNumber != "62" {
			return nil, errors.New("Invalid phone number")
		}

		// Construct message data
		messageData := MessageData{
			PhoneNumber: phoneNumber,
			Message:     body.MessageData,
		}

		templateCode := messageData.Message.Template.TemplateCodeID

		// Define the regular expression pattern for the template code
		templateCodePattern := `^[a-f0-9]{8}_[a-f0-9]{4}_[a-f0-9]{4}_[a-f0-9]{4}_[a-f0-9]{12}:[a-z0-9]+$`
		templateCodeRegex := regexp.MustCompile(templateCodePattern)

		// If template code is empty or invalid, return error template code is required
		if templateCode == "" {
			return nil, errors.New("template code is required")
		}
		if !templateCodeRegex.MatchString(templateCode) {
			return nil, errors.New("invalid template code")
		}

		// Convert message data to JSON
		messageDataJSON, err := json.Marshal(messageData)
		if err != nil {
			return nil, err
		}

		// Send HTTP request
		url := g.OCAWABASEURL + "/api/v2/push/message"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageDataJSON))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+g.OCAWAToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("Failed to send notification")
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, err
		}

	}

	response := map[string]interface{}{
		"message": "Whatsapp sent successfully",
		"status":  "success",
	}

	log.Println("Response:", response)

	return response, nil
}
