/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"

	"github.com/presslabs/dashboard/pkg/internal/project"
)

// NewPrometheusServiceAccountSyncer returns a new syncer.Interface for reconciling Prometheus ServiceAccount
func NewPrometheusServiceAccountSyncer(proj *project.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(project.PrometheusServiceAccount)

	obj := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(project.PrometheusServiceAccount),
			Namespace: proj.ComponentName(project.Namespace),
		},
	}

	return syncer.NewObjectSyncer("PrometheusServiceAccount", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.ServiceAccount)
		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		return nil
	})
}
