/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	// EventReasonNamespaceFailed is the event reason for a failed Namespace reconcile
	EventReasonNamespaceFailed EventReason = "NamespaceFailed"
	// EventReasonNamespaceUpdated is the event reason for a successful Namespace reconcile
	EventReasonNamespaceUpdated EventReason = "NamespaceUpdated"
)

// namespaceSyncer defines the Syncer for Namespace
type namespaceSyncer struct {
	scheme   *runtime.Scheme
	proj     *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.Namespace
}

// NewNamespaceSyncer returns a new sync.Interface for reconciling Namespace
func NewNamespaceSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &namespaceSyncer{
		scheme:   r,
		proj:     p,
		existing: &corev1.Namespace{},
		key:      p.GetNamespaceKey(),
	}
}

// GetKey returns the namespaceSyncer key through which an existing object may be identified
func (s *namespaceSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *namespaceSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Namespace object
func (s *namespaceSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.Namespace)

	out.Labels = s.proj.GetDefaultLabels()

	err := controllerutil.SetControllerReference(s.proj, out, s.scheme)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *namespaceSyncer) GetErrorEventReason(err error) EventReason {
	if err != nil {
		return EventReasonNamespaceFailed
	}
	return EventReasonNamespaceUpdated
}
