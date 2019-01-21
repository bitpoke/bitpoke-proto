/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package defaultserver

import (
	"fmt"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/labels"
	// "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"

	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
)

var (
	log        = logf.Log.WithName("default_server")
	builderMap = map[string]*builder.WebhookBuilder{}
	// HandlerMap contains all admission webhook handlers.
	HandlerMap = map[string][]admission.Handler{}
)

// Add adds itself to the manager
// +kubebuilder:webhook:port=7890,cert-dir=/path/to/cert
// +kubebuilder:webhook:service=test-system:webhook-service,selector=app:webhook-server
// +kubebuilder:webhook:secret=test-system:webhook-secret
// +kubebuilder:webhook:mutating-webhook-config-name=mutating-webhook-configuration
// +kubebuilder:webhook:validating-webhook-config-name=validating-webhook-configuration
func Add(mgr manager.Manager) error {
	selectors, err := labels.ConvertSelectorToLabelsMap(options.WebhookServiceSelector)
	if err != nil {
		log.Error(err, "error while starting webhook server")
		return err
	}

	opts := webhook.ServerOptions{
		Port:                          int32(options.WebhookPort),
		CertDir:                       options.WebhookCertDir,
		BootstrapOptions:              &webhook.BootstrapOptions{},
		DisableWebhookConfigInstaller: &options.WebhookDisableBootstrapping,
	}

	if !options.WebhookDisableBootstrapping {
		if len(options.WebhookSecretName) > 0 {
			opts.BootstrapOptions.Secret = &types.NamespacedName{
				Namespace: options.WebhookNamespace,
				Name:      options.WebhookSecretName,
			}
		}

		if len(options.WebhookServiceName) > 0 {
			opts.BootstrapOptions.Service = &webhook.Service{
				Namespace: options.WebhookNamespace,
				Name:      options.WebhookServiceName,
				// Selectors should select the pods that runs this webhook server.
				Selectors: selectors,
			}
		} else if len(options.WebhookHost) > 0 {
			opts.BootstrapOptions.Host = &options.WebhookHost
		}
	}
	srv, err := webhook.NewServer("presslabs-dashboard-admission-server", mgr, opts)
	if err != nil {
		log.Error(err, "error while starting webhook server")
		return err
	}

	webhooks := make([]webhook.Webhook, len(builderMap))
	i := 0
	for k, builder := range builderMap {
		handlers, ok := HandlerMap[k]
		if !ok {
			log.V(1).Info(fmt.Sprintf("can't find handlers for builder: %v", k))
			handlers = []admission.Handler{}
		}
		wh, err := builder.
			Handlers(handlers...).
			WithManager(mgr).
			Build()
		if err != nil {
			log.Error(err, "error while starting webhook server")
			return err
		}
		webhooks[i] = wh
		i++
	}

	return srv.Register(webhooks...)
}
