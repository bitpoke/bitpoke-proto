/*
Copyright 2019 Pressinfra SRL

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

// NewNamespaceSyncer returns a new syncer.Interface for reconciling Project Namespace
func NewNamespaceSyncer(proj *project.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: proj.Spec.NamespaceName,
		},
	}

	return syncer.NewObjectSyncer("Namespace", nil, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Namespace)

		out.Labels = labels.Merge(labels.Merge(out.Labels, proj.Labels()), controllerLabels)
		out.Labels["presslabs.com/kind"] = "project"
		out.Annotations = map[string]string{
			"presslabs.com/created-by":   proj.Annotations["presslabs.com/created-by"],
			"presslabs.com/display-name": proj.Annotations["presslabs.com/display-name"],
		}

		return nil
	})
}
