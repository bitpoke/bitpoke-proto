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

package webhook

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
)

// addToServerFuncs is a list of functions to add all webhooks to the server
var addToServerFuncs []func(manager.Manager, *webhook.Server) error

// newServer creates a new webhook server
func NewServer(m manager.Manager) (*webhook.Server, error) {
	selectors, err := labels.ConvertSelectorToLabelsMap(options.WebhookServiceSelector)
	if err != nil {
		return nil, err
	}

	opts := webhook.ServerOptions{
		CertDir:          "/tmp/cert",
		BootstrapOptions: &webhook.BootstrapOptions{},
	}

	if len(options.WebhookSecretName) > 0 {
		opts.BootstrapOptions.Secret = &types.NamespacedName{
			Namespace: options.WebhookNamespace,
			Name:      options.WebhookSecretName,
		}
	}

	if len(options.WebhookService) > 0 {
		opts.BootstrapOptions.Service = &webhook.Service{
			Namespace: options.WebhookNamespace,
			Name:      options.WebhookService,
			// Selectors should select the pods that runs this webhook server.
			Selectors: selectors,
		}
	}

	if len(options.WebhookHost) > 0 {
		opts.BootstrapOptions.Host = &options.WebhookHost
	}
	server, err := webhook.NewServer("foo-admission-server", m, opts)

	if err != nil {
		return nil, err
	}
	return server, nil
}

// AddToManager add all webhooks to the server
func AddToManager(m manager.Manager) error {
	server, err := NewServer(m)
	if err != nil {
		return err
	}

	for _, fn := range addToServerFuncs {
		if err := fn(m, server); err != nil {
			return err
		}
	}
	return m.Add(server)
}
