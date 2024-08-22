package config

type EmailConfig struct {
	EmailHost     string `env:"EMAIL_HOST" envDefault:"" json:"notif_email_host" validate:"required"`
	EmailPort     string `env:"EMAIL_PORT" envDefault:"" json:"notif_email_port" validate:"required"`
	EmailUserName string `env:"EMAIL_USERNAME" envDefault:"" json:"notif_email_username" validate:"required"`
	EmailPassword string `env:"EMAIL_PASSWORD" envDefault:"" json:"notif_email_password" validate:"required"`
}

type OCAConfig struct {
	OCAWABASEURL string `env:"OCA_WA_BASE_URL" envDefault:"" json:"notif_oca_wa_base_url" validate:"required"`
	OCAWAToken   string `env:"OCA_WA_TOKEN" envDefault:"" json:"notif_oca_wa_token" validate:"required"`
}

type BellConfig struct {
	FabdCoreUrl string `env:"FABD_CORE_URL" envDefault:"" json:"notif_fabd_core_url" validate:"required"`
	BellApiKey  string `env:"BELL_API_KEY" envDefault:"" json:"notif_bell_api_key" validate:"required"`
}

func InitEnv(c any) error {
	return envLoader(c, OptionsEnv{DotEnv: true, Prefix: "NOTIF_"})
}
