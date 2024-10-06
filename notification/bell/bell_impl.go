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
		ApiKey:      config.BellConfig.BellApiKey,
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

func (g *gateway) SendBellBroadcast(ctx context.Context, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error {
	start := time.Now()
	defer func() {
		log.Printf("sendNotif took %v", time.Since(start))
	}()

	var wg sync.WaitGroup
	errChan := make(chan error, len(userIdentifiers))

	notificationPayloads := []NotificationPayload{}

	var mu sync.Mutex

	if len(userIdentifiers) > 0 {
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
					return
				}

				mu.Lock()
				notificationPayloads = append(notificationPayloads, notificationPayload)
				mu.Unlock()

			}(user)
		}
		wg.Wait()
	} else {
		notificationPayloads = append(notificationPayloads, NotificationPayload{
			UserID:      payload.UserID,
			Type:        payload.Type,
			Icon:        payload.Icon,
			Path:        payload.Path,
			Content:     payload.Content,
			Color:       payload.Color,
			IsRead:      payload.IsRead,
			MsgType:     payload.MsgType,
			Channel:     payload.Channel,
			EcosystemID: payload.EcosystemID,
		})
	}

	if err := g.pushNotifBulk(notificationPayloads); err != nil {
		log.Printf("Error sending notifications: %v", err)
		select {
		case errChan <- fmt.Errorf("failed to send broadcast notifications"):
		default:
		}
	}

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
