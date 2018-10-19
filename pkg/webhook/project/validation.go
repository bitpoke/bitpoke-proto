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

package project

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

	"github.com/presslabs/dashboard/pkg/internal/project"
)

type projectValidation struct {
	client  client.Client
	decoder types.Decoder
}

// projectValidation implements admission.Handler.
var _ admission.Handler = &projectValidation{}

// StatusError is a custom error type that also contains a status code
type StatusError struct {
	err        error
	statusCode int32
}

// Error returns the error string
func (e *StatusError) Error() string {
	return e.err.Error()
}

// StatusCode returns the StatusError status code
func (e *StatusError) StatusCode() int32 {
	return e.statusCode
}

// NewStatusError returns a new StatusError object (surprisingly)
func NewStatusError(status int32, err error) *StatusError {
	return &StatusError{
		err:        err,
		statusCode: status,
	}
}

func (a *projectValidation) Handle(ctx context.Context, req types.Request) types.Response {
	proj := project.New(&corev1.Namespace{})

	if err := a.decoder.Decode(req, proj.Unwrap()); err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	if err := a.validateProjectFn(proj); err != nil {
		return admission.ErrorResponse(err.StatusCode(), err)
	}

	// admission.PatchResponse generates a Response containing patches.
	return admission.ValidationResponse(true, "the project is valid")
}

func (a *projectValidation) validateProjectFn(o *project.Project) *StatusError {
	if err := o.ValidateMetadata(); err != nil {
		return NewStatusError(http.StatusBadRequest, err)
	}

	return nil
}

// projectValidation implements inject.Client.
var _ inject.Client = &projectValidation{}

// InjectClient injects the client into the projectValidation
func (a *projectValidation) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// projectValidation implements inject.Decoder.
var _ inject.Decoder = &projectValidation{}

// InjectDecoder injects the decoder into the projectValidation
func (a *projectValidation) InjectDecoder(d types.Decoder) error {
	a.decoder = d
	return nil
}

// AddToServer register an webhook to the server
func AddToServer(m manager.Manager, server *webhook.Server) error {
	wh, err := builder.NewWebhookBuilder().
		Name("projectvalidation.dashboard.presslabs.com").
		Validating().
		Operations(admissionregistrationv1beta1.Create).
		Path("/projectvalidation").
		ForType(&corev1.Pod{}).
		Handlers(&projectValidation{}).
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
