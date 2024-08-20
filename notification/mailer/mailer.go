package mailer

import (
	"context"
)

type SmtpClient interface {
	SendEmailWithFilePaths(ctx context.Context, mail MailWithoutAttachments, filePaths []string) (data interface{}, err error)
	SendEmail(ctx context.Context, mail Mail) (data interface{}, err error)
}
