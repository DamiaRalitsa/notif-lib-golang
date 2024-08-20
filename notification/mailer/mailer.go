package mailer

import (
	"context"
	"time"
)

type SmtpClient interface {
	SendEmailWithFilePaths(ctx context.Context, mail MailWithoutAttachments, filePaths []string) (data interface{}, err error)
	SendEmail(ctx context.Context, mail Mail, timer time.Time) (data interface{}, err error)
}
