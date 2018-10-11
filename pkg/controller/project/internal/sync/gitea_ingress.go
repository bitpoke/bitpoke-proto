/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

// NewGiteaIngressSyncer returns a new syncer.Interface for reconciling Gitea Ingress
func NewGiteaIngressSyncer(proj *project.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(project.GiteaIngress)

	obj := &extv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(project.GiteaIngress),
			Namespace: proj.ComponentName(project.Namespace),
		},
	}

	return syncer.NewObjectSyncer("GiteaIngress", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*extv1beta1.Ingress)
		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		out.Spec.Rules = []extv1beta1.IngressRule{
			{
				Host: giteaDomain(proj),
				IngressRuleValue: extv1beta1.IngressRuleValue{
					HTTP: &extv1beta1.HTTPIngressRuleValue{
						Paths: []extv1beta1.HTTPIngressPath{
							{
								Path: "/",
								Backend: extv1beta1.IngressBackend{
									ServiceName: proj.ComponentName(project.GiteaService),
									ServicePort: intstr.FromString("http"),
								},
							},
						},
					},
				},
			},
		}
		return nil
	})
}
