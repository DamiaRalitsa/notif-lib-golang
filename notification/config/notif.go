package config

type NotifConfig struct {
	FabdCoreUrl   string `env:"FABD_CORE_URL" envDefault:"" json:"notif_fabd_core_url"`
	EmailHost     string `env:"EMAIL_HOST" envDefault:"" json:"notif_email_host"`
	EmailPort     string `env:"EMAIL_PORT" envDefault:"" json:"notif_email_port"`
	EmailUserName string `env:"EMAIL_USERNAME" envDefault:"" json:"notif_email_username"`
	EmailPassword string `env:"EMAIL_PASSWORD" envDefault:"" json:"notif_email_password"`
	OCAWABASEURL  string `env:"OCA_WA_BASE_URL" envDefault:"" json:"notif_oca_wa_base_url"`
	OCAWAToken    string `env:"OCA_WA_TOKEN" envDefault:"" json:"notif_oca_wa_token"`
	BellType      string `env:"BELL_TYPE" envDefault:"" json:"notif_bell_type"`
	BellHost      string `env:"BELL_HOST" envDefault:"" json:"notif_bell_host"`
	BellPort      string `env:"BELL_PORT" envDefault:"" json:"notif_bell_port"`
	BellUsername  string `env:"BELL_USERNAME" envDefault:"" json:"notif_bell_username"`
	BellPassword  string `env:"BELL_PASSWORD" envDefault:"" json:"notif_bell_password"`
	BellDatabase  string `env:"BELL_DATABASE" envDefault:"" json:"notif_bell_database"`
}

func (c *NotifConfig) InitEnv() error {
	return envLoader(c, OptionsEnv{DotEnv: true, Prefix: "NOTIF_"})
}
