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
	GiteaServiceFailed  EventReason = "GiteaServiceFailed"
	GiteaServiceUpdated EventReason = "GiteaServiceUpdated"
)

type GiteaServiceSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.Service
}

func NewGiteaServiceSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) *GiteaServiceSyncer {
	return &GiteaServiceSyncer{
		scheme:   r,
		existing: &corev1.Service{},
		p:        p,
		key:      p.GetGiteaServiceKey(),
	}
}

func (s *GiteaServiceSyncer) GetKey() types.NamespacedName                 { return s.key }
func (s *GiteaServiceSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func (s *GiteaServiceSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.Service)
	out.Labels = GetGiteaPodLabels(s.p)

	out.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       int32(80),
			TargetPort: intstr.FromInt(giteaHTTPInternalPort),
		},
	}
	out.Spec.Selector = GetGiteaLabels(s.p)

	return out, nil
}

func (s *GiteaServiceSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaServiceUpdated
	} else {
		return GiteaServiceFailed
	}
}
