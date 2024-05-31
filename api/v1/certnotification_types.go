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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CertNotificationSpec defines the desired state of CertNotification
type CertNotificationSpec struct {
	EmailConfig            EmailConfig    `json:"emailConfig,omitempty"`
	SlackConfig            SlackConfig    `json:"slackConfig,omitempty"`
	TelegramConfig         TelegramConfig `json:"telegramConfig,omitempty"`
	ExpiryNotificationDays int            `json:"expiryNotificationDays"`
}

type EmailConfig struct {
	From     string                   `json:"from"`
	Password corev1.SecretKeySelector `json:"password"`
	To       string                   `json:"to"`
	SMTPHost string                   `json:"smtpHost"`
	SMTPPort string                   `json:"smtpPort"`
}

type SlackConfig struct {
	Token     corev1.SecretKeySelector `json:"token"`
	ChannelID string                   `json:"channelID"`
}

type TelegramConfig struct {
	BotToken corev1.SecretKeySelector `json:"botToken"`
	ChatID   int64                    `json:"chatID"`
}

// CertNotificationStatus defines the observed state of CertNotification
type CertNotificationStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CertNotification is the Schema for the certnotifications API
type CertNotification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertNotificationSpec   `json:"spec,omitempty"`
	Status CertNotificationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CertNotificationList contains a list of CertNotification
type CertNotificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CertNotification `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CertNotification{}, &CertNotificationList{})
}
