package bell

type NotifBellClient interface {
	SendBell(payload NotificationPayload) error
	SendBellBroadcast(userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error
}
