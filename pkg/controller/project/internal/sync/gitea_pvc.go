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

var (
	resGiteaRequestsStorage = resource.MustParse(giteaRequestsStorage)
)

// NewGiteaPVCSyncer returns a new syncer.Interface for reconciling Gitea PVC
func NewGiteaPVCSyncer(proj *dashboardv1alpha1.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      giteaPVCName(proj),
			Namespace: getNamespaceName(proj),
		},
	}

	return syncer.NewObjectSyncer("GiteaPVC", proj, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.PersistentVolumeClaim)
		out.Labels = giteaPodLabels(proj)

		out.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		}
		out.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: giteaLabels(proj),
		}
		out.Spec.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				"storage": resGiteaRequestsStorage,
			},
		}

		return nil
	})
}
