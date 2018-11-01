/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package validating

import (
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

var (
	// Builders contain admission webhook builders
	Builders = map[string]*builder.WebhookBuilder{}
	// HandlerMap contains admission webhook handlers
	HandlerMap = map[string][]admission.Handler{}
)
