# Cert Notifier Controller

Cert Notifier Controller is a Kubernetes controller built with Kubebuilder that monitors certificates issued by cert-manager and sends notifications when certificates are nearing expiration. Notifications can be sent via email, Slack, and Telegram.

## Features

- Monitors certificates issued by cert-manager
- Sends notifications when certificates are within a configurable number of days from expiration
- Supports multiple notification channels: Email, Slack, and Telegram
- Secrets for sensitive information (e.g., API tokens, passwords) are securely fetched from Kubernetes secrets

## Prerequisites

- Kubernetes cluster (v1.19+)
- cert-manager (v1.0+)
- Kubebuilder (v3.0+)
- Go (v1.19+)

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/your-username/cert-notifier-controller.git
   cd cert-notifier-controller
   ```

2. **Build the controller:**

   ```sh
   make build
   ```

3. **Deploy the controller to your Kubernetes cluster:**

   ```sh
   make deploy
   ```

4. **Apply the `CertNotification` CR:**
   ```sh
   kubectl apply -f config/samples/certmanagerapp_v1_certnotification.yaml
   ```

## Custom Resource Definition

### CertNotification

The `CertNotification` custom resource allows you to configure the notification settings for certificate expiration.

- **expiryNotificationDays:** The number of days before certificate expiration to send a notification.
- **emailConfig:** Configuration for email notifications.
  - **smtpHost:** SMTP server host.
  - **smtpPort:** SMTP server port.
  - **from:** Email address to send from.
  - **to:** Email address to send to.
  - **password:** Reference to a Kubernetes secret containing the SMTP password.
- **slackConfig:** Configuration for Slack notifications.
  - **token:** Reference to a Kubernetes secret containing the Slack API token.
  - **channelID:** Slack channel ID to send messages to.
- **telegramConfig:** Configuration for Telegram notifications.
  - **botToken:** Reference to a Kubernetes secret containing the Telegram bot token.
  - **chatID:** Telegram chat ID to send messages to.

### Example CR

Here is an example of a `CertNotification` custom resource:

```yaml
apiVersion: certmanagerapp.cert-notifier.com/v1
kind: CertNotification
metadata:
  labels:
    app.kubernetes.io/name: cert-notifier-controller
    app.kubernetes.io/managed-by: kustomize
  name: certnotification-sample
spec:
  emailConfig:
    from: "your-email@example.com"
    password:
      name: notification-secrets
      key: email-password
    to: "recipient@example.com"
    smtpHost: "smtp.example.com"
    smtpPort: "587"
  slackConfig:
    token:
      name: notification-secrets
      key: slack-token
    channelID: "your-channel-id"
  telegramConfig:
    botToken:
      name: notification-secrets
      key: telegram-bot-token
    chatID: 123456789
  expiryNotificationDays: 30
```
