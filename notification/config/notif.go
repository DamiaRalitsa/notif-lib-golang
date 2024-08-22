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
	FabdCoreUrl  string `env:"FABD_CORE_URL" envDefault:"" json:"notif_fabd_core_url" validate:"required"`
	BellApiKey   string `env:"BELL_API_KEY" envDefault:"" json:"notif_bell_api_key" validate:"required"`
	BellType     string `env:"BELL_TYPE" envDefault:"" json:"notif_bell_type" validate:"required"`
	BellHost     string `env:"BELL_HOST" envDefault:"" json:"notif_bell_host" validate:"required"`
	BellPort     string `env:"BELL_PORT" envDefault:"" json:"notif_bell_port" validate:"required"`
	BellUsername string `env:"BELL_USERNAME" envDefault:"" json:"notif_bell_username" validate:"required"`
	BellPassword string `env:"BELL_PASSWORD" envDefault:"" json:"notif_bell_password" validate:"required"`
	BellDatabase string `env:"BELL_DATABASE" envDefault:"" json:"notif_bell_database" validate:"required"`
}

func InitEnv(c any) error {
	return envLoader(c, OptionsEnv{DotEnv: true, Prefix: "NOTIF_"})
}
