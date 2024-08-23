# NotifConfig Library

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