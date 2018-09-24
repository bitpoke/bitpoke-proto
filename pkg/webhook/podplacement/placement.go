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

package podplacement

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
)

type podPlacement struct {
	client  client.Client
	decoder types.Decoder
}

// podPlacement implements admission.Handler.
var _ admission.Handler = &podPlacement{}

func (a *podPlacement) Handle(ctx context.Context, req types.Request) types.Response {
	// example

	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)

	}
	copy := pod.DeepCopy()

	err = a.mutatePodsFn(copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)

	}
	// admission.PatchResponse generates a Response containing patches.
	return admission.PatchResponse(pod, copy)

	//return types.Response{}
}

// mutatePodsFn add an annotation to the given pod
// example
func (a *podPlacement) mutatePodsFn(pod *corev1.Pod) error {
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}

	}
	pod.Annotations["example-mutating-admission-webhook"] = "foo"
	return nil

}

// podPlacement implements inject.Client.
var _ inject.Client = &podPlacement{}

// InjectClient injects the client into the podPlacement
func (a *podPlacement) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// podPlacement implements inject.Decoder.
var _ inject.Decoder = &podPlacement{}

// InjectDecoder injects the decoder into the podPlacement
func (a *podPlacement) InjectDecoder(d types.Decoder) error {
	a.decoder = d
	return nil
}

// AddToServer register an webhook to the server
func AddToServer(m manager.Manager, server *webhook.Server) error {
	wh, err := builder.NewWebhookBuilder().
		Mutating().
		Operations(admissionregistrationv1beta1.Create).
		ForType(&corev1.Pod{}).
		Handlers(&podPlacement{}).
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
