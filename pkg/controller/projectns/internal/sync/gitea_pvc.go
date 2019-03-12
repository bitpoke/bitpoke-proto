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
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

var (
	resGiteaRequestsStorage = resource.MustParse(giteaRequestsStorage)
)

// NewGiteaPVCSyncer returns a new syncer.Interface for reconciling Gitea PVC
func NewGiteaPVCSyncer(proj *projectns.ProjectNamespace, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(projectns.GiteaPVC)

	obj := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(projectns.GiteaPVC),
			Namespace: proj.ComponentName(projectns.Namespace),
		},
	}

	return syncer.NewObjectSyncer("GiteaPVC", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.PersistentVolumeClaim)

		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		out.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		}
		out.Spec.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				"storage": resGiteaRequestsStorage,
			},
		}

		return nil
	})
}
