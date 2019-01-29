/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package validating

import (
	"context"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	organization "github.com/presslabs/dashboard/pkg/internal/organization"
)

func init() {
	webhookName := "validating-create-organization"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &NamespaceCreateHandler{})
}

// NamespaceCreateHandler handles Namespace
type NamespaceCreateHandler struct {
	Client client.Client

	// Decoder decodes objects
	Decoder types.Decoder
}

func (h *NamespaceCreateHandler) validatingNamespaceFn(obj *organization.Organization) (bool, string, error) {
	if obj.Namespace.Labels["presslabs.com/kind"] != "organization" {
		return true, "not an organization, skipping validation", nil
	}

	if err := obj.ValidateMetadata(); err != nil {
		return false, "validation failed", err
	}

	return true, "allowed to be admitted", nil
}

var _ admission.Handler = &NamespaceCreateHandler{}

// Handle handles admission requests.
func (h *NamespaceCreateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	org := organization.New(&corev1.Namespace{})

	err := h.Decoder.Decode(req, org.Unwrap())
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := h.validatingNamespaceFn(org)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

var _ inject.Client = &NamespaceCreateHandler{}

// InjectClient injects the client into the NamespaceCreateHandler
func (h *NamespaceCreateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &NamespaceCreateHandler{}

// InjectDecoder injects the decoder into the NamespaceCreateHandler
func (h *NamespaceCreateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
