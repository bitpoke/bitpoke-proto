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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// NewNamespaceSyncer returns a new syncer.Interface for reconciling Namespace
func NewNamespaceSyncer(proj *dashboardv1alpha1.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: getNamespaceName(proj),
		},
	}

	return syncer.NewObjectSyncer("Namespace", proj, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Namespace)

		out.Labels = getDefaultLabels(proj)

		return nil
	})
}
