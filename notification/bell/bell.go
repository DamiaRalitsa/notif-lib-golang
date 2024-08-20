package bell

import "database/sql"

type NotifBellClient interface {
	SendBell(db *sql.DB, payload NotificationPayload) error
	SendBellBroadcast(db *sql.DB, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error
}
