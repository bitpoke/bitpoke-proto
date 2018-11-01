/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package validating

import (
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

func init() {
	builderName := "validating-create-organization"
	Builders[builderName] = builder.
		NewWebhookBuilder().
		Name(builderName + ".presslabs.com").
		Path("/" + builderName).
		Validating().
		Operations(admissionregistrationv1beta1.Create).
		FailurePolicy(admissionregistrationv1beta1.Fail).
		ForType(&corev1.Namespace{})
}
