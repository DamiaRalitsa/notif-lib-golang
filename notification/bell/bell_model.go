package bell

type NotificationPayload struct {
	UserID      string      `json:"user_id"`
	Type        string      `json:"type"`
	Icon        string      `json:"icon"`
	Path        string      `json:"path"`
	Content     interface{} `json:"content"`
	Color       string      `json:"color"`
	IsRead      bool        `json:"is_read"`
	MsgType     string      `json:"msg_type"`
	Channel     string      `json:"channel"`
	EcosystemID string      `json:"ecosystem_id"`
}

type UserIdentifier struct {
	UserID string `json:"user_id"`
}
