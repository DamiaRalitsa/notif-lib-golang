# notif-lib-golang
FABD Notification Library

# configuration

The `NotifConfig` library is used to configure notification services by loading environment variables. This library ensures that all necessary configuration parameters are set up correctly for your notification services.

## Required Environment Variables

To use this library, make sure the following environment variables are set in your environment based on the service you are using.

### Bell Service

- `NOTIF_BELL_TYPE`: Type of the Bell service.
- `NOTIF_BELL_HOST`: Host for the Bell service.
- `NOTIF_BELL_PORT`: Port for the Bell service.
- `NOTIF_BELL_USERNAME`: Username for the Bell service.
- `NOTIF_BELL_PASSWORD`: Password for the Bell service.
- `NOTIF_BELL_DATABASE`: Database name for the Bell service.

### OCA Service

- `NOTIF_OCA_WA_BASE_URL`: Base URL for the OCA WA service.
- `NOTIF_OCA_WA_TOKEN`: Token for the OCA WA service.

### Email Service

- `NOTIF_EMAIL_HOST`: Host for the email service.
- `NOTIF_EMAIL_PORT`: Port for the email service.
- `NOTIF_EMAIL_USERNAME`: Username for the email service.
- `NOTIF_EMAIL_PASSWORD`: Password for the email service.

### FABD Core Service

- `NOTIF_FABD_CORE_URL`: URL for the FABD core service.

## Example

Here is an example of how to set these environment variables in a `.env` file:

```env
# FABD Core Service
NOTIF_FABD_CORE_URL=https://example.com/fabd-core

# Email Service
NOTIF_EMAIL_HOST=smtp.example.com
NOTIF_EMAIL_PORT=587
NOTIF_EMAIL_USERNAME=user@example.com
NOTIF_EMAIL_PASSWORD=yourpassword

# OCA Service
NOTIF_OCA_WA_BASE_URL=https://example.com/oca-wa
NOTIF_OCA_WA_TOKEN=yourtoken

# Bell Service
NOTIF_BELL_TYPE=yourbelltype
NOTIF_BELL_HOST=bell.example.com
NOTIF_BELL_PORT=5432
NOTIF_BELL_USERNAME=belluser
NOTIF_BELL_PASSWORD=bellpassword
NOTIF_BELL_DATABASE=belldb

# notif bell

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-go/bell@latest
go get github.com/lib/pq
```

## Configuration

Create a configuration struct to hold your database connection details:

```sh
package main

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/DamiaRalitsa/notif-lib-go/bell"
)

func main() {

    // Initialize the notification handler
    notifHandler := bell.NewNotifBellHandler()

    // Example usage
    db, err := sql.Open("postgres", "user=username password=password dbname=database sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

```

## Sending Notifications

Send a Single Notification

To send a single notification, use the SendBell method:

```sh
payload := bell.NotificationPayload{
    UserID:  "123",
    Type:    "info",
    Name:    "John Doe",
    Email:   "john.doe@example.com",
    Phone:   "1234567890",
    Icon:    "icon.png",
    Path:    "/path/to/resource",
    Content: map[string]interface{}{"message": "Hello, World!"},
    Color:   "primary",
}

err := notifHandler.SendBell(db, payload)
if err != nil {
    log.Fatal(err)
}
```

Send Broadcast Notifications

To send broadcast notifications to multiple users, use the SendBellBroadcast method:
```sh
userIdentifiers := []bell.UserIdentifier{
    {UserID: "123", Name: "John Doe", Email: "john.doe@example.com", Phone: "1234567890"},
    {UserID: "456", Name: "Jane Smith", Email: "jane.smith@example.com", Phone: "0987654321"},
}

broadcastPayload := bell.NotificationPayloadBroadcast{
    Type:    "info",
    Icon:    "icon.png",
    Path:    "/path/to/resource",
    Content: map[string]interface{}{"message": "Hello, Everyone!"},
    Color:   "primary",
}

err := notifHandler.SendBellBroadcast(db, userIdentifiers, broadcastPayload)
if err != nil {
    log.Fatal(err)
}
```

# notif Mailer

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-go/mailer@latest
```

## Configuration

Create a configuration struct to hold your mailer connection details:

```sh
package main

import (
    "context"
    "log"
    "os"

    "github.com/DamiaRalitsa/notif-lib-go/mailer"
)

func main() {

    // Initialize the Mailer handler
    mailerHandler := mailer.NewMailerHandler()
    
}
```

## Sending Notifications

Send Emails

To send a notification, use the SendEmail method

```sh
to := []string{"recipient1@example.com", "recipient2@example.com"}
subject := "Test Subject"
message := "This is a test email with attachments."

 mail := mailer.Mail{
  	To:      emailRecipients,
  	Subject: emailSubject,
  	Message: emailMessage,
  }

response, err := mailerHandler.SendEmail(context.Background(), mail)
if err != nil {
    log.Fatal(err)
}

log.Println("Response:", response)
```

# notif OCA

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-go
```

## Configuration

Create a configuration struct to hold your OCA connection details:

```sh
package main

import (
    "context"
    "log"
    "os"

    "github.com/DamiaRalitsa/notif-lib-go/oca@latest"
)

func main() {

    // Initialize the OCA handler
    ocaHandler := oca.NewOCAHandler()

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
				TemplateCodeID: os.Getenv("OCA_WA_TEMPLATE_CODE"),
				Payload: []oca.Payload{
					{
						Position: "body",
						Parameters: []interface{}{
							oca.Parameter{
								Type: "text",
								Text: "tesst",
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
										Text: "tessssstttt",
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