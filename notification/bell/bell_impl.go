package bell

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/DamiaRalitsa/notif-lib-golang/notification/config"
)

type BellHandler struct {
	notifBellService NotifBellClient
}

type gateway struct {
	Type     string
	Host     string
	Port     string
	Username string
	Password string
	Database string
	// HttpClient *helpers.ToolsAPI
}

// NewNotifBellHandler creates a new NotifBellHandler instance.
func NewNotifBellHandler() NotifBellClient {
	config := &config.NotifConfig{}
	config.InitEnv()
	g := &gateway{
		Type:     config.BellType,
		Host:     config.BellHost,
		Port:     config.BellPort,
		Username: config.BellUsername,
		Password: config.BellPassword,
		Database: config.BellDatabase,
	}
	return g
}

func (g *gateway) SendBell(db *sql.DB, payload NotificationPayload) error {
	// Simple input validation
	if payload.UserID == "" || payload.Type == "" || payload.Name == "" || payload.Email == "" || payload.Icon == "" || payload.Path == "" || payload.Content == nil {
		return errors.New("missing required fields in the payload")
	}

	// Start timing the operation
	start := time.Now()

	insertQuery := `
        INSERT INTO dev_fabd_user_core_owner.notifications (
            user_id, "type", "name", email, phone, icon, "path", "content", color
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	content, err := json.Marshal(payload.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %v", err)
	}

	values := []interface{}{
		payload.UserID,
		payload.Type,
		payload.Name,
		payload.Email,
		payload.Phone,
		payload.Icon,
		payload.Path,
		content,
		payload.Color,
	}

	_, err = db.Exec(insertQuery, values...)
	if err != nil {
		log.Printf("Error inserting notification: %v", err)
		return errors.New("failed to send notification")
	}

	// End timing and log the duration
	log.Printf("sendBell took %v", time.Since(start))
	return nil
}

func (g *gateway) SendBellBroadcast(db *sql.DB, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error {
	// Validate that the userIdentifiers array is not empty
	if len(userIdentifiers) == 0 {
		return errors.New("user identifiers array is empty")
	}

	// Start timing the operation
	start := time.Now()

	// Prepare the base insert query
	insertQuery := `
        INSERT INTO dev_fabd_user_core_owner.notifications (
            user_id, "type", "name", email, phone, icon, "path", "content", color
        ) VALUES
    `

	// Construct the values part of the query
	valueRows := ""
	values := []interface{}{}
	for i, user := range userIdentifiers {
		if i > 0 {
			valueRows += ", "
		}
		baseIndex := i * 9
		valueRows += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", baseIndex+1, baseIndex+2, baseIndex+3, baseIndex+4, baseIndex+5, baseIndex+6, baseIndex+7, baseIndex+8, baseIndex+9)
		values = append(values, user.UserID, payload.Type, user.Name, user.Email, user.Phone, payload.Icon, payload.Path, payload.Content, payload.Color)
	}

	fullQuery := insertQuery + valueRows

	_, err := db.Exec(fullQuery, values...)
	if err != nil {
		log.Printf("Error broadcasting notifications: %v", err)
		return errors.New("failed to send broadcast notifications")
	}

	// End timing and log the duration
	log.Printf("sendBellBroadcast took %v", time.Since(start))
	return nil
}
