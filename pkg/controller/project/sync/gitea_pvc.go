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
	// GiteaPVCFailed is the event reason for a failed Gitea PVC reconcile
	GiteaPVCFailed EventReason = "GiteaPVCFailed"
	// GiteaPVCUpdated is the event reason for a successful Gitea PVC reconcile
	GiteaPVCUpdated EventReason = "GiteaPVCUpdated"
)

type giteaPVCSyncer struct {
	scheme   *runtime.Scheme
	proj     *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.PersistentVolumeClaim
}

// NewGiteaPVCSyncer returns a new sync.Interface for reconciling Gitea PVC
func NewGiteaPVCSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &giteaPVCSyncer{
		scheme:   r,
		existing: &corev1.PersistentVolumeClaim{},
		proj:     p,
		key:      p.GetGiteaPVCKey(),
	}
}

// GetKey returns the giteaPVCSyncer key through which an existing object may be identified
func (s *giteaPVCSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *giteaPVCSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Gitea PVC
func (s *giteaPVCSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.PersistentVolumeClaim)
	out.Labels = GetGiteaPodLabels(s.proj)

	out.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteOnce,
	}
	out.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: GetGiteaLabels(s.proj),
	}
	out.Spec.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"storage": resource.MustParse(giteaRequestsStorage),
		},
	}

	return out, nil
}

func (s *giteaPVCSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaPVCUpdated
	}
	return GiteaPVCFailed
}
