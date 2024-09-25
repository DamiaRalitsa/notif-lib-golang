package bell

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gatewayApi struct {
	FabdCoreUrl string
	ApiKey      string
}

func NewNotifBellApiHandler() (NotifBellClient, error) {
	config, err := cfg.InitEnv(cfg.API)
	if err != nil {
		return nil, err
	}
	g := &gatewayApi{
		FabdCoreUrl: config.ApiConfig.FabdCoreUrl,
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

func (g *gatewayApi) SendBellBroadcast(ctx context.Context, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error {
	if len(userIdentifiers) == 0 {
		return errors.New("user identifiers array is empty")
	}

	start := time.Now()
	defer func() {
		log.Printf("sendNotif took %v", time.Since(start))
	}()

	var wg sync.WaitGroup
	errChan := make(chan error, len(userIdentifiers))

	notificationPayloads := NotificationPayloads{}.NotificationPayload

	var mu sync.Mutex

	for _, user := range userIdentifiers {
		wg.Add(1)
		go func(user UserIdentifier) {
			defer wg.Done()

			notificationPayload := NotificationPayload{
				UserID:  user.UserID,
				Type:    payload.Type,
				Name:    user.Name,
				Email:   user.Email,
				Phone:   user.Phone,
				Icon:    payload.Icon,
				Path:    payload.Path,
				Content: payload.Content,
				Color:   payload.Color,
			}

			if err := validatePayload(notificationPayload); err != nil {
				return
			}

			mu.Lock()
			notificationPayloads = append(notificationPayloads, notificationPayload)
			mu.Unlock()

		}(user)
	}

	if err := g.pushNotifBulk(notificationPayloads); err != nil {
		log.Printf("Error sending notifications: %v", err)
		select {
		case errChan <- fmt.Errorf("failed to send broadcast notifications"):
		default:
		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *gatewayApi) pushNotif(payload NotificationPayload) error {
	url := g.FabdCoreUrl + "/v4/notification-service/notifications/bell"
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	log.Printf("Payload: %s", string(jsonData))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", g.ApiKey)
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
	url := g.FabdCoreUrl + "/v4/notification-service/notifications/bell/bulk"
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	log.Printf("Payload: %s", string(jsonData))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", g.ApiKey)
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
