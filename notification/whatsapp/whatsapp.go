package whatsapp

import "context"

type WhatsappClient interface {
	SendWhatsapp(ctx context.Context, body Whatsapp) (data interface{}, err error)
}
