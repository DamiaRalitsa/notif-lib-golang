package mailer

type Mail struct {
	To           []string               `json:"to" validate:"required"`
	CC           []string               `json:"cc"`
	BCC          []string               `json:"bcc"`
	Subject      string                 `json:"subject" validate:"required"`
	TemplateCode string                 `json:"template_code" validate:"required"`
	Data         map[string]interface{} `json:"data" validate:"required"`
	Attachments  []Attachment           `json:"attachments"`
}
type Attachment struct {
	FileName string `json:"file_name"`
	Path     string `json:"path"`
}

type MailWithoutAttachments struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
	Text    string   `json:"text,omitempty"`
}
