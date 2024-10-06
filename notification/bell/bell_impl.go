package bell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gateway struct {
	FabdBaseUrl string
	ApiKey      string
}

func NewNotifBellHandler() (NotifBellClient, error) {
	config, err := cfg.InitEnv(cfg.BELL)
	if err != nil {
		return nil, err
	}
	g := &gateway{
		FabdBaseUrl: config.BellConfig.FabdBaseUrl,
		ApiKey:      config.BellConfig.ApiKey,
	}
	return g, err
}

func (g *gateway) SendBell(ctx context.Context, payload NotificationPayload) error {
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

func (g *gateway) SendBellBroadcast(ctx context.Context, userIdentifiers []UserIdentifier, payloads []NotificationPayloadBroadcast) error {
	start := time.Now()
	defer func() {
		log.Printf("sendNotif took %v", time.Since(start))
	}()

	var wg sync.WaitGroup
	errChan := make(chan error, len(userIdentifiers))

	notificationPayloads := make(chan NotificationPayload, len(userIdentifiers))

	if len(userIdentifiers) > 0 {
		payload := payloads[0]
		for _, user := range userIdentifiers {
			wg.Add(1)
			go func(user UserIdentifier) {
				defer wg.Done()

				notificationPayload := NotificationPayload{
					UserID:      user.UserID,
					Type:        payload.Type,
					Icon:        payload.Icon,
					Path:        payload.Path,
					Content:     payload.Content,
					Color:       payload.Color,
					IsRead:      payload.IsRead,
					MsgType:     payload.MsgType,
					Channel:     payload.Channel,
					EcosystemID: payload.EcosystemID,
				}

				if err := validatePayload(notificationPayload); err != nil {
					log.Printf("Validation error for user %s: %v", user.UserID, err)
					return
				}

				notificationPayloads <- notificationPayload
			}(user)
		}
	} else {
		for _, payload := range payloads {
			notificationPayloads <- NotificationPayload{
				UserID:      "",
				Type:        payload.Type,
				Icon:        payload.Icon,
				Path:        payload.Path,
				Content:     payload.Content,
				Color:       payload.Color,
				IsRead:      payload.IsRead,
				MsgType:     payload.MsgType,
				Channel:     payload.Channel,
				EcosystemID: payload.EcosystemID,
			}
		}
	}

	go func() {
		wg.Wait()
		close(notificationPayloads)
	}()

	var payloadList []NotificationPayload
	for payload := range notificationPayloads {
		payloadList = append(payloadList, payload)
	}

	log.Printf("Prepared %d notification payloads", len(payloadList))

	pushStart := time.Now()
	if err := g.pushNotifBulk(payloadList); err != nil {
		log.Printf("Error sending notifications: %v", err)
		select {
		case errChan <- fmt.Errorf("failed to send broadcast notifications"):
		default:
		}
	}
	log.Printf("pushNotifBulk took %v", time.Since(pushStart))

	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *gateway) pushNotif(payload NotificationPayload) error {
	url := g.FabdBaseUrl + "/v4/webhooks/notification"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", g.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("Response from external endpoint:", resp.Status)
	return nil
}

func (g *gateway) pushNotifBulk(payload []NotificationPayload) error {
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
