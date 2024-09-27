package bell

type NotificationPayload struct {
	UserID  string      `json:"user_id"`
	Type    string      `json:"type"`
	Name    string      `json:"name"`
	Email   string      `json:"email"`
	Phone   string      `json:"phone"`
	Icon    string      `json:"icon"`
	Path    string      `json:"path"`
	Content interface{} `json:"content"`
	Color   string      `json:"color,omitempty"`
}

type UserIdentifier struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

type NotificationPayloadBroadcast struct {
	Type    string      `json:"type"`
	Icon    string      `json:"icon"`
	Path    string      `json:"path"`
	Content interface{} `json:"content"`
	Color   string      `json:"color,omitempty"`
}
