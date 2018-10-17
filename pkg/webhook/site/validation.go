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

package site

import (
	"context"
	"net/http"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	"github.com/presslabs/dashboard/pkg/internal/site"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

type siteValidation struct {
	client  client.Client
	decoder types.Decoder
}

// siteValidation implements admission.Handler.
var _ admission.Handler = &siteValidation{}

func (a *siteValidation) Handle(ctx context.Context, req types.Request) types.Response {
	wp := site.New(&wordpressv1alpha1.Wordpress{})

	err := a.decoder.Decode(req, wp.Unwrap())
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	err = a.validateSiteFn(wp)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	// admission.PatchResponse generates a Response containing patches.
	return admission.ValidationResponse(true, "the site is valid")
}

func (a *siteValidation) validateSiteFn(o *site.Site) error {
	return o.ValidateMetadata()
}

// siteValidation implements inject.Client.
var _ inject.Client = &siteValidation{}

// InjectClient injects the client into the siteValidation
func (a *siteValidation) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// siteValidation implements inject.Decoder.
var _ inject.Decoder = &siteValidation{}

// InjectDecoder injects the decoder into the siteValidation
func (a *siteValidation) InjectDecoder(d types.Decoder) error {
	a.decoder = d
	return nil
}

// AddToServer register an webhook to the server
func AddToServer(m manager.Manager, server *webhook.Server) error {
	wh, err := builder.NewWebhookBuilder().
		Mutating().
		Operations(admissionregistrationv1beta1.Create).
		ForType(&corev1.Pod{}).
		Handlers(&siteValidation{}).
		WithManager(m).
		Build()
	if err != nil {
		return err
	}

	err = server.Register(wh)
	if err != nil {
		return err
	}

	return nil
}
