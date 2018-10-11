/*
Copyright 2018 Pressinfra SRL.

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

package options

import (
	"net/url"

	"github.com/spf13/pflag"
)

// GitBaseDomainURL is the the base domain used to obtain the git repo domain for projects
var GitBaseDomainURL = "git.presslabs.net"

var (
	// WebhookNamespace is the namespace for webhooks
	WebhookNamespace = "default"
	// WebhookSecretName is the secret for webhooks
	// nolint: gosec
	WebhookSecretName = "presslabs-dashboard-admission-webhook-cert"
	// WebhookServiceName is the service for webhooks
	WebhookServiceName = "presslabs-dashboard-admission-webhhok"
	// WebhookServiceSelector is the selector for webhook
	WebhookServiceSelector = "app.kubernetes.io/name=presslabs-dashboard,app.kubernetes.io/component=controller-manager"
	// WebhookHost is the host of the webhook server
	WebhookHost = "localhost"
	// WebhookPort is the port of the webbhook server
	WebhookPort = 4433
)

// AddToFlagSet add options to a FlagSet
func AddToFlagSet(flag *pflag.FlagSet) {
	flag.StringVar(&GitBaseDomainURL, "git-base-domain", GitBaseDomainURL, "The base git domain")

	flag.StringVar(&WebhookNamespace, "webhook-namespace", WebhookNamespace, "The webhook namespace")
	flag.StringVar(&WebhookSecretName, "webhook-secret-name", WebhookSecretName, "The webhook secret name")
	flag.StringVar(&WebhookServiceName, "webhook-service-name", WebhookServiceName, "The webhook server name")
	flag.StringVar(&WebhookServiceSelector, "webhook-service-selector", WebhookServiceSelector, "The selector for webhook service")
	flag.StringVar(&WebhookHost, "webhook-host", WebhookHost, "The webhook server host")
	flag.IntVar(&WebhookPort, "webhook-port", WebhookPort, "The webhook server port")
}

// Validate validates the arguments
func Validate() error {
	_, err := url.Parse(GitBaseDomainURL)
	if err != nil {
		return err
	}
	return nil
}
