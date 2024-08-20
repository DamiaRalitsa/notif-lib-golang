# NotifConfig Library

The `NotifConfig` library is used to configure notification services by loading environment variables. This library ensures that all necessary configuration parameters are set up correctly for your notification services.

## Required Environment Variables

To use this library, make sure the following environment variables are set in your environment:

- `FABD_CORE_URL`: URL for the FABD core service.
- `EMAIL_HOST`: Host for the email service.
- `EMAIL_PORT`: Port for the email service.
- `EMAIL_USERNAME`: Username for the email service.
- `EMAIL_PASSWORD`: Password for the email service.
- `OCA_WA_BASE_URL`: Base URL for the OCA WA service.
- `OCA_WA_TOKEN`: Token for the OCA WA service.
- `BELL_TYPE`: Type of the Bell service.
- `BELL_HOST`: Host for the Bell service.
- `BELL_PORT`: Port for the Bell service.
- `BELL_USERNAME`: Username for the Bell service.
- `BELL_PASSWORD`: Password for the Bell service.
- `BELL_DATABASE`: Database name for the Bell service.

## Example

Here is an example of how to set these environment variables in a `.env` file:

```env
FABD_CORE_URL=https://example.com/fabd-core
EMAIL_HOST=smtp.example.com
EMAIL_PORT=587
EMAIL_USERNAME=user@example.com
EMAIL_PASSWORD=yourpassword
OCA_WA_BASE_URL=https://example.com/oca-wa
OCA_WA_TOKEN=yourtoken
BELL_TYPE=yourbelltype
BELL_HOST=bell.example.com
BELL_PORT=5432
BELL_USERNAME=belluser
BELL_PASSWORD=bellpassword
BELL_DATABASE=belldb