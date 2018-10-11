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
	"fmt"
	"net/http"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	projectinternal "github.com/presslabs/dashboard/pkg/internal/project"
)

var log = logf.Log.WithName("project-controller")

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
	project := &dashboardv1alpha1.Project{}

	if err := a.decoder.Decode(req, project); err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	if err := a.validateProjectFn(ctx, project); err != nil {
		return admission.ErrorResponse(err.StatusCode(), err)
	}

	// admission.PatchResponse generates a Response containing patches.
	return admission.ValidationResponse(true, "the project is valid")
}

func (a *projectValidation) validateProjectFn(ctx context.Context, project *dashboardv1alpha1.Project) *StatusError {
	if err := projectinternal.ValidateMetadata(project); err != nil {
		return NewStatusError(http.StatusBadRequest, err)
	}

	namespaceLabels := map[string]string{
		"presslabs.com/kind": "project",
	}
	for _, key := range projectinternal.RequiredLabels() {
		namespaceLabels[key] = project.Labels[key]
	}
	namespaceAnnotations := map[string]string{}
	for _, key := range projectinternal.RequiredAnnotations() {
		namespaceAnnotations[key] = project.Annotations[key]
	}

	// Try to create the Project's Namespace
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        projectinternal.GetNamespaceName(project),
			Labels:      namespaceLabels,
			Annotations: namespaceAnnotations,
		},
	}
	if err := a.client.Create(ctx, namespace); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return NewStatusError(http.StatusBadRequest, fmt.Errorf("project \"%s\" is not available", project.Name))
		}

		log.Error(err, "unable to create project namespace")

		return NewStatusError(http.StatusInternalServerError, fmt.Errorf("could not create project at the time"))
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
