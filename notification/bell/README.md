# notif bell

## Installation

First, install the library using `go get`:

```sh
go get github.com/DamiaRalitsa/notif-lib-go
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
    config := bell.NotifBellConfig{
        Type:     "postgres",
        Host:     os.Getenv("DB_HOST"),
		Port:     5432,
        Username: os.Getenv("DB_USERNAME"),
        Password: os.Getenv("DB_PASSWORD"),
        Database: os.Getenv("DB_DATABASE"),
    }

    // Initialize the notification handler
    notifHandler := bell.NewNotifBellHandler(config)

    // Example usage
    db, err := sql.Open("postgres", "user=username password=password dbname=database sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

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

    err = notifHandler.SendBell(db, payload)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Notification sent successfully!")
}
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

