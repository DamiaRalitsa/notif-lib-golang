package oca

type OCA struct {
	PhoneNumber []string `json:"phone_number"`
	MessageData Message  `json:"message_data"`
}

type MessageData struct {
	PhoneNumber string  `json:"phone_number"`
	Message     Message `json:"message"`
}

type Message struct {
	Type     string   `json:"type"`
	Template Template `json:"template"`
}

type Template struct {
	TemplateCodeID string    `json:"template_code_id" validate:"required"`
	Payload        []Payload `json:"payload"`
}

type Payload struct {
	Position   string        `json:"position"`
	Parameters []interface{} `json:"parameters"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SubParameter struct {
	SubType    string      `json:"sub_type"`
	Index      string      `json:"index"`
	Parameters []Parameter `json:"parameters"`
}
