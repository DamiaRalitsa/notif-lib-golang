package oca

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gatewayApi struct {
	FabdCoreUrl string
	ApiKey      string
}

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewOCAApiHandler() (OCAClient, error) {
	config, err := cfg.InitEnv(cfg.API)
	if err != nil {
		log.Fatalf("Failed to initialize environment: %v", err)
		return nil, err
	}
	g := &gatewayApi{
		FabdCoreUrl: config.ApiConfig.FabdCoreUrl,
		ApiKey:      config.ApiConfig.ApiKey,
	}
	return g, nil
}

func (g gatewayApi) SendWhatsapp(ctx context.Context, payload OCA) (data interface{}, err error) {
	url := g.FabdCoreUrl + "/v4/notification-service/notifications/whatsapp/oca"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", g.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))
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
