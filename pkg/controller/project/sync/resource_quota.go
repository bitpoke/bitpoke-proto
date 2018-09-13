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

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
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

// NewResourceQuotaSyncer returns a new syncer.Interface for reconciling ResourceQuota
func NewResourceQuotaSyncer(proj *dashboardv1alpha1.Project) syncer.Interface {
	obj := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.GetResourceQuotaName(),
			Namespace: proj.GetNamespaceName(),
		},
	}

	return syncer.New("ResourceQuota", proj, obj, func(existing runtime.Object) error {
		out := existing.(*corev1.ResourceQuota)

		out.Labels = proj.GetDefaultLabels()

		out.Spec = corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceRequestsCPU:    defaultOrMaxValue(out.Spec.Hard, corev1.ResourceRequestsCPU),
				corev1.ResourceRequestsMemory: defaultOrMaxValue(out.Spec.Hard, corev1.ResourceRequestsMemory),
				corev1.ResourceLimitsCPU:      defaultOrMaxValue(out.Spec.Hard, corev1.ResourceLimitsCPU),
				corev1.ResourceLimitsMemory:   defaultOrMaxValue(out.Spec.Hard, corev1.ResourceLimitsMemory),
				corev1.ResourcePods:           defaultOrMaxValue(out.Spec.Hard, corev1.ResourcePods),
			},
		}

		return nil
	})
}
