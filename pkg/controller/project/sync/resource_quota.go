/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	// EventReasonResourceQuotaFailed is the event reason for a failed ResourceQuota reconcile
	EventReasonResourceQuotaFailed EventReason = "ResourceQuotaFailed"
	// EventReasonResourceQuotaUpdated is the event reason for a successful ResourceQuota reconcile
	EventReasonResourceQuotaUpdated EventReason = "ResourceQuotaUpdated"
)

var (
	defaultQuotaValues = corev1.ResourceList{
		corev1.ResourceRequestsCPU:    resource.MustParse("4"),
		corev1.ResourceRequestsMemory: resource.MustParse("15Gi"),
		corev1.ResourceLimitsCPU:      resource.MustParse("8"),
		corev1.ResourceLimitsMemory:   resource.MustParse("32Gi"),
		corev1.ResourcePods:           resource.MustParse("20"),
	}
)

// resourceQuotaSyncer defines the Syncer for ResourceQuota
type resourceQuotaSyncer struct {
	scheme   *runtime.Scheme
	proj     *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.ResourceQuota
}

// NewResourceQuotaSyncer returns a new sync.Interface for reconciling ResourceQuota
func NewResourceQuotaSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &resourceQuotaSyncer{
		scheme:   r,
		existing: &corev1.ResourceQuota{},
		proj:     p,
		key:      p.GetResourceQuotaKey(),
	}
}

// GetKey returns the resourceQuotaSyncer key through which an existing object may be identified
func (s *resourceQuotaSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *resourceQuotaSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func defaultOrMaxValue(rl corev1.ResourceList, resource corev1.ResourceName) resource.Quantity {
	defaultResource := defaultQuotaValues[resource]
	if existingResource, ok := rl[resource]; !ok {
		return defaultResource
	} else { // nolint
		if defaultResource.Value() > existingResource.Value() {
			return defaultResource
		}
		return existingResource
	}
}

// T is the transform function used to reconcile the ResourceQuota object
func (s *resourceQuotaSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.ResourceQuota)

	out.Labels = s.proj.GetDefaultLabels()

	out.Spec = corev1.ResourceQuotaSpec{
		Hard: corev1.ResourceList{
			corev1.ResourceRequestsCPU:    defaultOrMaxValue(out.Spec.Hard, corev1.ResourceRequestsCPU),
			corev1.ResourceRequestsMemory: defaultOrMaxValue(out.Spec.Hard, corev1.ResourceRequestsMemory),
			corev1.ResourceLimitsCPU:      defaultOrMaxValue(out.Spec.Hard, corev1.ResourceLimitsCPU),
			corev1.ResourceLimitsMemory:   defaultOrMaxValue(out.Spec.Hard, corev1.ResourceLimitsMemory),
			corev1.ResourcePods:           defaultOrMaxValue(out.Spec.Hard, corev1.ResourcePods),
		},
	}

	err := controllerutil.SetControllerReference(s.proj, out, s.scheme)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *resourceQuotaSyncer) GetErrorEventReason(err error) EventReason {
	if err != nil {
		return EventReasonResourceQuotaFailed
	}
	return EventReasonResourceQuotaUpdated
}
