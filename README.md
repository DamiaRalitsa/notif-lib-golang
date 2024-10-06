# notif-lib-golang
FABD Notification Library

# configuration

The `NotifConfig` library is used to configure notification services by loading environment variables. This library ensures that all necessary configuration parameters are set up correctly for your notification services.

## Required Environment Variables

To use this library, make sure the following environment variables are set in your environment based on the service you are using.

### Bell Service

- `NOTIF_BELL_API_KEY`: ApiKey for the Bell service.

### OCA Service

- `NOTIF_OCA_WA_BASE_URL`: Base URL for the OCA WA service.
- `NOTIF_OCA_WA_TOKEN`: Token for the OCA WA service.

### Email Service

- `NOTIF_EMAIL_HOST`: Host for the email service.
- `NOTIF_EMAIL_PORT`: Port for the email service.
- `NOTIF_EMAIL_USERNAME`: Username for the email service.
- `NOTIF_EMAIL_PASSWORD`: Password for the email service.

### FABD Core Service

- `FABD_BASE_URL`: URL for the FABD core service.
- `API_KEY`: ApiKey for the FABD core service.

## Example

Here is an example of how to set these environment variables in a `.env` file:

```env
# Bell Service
NOTIF_BELL_API_KEY=yourapikey

# OCA Service
NOTIF_OCA_WA_BASE_URL=https://example.com/oca-wa
NOTIF_OCA_WA_TOKEN=yourtoken
NOTIF_OCA_WA_TEMPLATE_CODE=yourtemplatecode

# Email Service
NOTIF_EMAIL_HOST=smtp.example.com
NOTIF_EMAIL_PORT=587
NOTIF_EMAIL_USERNAME=user@example.com
NOTIF_EMAIL_PASSWORD=yourpassword

# FABD Core Service
FABD_BASE_URL=yourdomain.com
API_KEY=yourapikey

```

# notif bell

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-golang/notification/bell@latest
```

## Configuration

Create a configuration struct to hold your database connection details:

```sh
package main

import (
    "database/sql"
    "log"
    "os"

    "github.com/DamiaRalitsa/notif-lib-golang/notification/bell"
)

func main() {

    // Initialize the notification handler
    notifHandler, err := bell.NewNotifBellApiHandler()
    if err != nil {
        log.Fatalf("Failed to initialize notification gateway: %v", err)
    }

```

## Sending Notifications

Send a Single Notification

To send a single notification, use the SendBell method:

```sh
payload := bell.NotificationPayload{
    UserID:      "123",
    Type:        "info",
    Icon:        "icon.png",
    Path:        "/path/to/resource",
    Content:     map[string]interface{}{"message": "Hello, World!"},
    Color:       "primary",
    MsgType:     "alert",
    Channel:     "email",
    EcosystemID: "ecosystem_123",
}

err := notifHandler.SendBell(ctx, payload)
if err != nil {
    log.Fatal(err)
}

log.Println("Single notification sent successfully")
```

Send Broadcast Notifications

To send broadcast notifications to multiple users, use the SendBellBroadcast method:
```sh
userIdentifiers := []bell.UserIdentifier{
    {
        UserID: "123",
    },
    {
        UserID: "456",
    },
}

broadcastPayload := bell.NotificationPayloadBroadcast{
    Type:        "info",
    Icon:        "icon.png",
    Path:        "/path/to/resource",
    Content:     map[string]interface{}{"message": "Hello, Everyone!"},
    Color:       "primary",
    MsgType:     "alert",
    Channel:     "email",
    EcosystemID: "ecosystem_123",
}

err := notifHandler.SendBellBroadcast(ctx, userIdentifiers, broadcastPayload)
if err != nil {
    log.Fatal(err)
}

log.Println("Broadcast notifications sent successfully")
```

# notif Mailer

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-golang/notification/mailer@latest
```

## Configuration

Create a configuration struct to hold your mailer connection details:

```sh
package main

import (
    "context"
    "log"
    "os"

    "github.com/DamiaRalitsa/notif-lib-golang/notification/mailer"
)

func main() {

    // Initialize the Mailer handler
    mailerHandler, err := mailer.NewMailerHandler()
    if err != nil {
	log.Fatalf("Failed to initialize mailer gateway: %v", err)
	}
}
```

## Sending Notifications

Send Emails

To send a notification, use the SendEmail method

```sh
emailPayload := mailer.Mail{
		To:           []string{"example@gmail.com", "example2@gmail.com"},
		CC:           []string{"example@gmail.com"},
		BCC:          []string{"example@gmail.com"},
		Subject:      "Test Subject",
		TemplateCode: "test_template_code",
		Data:         map[string]interface{}{"otp_code": "123456"},
		Attachments: []mailer.Attachment{
			{
				FileName: "test.txt",
				Path:     "./test.txt",
			},
		},
	}

response, err := emailHandler.SendEmail(ctx, emailPayload)
if err != nil {
	log.Fatalf("Error sending email: %v", err)
}
log.Printf("Email sent successfully: %v", response)
```

# notif OCA

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-golang/notification/oca@latest
```

## Configuration

Create a configuration struct to hold your OCA connection details:

```sh
package main

import (
    "context"
    "log"
    "os"

    "github.com/DamiaRalitsa/notif-lib-golang/notification/oca"
)

func main() {

    // Initialize the OCA handler
    ocaHandler, err := oca.NewOCAHandler()

}
```

## Sending Notifications

Send Notification

To send notification, use the SendWhatsapp method

```sh
body := oca.OCA{
    PhoneNumber: []string{"0812345678","0812345678"},
    MessageData: oca.Message{
			Type: "template",
			Template: oca.Template{
				TemplateCodeID: os.Getenv("NOTIF_OCA_WA_TEMPLATE_CODE"),
				Payload: []oca.Payload{
					{
						Position: "body",
						Parameters: []interface{}{
							oca.Parameter{
								Type: "text",
								Text: "123456",
							},
						},
					},
					{
						Position: "button",
						Parameters: []interface{}{
							oca.SubParameter{
								SubType: "url",
								Index:   "0",
								Parameters: []oca.Parameter{
									{
										Type: "text",
										Text: "123456",
									},
								},
							},
						},
					},
				},
			},
		},
}

response, err := ocaHandler.SendWhatsapp(context.Background(), body)
if err != nil {
    log.Fatal(err)
}

log.Println("Response:", response)
```
