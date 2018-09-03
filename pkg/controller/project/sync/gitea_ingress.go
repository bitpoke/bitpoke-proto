/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	// GiteaIngressFailed is the event reason for a failed Gitea Ingress reconcile
	GiteaIngressFailed EventReason = "GiteaIngressFailed"
	// GiteaIngressUpdated is the event reason for a successful Gitea Ingress reconcile
	GiteaIngressUpdated EventReason = "GiteaIngressUpdated"
)

type giteaIngressSyncer struct {
	scheme   *runtime.Scheme
	proj     *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *extv1beta1.Ingress
}

// NewGiteaIngressSyncer returns a new sync.Interface for reconciling Gitea Ingress
func NewGiteaIngressSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &giteaIngressSyncer{
		scheme:   r,
		existing: &extv1beta1.Ingress{},
		proj:     p,
		key:      p.GetGiteaIngressKey(),
	}
}

// GetKey returns the giteaIngressSyncer key through which an existing object may be identified
func (s *giteaIngressSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *giteaIngressSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Gitea Ingress
func (s *giteaIngressSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*extv1beta1.Ingress)
	out.Labels = GetGiteaPodLabels(s.proj)

	out.Spec.Rules = []extv1beta1.IngressRule{
		{
			Host: s.proj.GetGiteaDomain(),
			IngressRuleValue: extv1beta1.IngressRuleValue{
				HTTP: &extv1beta1.HTTPIngressRuleValue{
					Paths: []extv1beta1.HTTPIngressPath{
						{
							Path: "/",
							Backend: extv1beta1.IngressBackend{
								ServiceName: s.proj.GetGiteaServiceName(),
								ServicePort: intstr.FromString("http"),
							},
						},
					},
				},
			},
		},
	}
	return out, nil
}

func (s *giteaIngressSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaIngressUpdated
	}
	return GiteaIngressFailed
}
