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