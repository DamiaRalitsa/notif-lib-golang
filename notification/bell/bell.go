package bell

import "context"

type NotifBellClient interface {
	SendBell(ctx context.Context, payload NotificationPayload) error
	SendBellBroadcast(ctx context.Context, userIdentifiers []UserIdentifier, payload []NotificationPayload) error
}
