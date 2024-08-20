package mailer

type Mail struct {
	To          []string      `json:"to"`
	Subject     string        `json:"subject"`
	Message     string        `json:"message"`
	Attachments []Attachments `json:"attachments"`
}
type Attachments struct {
	FileName    string `json:"file_name"`
	Content     []byte `json:"content"`
	Encoding    string `json:"encoding,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

type MailWithoutAttachments struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
}
