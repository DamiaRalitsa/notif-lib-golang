package oca

import "context"

type OCAClient interface {
	SendWhatsapp(ctx context.Context, body OCA) (data interface{}, err error)
}
