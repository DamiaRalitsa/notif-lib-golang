package whatsapp

import (
	"bytes"
	"context"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
)

type WhatsappHandler struct {
	whatsappService WhatsappClient
}

type WhatsappConfig struct {
	BaseURL string
	AppKey  string
	AuthKey string
}

// Gateway represents the structure of the gateway.
type gateway struct {
	BaseURL string
	AppKey  string
	AuthKey string
	// HttpClient *helpers.ToolsAPI
}

// NewWhatsappHandler creates a new WhatsappHandler instance.
func NewWhatsappHandler(whatsappConfig WhatsappConfig) WhatsappClient {
	g := &gateway{
		BaseURL: whatsappConfig.BaseURL,
		AppKey:  whatsappConfig.AppKey,
		AuthKey: whatsappConfig.AuthKey,
	}
	return g
}

// SendWhatsapp sends a WhatsApp message to multiple phone numbers.
func (g gateway) SendWhatsapp(ctx context.Context, body Whatsapp) (data interface{}, err error) {
	url := g.BaseURL

	for _, phoneNumber := range body.To {
		// Create a buffer
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		checkPhoneNumber := phoneNumber[:2]
		if checkPhoneNumber == "08" {
			phoneNumber = "62" + phoneNumber[1:]

		} else if checkPhoneNumber == "+6" {
			phoneNumber = phoneNumber[1:]
		} else if checkPhoneNumber != "62" {
			return nil, errors.New("Invalid phone number")
		}

		log.Println("WHATSAPP SENT TO PHONE NUMBER:", phoneNumber)

		if body.Type == "PO" {
			body.Message = "Halo, Selamat kamu mendapatkan pesanan baru dengan nomor pesanan " + body.ID + ".\n\nSilahkan cek aplikasi untuk melihat detail pesanan.\n\nTerima kasih."
		}
		if body.Type == "customer" {
			body.Message = "Halo, terima kasih telah melakukan pemesanan dengan nomor pesanan " + body.ID + ".\n\nPesanan akan segera kami proses. Mohon ditunggu.\n\nTerima kasih."
		}

		// Write the fields
		_ = writer.WriteField("appkey", g.AppKey)
		_ = writer.WriteField("authkey", g.AuthKey)
		_ = writer.WriteField("to", phoneNumber)
		_ = writer.WriteField("message", body.Message)

		// Close the writer
		err = writer.Close()
		if err != nil {
			return nil, err
		}

		// Create a new request
		req, err := http.NewRequest("POST", url, buf)
		if err != nil {
			return nil, err
		}

		// Set the content type
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("res.Body is not a map")
			log.Println("PATH:", url)
			log.Println("STATUS:", res.StatusCode)
			log.Println("MESSAGE:", res.Status)
			return nil, err
		}
		defer res.Body.Close()
	}

	response := map[string]interface{}{
		"message": "Whatsapp sent successfully",
		"status":  "success",
	}

	return response, nil
}

func (g *gateway) NewWhatsappClient() WhatsappClient {
	// baseUrl := os.Getenv("WHATSAPP_BASE_URL")
	// appKey := os.Getenv("WHATSAPP_APP_KEY")
	// authKey := os.Getenv("WHATSAPP_AUTH_KEY")
	// httpClient := &helpers.ToolsAPI{}
	return &gateway{
		BaseURL: g.BaseURL,
		AppKey:  g.AppKey,
		AuthKey: g.AuthKey,
		// HttpClient: httpClient,
	}
}
