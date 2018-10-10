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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// NewLimitRangeSyncer returns a new syncer.Interface for reconciling Gitea Secret
func NewLimitRangeSyncer(proj *dashboardv1alpha1.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "presslabs-dashboard",
			Namespace: getNamespaceName(proj),
		},
	}

	return syncer.NewObjectSyncer("LimitRange", proj, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.LimitRange)
		out.Labels = getDefaultLabels(proj)

		out.Spec.Limits = []corev1.LimitRangeItem{
			{
				Type: corev1.LimitTypeContainer,
				Default: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("1Gi"),
					corev1.ResourceCPU:    resource.MustParse("200m"),
				},
				DefaultRequest: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("128Mi"),
					corev1.ResourceCPU:    resource.MustParse("100m"),
				},
			},
		}

		return nil
	})
}
