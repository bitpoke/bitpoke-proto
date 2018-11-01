/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package validating

import (
	"context"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	"github.com/presslabs/dashboard/pkg/internal/site"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

func init() {
	webhookName := "validating-create-site"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &WordpressCreateHandler{})
}

// WordpressCreateHandler handles Wordpress
type WordpressCreateHandler struct {
	// Client  client.Client

	// Decoder decodes objects
	Decoder types.Decoder
}

func (h *WordpressCreateHandler) validatingWordpressFn(obj *site.Site) (bool, string, error) {
	if err := obj.ValidateMetadata(); err != nil {
		return false, "validation failed", err
	}
	return true, "allowed to be admitted", nil
}

var _ admission.Handler = &WordpressCreateHandler{}

// Handle handles admission requests.
func (h *WordpressCreateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	wp := site.New(&wordpressv1alpha1.Wordpress{})

	err := h.Decoder.Decode(req, wp.Unwrap())
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := h.validatingWordpressFn(wp)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

//var _ inject.Client = &WordpressCreateHandler{}
//
//// InjectClient injects the client into the WordpressCreateHandler
//func (h *WordpressCreateHandler) InjectClient(c client.Client) error {
//	h.Client = c
//	return nil
//}

var _ inject.Decoder = &WordpressCreateHandler{}

// InjectDecoder injects the decoder into the WordpressCreateHandler
func (h *WordpressCreateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
