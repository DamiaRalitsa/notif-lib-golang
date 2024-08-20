# notif Mailer

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-go
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
    config := mailer.MailerConfig{
        FabdCoreUrl:   os.Getenv("FABD_CORE_URL"),
        EmailHost:     os.Getenv("EMAIL_HOST"),
        EmailPort:     os.Getenv("EMAIL_PORT"),
        EmailUserName: os.Getenv("EMAIL_USERNAME"),
        EmailPassword: os.Getenv("EMAIL_PASSWORD"),
    }

    // Initialize the Mailer handler
    mailerHandler := mailer.NewMailerHandler(config)

    // Example usage
    to := []string{"recipient@example.com"}
    subject := "Test Subject"
    message := "This is a test email."
    filePaths := []string{"path/to/attachment1", "path/to/attachment2"}

    response, err := mailerHandler.SendEmailWithAttachments(context.Background(), to, subject, message, filePaths)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Response:", response)
}
```

## Sending Notifications

Send Emails

To send a notification, use the SendEmailWIthAttachments method

```sh
to := []string{"recipient1@example.com", "recipient2@example.com"}
subject := "Test Subject"
message := "This is a test email with attachments."
filePaths := []string{"path/to/attachment1", "path/to/attachment2"}

response, err := mailerHandler.SendEmailWithAttachments(context.Background(), to, subject, message, filePaths)
if err != nil {
    log.Fatal(err)
}

log.Println("Response:", response)
```