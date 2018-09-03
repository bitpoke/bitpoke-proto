/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	// GiteaServiceFailed is the event reason for a failed Gitea Service reconcile
	GiteaServiceFailed EventReason = "GiteaServiceFailed"
	// GiteaServiceUpdated is the event reason for a successful Gitea Service reconcile
	GiteaServiceUpdated EventReason = "GiteaServiceUpdated"
)

type giteaServiceSyncer struct {
	scheme   *runtime.Scheme
	proj     *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.Service
}

// NewGiteaServiceSyncer returns a new sync.Interface for reconciling Gitea Service
func NewGiteaServiceSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &giteaServiceSyncer{
		scheme:   r,
		existing: &corev1.Service{},
		proj:     p,
		key:      p.GetGiteaServiceKey(),
	}
}

// GetKey returns the giteaServiceSyncer key through which an existing object may be identified
func (s *giteaServiceSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *giteaServiceSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Gitea Service
func (s *giteaServiceSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.Service)
	out.Labels = GetGiteaPodLabels(s.proj)

	out.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       int32(80),
			TargetPort: intstr.FromInt(giteaHTTPPort),
		},
	}
	out.Spec.Selector = GetGiteaLabels(s.proj)

	return out, nil
}

func (s *giteaServiceSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaServiceUpdated
	}
	return GiteaServiceFailed
}
