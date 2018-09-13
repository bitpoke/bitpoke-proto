/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// NewGiteaServiceSyncer returns a new syncer.Interface for reconciling Gitea Service
func NewGiteaServiceSyncer(proj *dashboardv1alpha1.Project) syncer.Interface {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.GetGiteaServiceName(),
			Namespace: proj.GetNamespaceName(),
		},
	}

	return syncer.New("GiteaService", proj, obj, func(existing runtime.Object) error {
		out := existing.(*corev1.Service)
		out.Labels = GetGiteaPodLabels(proj)

		out.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "http",
				Port:       int32(80),
				TargetPort: intstr.FromInt(giteaHTTPPort),
			},
		}
		out.Spec.Selector = GetGiteaLabels(proj)

		return nil
	})
}
