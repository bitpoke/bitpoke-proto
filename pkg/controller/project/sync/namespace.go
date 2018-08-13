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
	EventReasonNamespaceFailed  EventReason = "NamespaceFailed"
	EventReasonNamespaceUpdated EventReason = "NamespaceUpdated"
)

type NamespaceSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.Namespace
}

func NewNamespaceSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) *NamespaceSyncer {
	return &NamespaceSyncer{
		scheme:   r,
		p:        p,
		existing: &corev1.Namespace{},
		key:      p.GetNamespaceKey(),
	}
}

func (s *NamespaceSyncer) GetKey() types.NamespacedName                 { return s.key }
func (s *NamespaceSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func (s *NamespaceSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.Namespace)

	out.Labels = s.p.GetDefaultLabels()

	controllerutil.SetControllerReference(s.p, out, s.scheme)

	return out, nil
}

func (s *NamespaceSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return EventReasonNamespaceUpdated
	} else {
		return EventReasonNamespaceFailed
	}
}
