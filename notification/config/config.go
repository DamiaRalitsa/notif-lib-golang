package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	godotenv.Load()
}

const (
	OCA           = "oca"
	BELL          = "bell"
	EMAIL         = "email"
	API           = "api"
	EnvPrefix     = "NOTIF_"
	EmailHost     = EnvPrefix + "EMAIL_HOST"
	EmailPort     = EnvPrefix + "EMAIL_PORT"
	EmailUserName = EnvPrefix + "EMAIL_USERNAME"
	EmailPassword = EnvPrefix + "EMAIL_PASSWORD"

	OCAWABASEURL = EnvPrefix + "OCA_WA_BASE_URL"
	OCAWAToken   = EnvPrefix + "OCA_WA_TOKEN"

	BellAPIKEY = EnvPrefix + "BELL_API_KEY"

	FabdCoreUrl = EnvPrefix + "BELL_FABD_CORE_URL"
	ApiKey      = EnvPrefix
)

type Config struct {
	EmailConfig EmailConfig
	OCAConfig   OCAConfig
	BellConfig  BellConfig
	ApiConfig   ApiConfig
}

type EmailConfig struct {
	EmailHost     string `json:"notif_email_host" validate:"required"`
	EmailPort     string `json:"notif_email_port" validate:"required"`
	EmailUserName string `json:"notif_email_username" validate:"required"`
	EmailPassword string `json:"notif_email_password" validate:"required"`
}

type OCAConfig struct {
	OCAWABASEURL string `json:"notif_oca_wa_base_url" validate:"required"`
	OCAWAToken   string `json:"notif_oca_wa_token" validate:"required"`
}

type BellConfig struct {
	FabdCoreUrl string `json:"notif_fabd_core_url" validate:"required"`
	BellApiKey  string `json:"notif_bell_api_key" validate:"required"`
}

type ApiConfig struct {
	FabdCoreUrl string `json:"notif_fabd_core_url" validate:"required"`
	ApiKey      string `json:"notif_api_key" validate:"required"`
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func InitEnv(configName string) (Config, error) {
	config := Config{}
	switch configName {
	case EMAIL:
		emailConfig := EmailConfig{
			EmailHost:     getEnv(EmailHost),
			EmailPort:     getEnv(EmailPort),
			EmailUserName: getEnv(EmailUserName),
			EmailPassword: getEnv(EmailPassword),
		}
		if err := validateEnv(&emailConfig); err != nil {
			log.Printf("Email configuration is not valid: %v", err)
			return Config{}, err
		}
		config.EmailConfig = emailConfig
	case OCA:
		ocaConfig := OCAConfig{
			OCAWABASEURL: getEnv(OCAWABASEURL),
			OCAWAToken:   getEnv(OCAWAToken),
		}
		if err := validateEnv(&ocaConfig); err != nil {
			log.Printf("OCA configuration is not valid: %v", err)
			return Config{}, err
		}
		config.OCAConfig = ocaConfig
	case BELL:
		bellConfig := BellConfig{
			FabdCoreUrl: getEnv(config.BellConfig.FabdCoreUrl),
			BellApiKey:  getEnv(BellAPIKEY),
		}
		if err := validateEnv(&bellConfig); err != nil {
			log.Printf("Bell configuration is not valid: %v", err)
			return Config{}, err
		}
		config.BellConfig = bellConfig
	case API:
		apiConfig := ApiConfig{
			FabdCoreUrl: getEnv(config.ApiConfig.FabdCoreUrl),
			ApiKey:      getEnv(ApiKey),
		}
		if err := validateEnv(&apiConfig); err != nil {
			log.Printf("Api configuration is not valid: %v", err)
			return Config{}, err
		}
		config.ApiConfig = apiConfig
	}
	return config, nil
}

func validateEnv(cfg any) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("Validation failed: %v", err)
		return err
	}
	return nil
}
