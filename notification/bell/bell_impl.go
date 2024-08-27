package bell

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	cfg "github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type gateway struct {
	FabdCoreUrl string
	ApiKey      string
}

func NewNotifBellHandler() (NotifBellClient, error) {
	config, err := cfg.InitEnv(cfg.BELL)
	if err != nil {
		return nil, err
	}
	g := &gateway{
		FabdCoreUrl: config.BellConfig.FabdCoreUrl,
		ApiKey:      config.BellConfig.BellApiKey,
	}
	return g, err
}

func (g *gateway) SendBell(ctx context.Context, payload NotificationPayload) error {

	// TODO : Go validator
	if payload.UserID == "" || payload.Type == "" || payload.Name == "" || payload.Email == "" || payload.Icon == "" || payload.Path == "" || payload.Content == nil {
		return errors.New("missing required fields in the payload")
	}

	start := time.Now()

	err := g.pushNotif(payload)
	if err != nil {
		return fmt.Errorf("failed to send bell notifications: %v", err)
	}

	log.Printf("sendBell took %v", time.Since(start))
	return nil
}

func (g *gateway) SendBellBroadcast(ctx context.Context, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error {
	if len(userIdentifiers) == 0 {
		return errors.New("user identifiers array is empty")
	}
	start := time.Now()

	for _, user := range userIdentifiers {
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

		err := g.pushNotif(notificationPayload)
		if err != nil {
			log.Printf("Error sending notification to user %s: %v", user.UserID, err)
			return errors.New("failed to send broadcast notifications")
		}
	}

	log.Printf("sendBellBroadcast took %v", time.Since(start))
	return nil
}

func (g *gateway) pushNotif(payload NotificationPayload) error {
	url := g.FabdCoreUrl + "/v4/notifications"
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
