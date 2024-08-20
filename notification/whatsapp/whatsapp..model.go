package whatsapp

// Whatsapp represents the structure of the WhatsApp message.
type Whatsapp struct {
	To      []string `json:"to"`
	Type    string   `json:"type"`
	ID      string   `json:"id"`
	Message string   `json:"message"`
}

type Respond struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
