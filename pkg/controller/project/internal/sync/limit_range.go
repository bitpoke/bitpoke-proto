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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"

	"github.com/presslabs/dashboard/pkg/internal/project"
)

// NewLimitRangeSyncer returns a new syncer.Interface for reconciling Gitea Secret
func NewLimitRangeSyncer(proj *project.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(project.LimitRange)

	obj := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(project.LimitRange),
			Namespace: proj.ComponentName(project.Namespace),
		},
	}

	return syncer.NewObjectSyncer("LimitRange", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.LimitRange)
		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

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
