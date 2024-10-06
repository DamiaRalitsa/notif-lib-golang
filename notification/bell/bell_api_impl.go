package bell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gatewayApi struct {
	FabdBaseUrl string
	ApiKey      string
}

func NewNotifBellApiHandler() (NotifBellClient, error) {
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

func (g *gatewayApi) SendBell(ctx context.Context, payload NotificationPayload) error {
	start := time.Now()
	defer func() {
		log.Printf("sendNotif took %v", time.Since(start))
	}()

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := validatePayload(payload); err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := g.pushNotif(payload); err != nil {
			select {
			case errChan <- fmt.Errorf("failed to send bell notifications: %v", err):
			default:
			}
		}
	}()

	wg.Wait()
	close(errChan)

	select {
	case err := <-errChan:
		return err
	default:
	}

	return nil
}

func (g *gatewayApi) SendBellBroadcast(ctx context.Context, userIdentifiers []UserIdentifier, payloads []NotificationPayload) error {
	start := time.Now()
	defer func() {
		log.Printf("sendNotif took %v", time.Since(start))
	}()

	if len(userIdentifiers) == 0 {
		for _, payload := range payloads {
			if err := validatePayload(payload); err != nil {
				log.Printf("Validation error: %v", err)
				return fmt.Errorf("validation error: %v", err)
			}
		}

		log.Printf("Prepared %d notification payloads", len(payloads))

		pushStart := time.Now()
		if err := g.pushNotifBulk(payloads); err != nil {
			log.Printf("Error sending notifications: %v", err)
			return fmt.Errorf("failed to send broadcast notifications")
		}
		log.Printf("pushNotifBulk took %v", time.Since(pushStart))

		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(userIdentifiers))
	notificationPayloads := make(chan NotificationPayload, len(userIdentifiers))

	for _, user := range userIdentifiers {
		wg.Add(1)
		go func(user UserIdentifier) {
			defer wg.Done()

			notificationPayload := payloads[0]
			notificationPayload.UserID = user.UserID

			if err := validatePayload(notificationPayload); err != nil {
				log.Printf("Validation error for user %s: %v", user.UserID, err)
				errChan <- err
				return
			}

			notificationPayloads <- notificationPayload
		}(user)
	}

	go func() {
		wg.Wait()
		close(notificationPayloads)
		close(errChan)
	}()

	var payloadList []NotificationPayload
	for payload := range notificationPayloads {
		payloadList = append(payloadList, payload)
	}

	log.Printf("Prepared %d notification payloads", len(payloadList))

	pushStart := time.Now()
	if err := g.pushNotifBulk(payloadList); err != nil {
		log.Printf("Error sending notifications: %v", err)
		return fmt.Errorf("failed to send broadcast notifications")
	}
	log.Printf("pushNotifBulk took %v", time.Since(pushStart))

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *gatewayApi) pushNotif(payload NotificationPayload) error {
	url := g.FabdBaseUrl + "/v4/webhooks/notifications"
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", g.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Response from external endpoint:", resp.Status)
	return nil
}

func (g *gatewayApi) pushNotifBulk(payload []NotificationPayload) error {
	url := g.FabdBaseUrl + "/v4/webhooks/notifications-bulk"
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", g.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Response from external endpoint:", resp.Status)
	return nil
}
