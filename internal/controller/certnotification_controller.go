/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	v1 "github.com/dcristobalhMad/cert-notifier-controller/api/v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slack-go/slack"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CertNotificationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *CertNotificationReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	var certNotification v1.CertNotification
	if err := r.Get(ctx, req.NamespacedName, &certNotification); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	var certs certmanagerv1.CertificateList
	if err := r.List(ctx, &certs); err != nil {
		return reconcile.Result{}, err
	}

	for _, cert := range certs.Items {
		if cert.Status.NotAfter != nil {
			daysToExpiry := cert.Status.NotAfter.Sub(time.Now()).Hours() / 24
			if daysToExpiry <= float64(certNotification.Spec.ExpiryNotificationDays) {
				if err := r.sendNotifications(ctx, certNotification, cert); err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	return reconcile.Result{RequeueAfter: 24 * time.Hour}, nil // Reconcile every day
}

func (r *CertNotificationReconciler) sendNotifications(ctx context.Context, certNotification v1.CertNotification, cert certmanagerv1.Certificate) error {
	// Send email
	if err := r.sendEmail(ctx, certNotification.Spec.EmailConfig); err != nil {
		return err
	}
	// Send Slack message
	if err := r.sendSlackMessage(ctx, certNotification.Spec.SlackConfig); err != nil {
		return err
	}
	// Send Telegram message
	if err := r.sendTelegramMessage(ctx, certNotification.Spec.TelegramConfig); err != nil {
		return err
	}
	return nil
}

func (r *CertNotificationReconciler) getSecretValue(ctx context.Context, selector corev1.SecretKeySelector, namespace string) (string, error) {
	var secret corev1.Secret
	if err := r.Get(ctx, client.ObjectKey{Name: selector.Name, Namespace: namespace}, &secret); err != nil {
		return "", err
	}

	value, exists := secret.Data[selector.Key]
	if !exists {
		return "", fmt.Errorf("key %s not found in secret %s", selector.Key, selector.Name)
	}

	return string(value), nil
}

func (r *CertNotificationReconciler) sendEmail(ctx context.Context, config v1.EmailConfig) error {
	password, err := r.getSecretValue(ctx, config.Password, "default")
	if err != nil {
		return err
	}

	msg := "Subject: Certificate Expiry Notice\n\nYour certificate is close to expiry."
	err = smtp.SendMail(fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort),
		smtp.PlainAuth("", config.From, password, config.SMTPHost),
		config.From, []string{config.To}, []byte(msg))
	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return err
	}
	fmt.Println("Sent email successfully")
	return nil
}

func (r *CertNotificationReconciler) sendSlackMessage(ctx context.Context, config v1.SlackConfig) error {
	token, err := r.getSecretValue(ctx, config.Token, "default")
	if err != nil {
		return err
	}

	api := slack.New(token)
	msg := "Your certificate is close to expiry."
	_, _, err = api.PostMessage(config.ChannelID, slack.MsgOptionText(msg, false))
	if err != nil {
		fmt.Printf("slack error: %s", err)
		return err
	}
	fmt.Println("Sent Slack message successfully")
	return nil
}

func (r *CertNotificationReconciler) sendTelegramMessage(ctx context.Context, config v1.TelegramConfig) error {
	botToken, err := r.getSecretValue(ctx, config.BotToken, "default")
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Printf("telegram error: %s", err)
		return err
	}
	msg := tgbotapi.NewMessage(config.ChatID, "Your certificate is close to expiry.")
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Printf("telegram error: %s", err)
		return err
	}
	fmt.Println("Sent Telegram message successfully")
	return nil
}

func (r *CertNotificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.CertNotification{}).
		Complete(r)
}
