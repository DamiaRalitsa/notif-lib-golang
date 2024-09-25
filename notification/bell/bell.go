package bell

import "context"

type NotifBellClient interface {
	SendBell(ctx context.Context, payload NotificationPayload) error
	SendBellBroadcast(ctx context.Context, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error
	SendBell2(ctx context.Context, payload NotificationPayloads) error
	SendBellBroadcast2(ctx context.Context, userIdentifiers []UserIdentifier, payload NotificationPayloadBroadcast) error
}
