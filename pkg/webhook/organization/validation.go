/*v
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

package organization

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

	organization "github.com/presslabs/dashboard/pkg/internal/organization"
)

type organizationValidation struct {
	client  client.Client
	decoder types.Decoder
}

// organizationValidation implements admission.Handler.
var _ admission.Handler = &organizationValidation{}

func (a *organizationValidation) Handle(ctx context.Context, req types.Request) types.Response {
	org := organization.Wrap(&corev1.Namespace{})

	err := a.decoder.Decode(req, org.Unwrap())
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	err = a.validateOrganizationFn(org)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	// admission.PatchResponse generates a Response containing patches.
	return admission.ValidationResponse(true, "the organization is valid")
}

func (a *organizationValidation) validateOrganizationFn(org *organization.Organization) error {
	return org.ValidateMetadata()
}

// organizationValidation implements inject.Client.
var _ inject.Client = &organizationValidation{}

// InjectClient injects the client into the organizationValidation
func (a *organizationValidation) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// organizationValidation implements inject.Decoder.
var _ inject.Decoder = &organizationValidation{}

// InjectDecoder injects the decoder into the organizationValidation
func (a *organizationValidation) InjectDecoder(d types.Decoder) error {
	a.decoder = d
	return nil
}

// AddToServer register an webhook to the server
func AddToServer(m manager.Manager, server *webhook.Server) error {
	wh, err := builder.NewWebhookBuilder().
		Mutating().
		Operations(admissionregistrationv1beta1.Create).
		ForType(&corev1.Pod{}).
		Handlers(&organizationValidation{}).
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
