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
}

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

	if payload.UserID == "" || payload.Type == "" || payload.Name == "" || payload.Email == "" || payload.Icon == "" || payload.Path == "" || payload.Content == nil {
		return errors.New("missing required fields in the payload")
	}

	start := time.Now()

	insertQuery := `
    INSERT INTO notifications (
        user_id, "type", "name", email, phone, icon, "path", "content", color
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	content, err := json.Marshal(payload.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %v", err)
	}

	_, err = stmt.Exec(
		payload.UserID,
		payload.Type,
		payload.Name,
		payload.Email,
		payload.Phone,
		payload.Icon,
		payload.Path,
		content,
		payload.Color,
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	log.Printf("sendBell took %v", time.Since(start))
	return nil
}

func (g *gateway) SendBellBroadcast(db *sql.DB, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error {

	if len(userIdentifiers) == 0 {
		return errors.New("user identifiers array is empty")
	}
	start := time.Now()

	insertQuery := `
        INSERT INTO notifications (
            user_id, "type", "name", email, phone, icon, "path", "content", color
        ) VALUES
    `

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

	stmt, err := db.Prepare(fullQuery)
	if err != nil {
		return errors.New("failed to prepare statement")
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return errors.New("failed to send broadcast notifications")
	}

	log.Printf("sendBellBroadcast took %v", time.Since(start))

	return nil
}
