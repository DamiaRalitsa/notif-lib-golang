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