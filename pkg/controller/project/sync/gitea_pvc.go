/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	GiteaPVCFailed  EventReason = "GiteaPVCFailed"
	GiteaPVCUpdated EventReason = "GiteaPVCUpdated"
)

type GiteaPVCSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.PersistentVolumeClaim
}

func NewGiteaPVCSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) *GiteaPVCSyncer {
	return &GiteaPVCSyncer{
		scheme:   r,
		existing: &corev1.PersistentVolumeClaim{},
		p:        p,
		key:      p.GetGiteaPVCKey(),
	}
}

func (s *GiteaPVCSyncer) GetKey() types.NamespacedName                 { return s.key }
func (s *GiteaPVCSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func (s *GiteaPVCSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.PersistentVolumeClaim)
	out.Labels = GetGiteaPodLabels(s.p)

	out.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteOnce,
	}
	out.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: GetGiteaLabels(s.p),
	}
	out.Spec.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"storage": resource.MustParse("2Gi"),
		},
	}

	return out, nil
}

func (s *GiteaPVCSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaPVCUpdated
	} else {
		return GiteaPVCFailed
	}
}
