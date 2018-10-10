/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// NewGiteaServiceSyncer returns a new syncer.Interface for reconciling Gitea Service
func NewGiteaServiceSyncer(proj *dashboardv1alpha1.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      giteaServiceName(proj),
			Namespace: getNamespaceName(proj),
		},
	}

	return syncer.NewObjectSyncer("GiteaService", proj, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Service)
		out.Labels = giteaPodLabels(proj)

		if !labels.Equals(giteaLabels(proj), out.Spec.Selector) {
			if out.ObjectMeta.CreationTimestamp.IsZero() {
				out.Spec.Selector = giteaLabels(proj)
			} else {
				return fmt.Errorf("service selector is immutable")
			}
		}

		if len(out.Spec.Ports) != 2 {
			out.Spec.Ports = make([]corev1.ServicePort, 2)
		}

		out.Spec.Ports[0].Name = "http"
		out.Spec.Ports[0].Port = int32(80)
		out.Spec.Ports[0].TargetPort = intstr.FromInt(giteaHTTPPort)

		out.Spec.Ports[1].Name = "ssh"
		out.Spec.Ports[1].Port = int32(22)
		out.Spec.Ports[1].TargetPort = intstr.FromInt(giteaSSHPort)

		return nil
	})
}
