/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// NewGiteaIngressSyncer returns a new syncer.Interface for reconciling Gitea Ingress
func NewGiteaIngressSyncer(proj *dashboardv1alpha1.Project) syncer.Interface {
	obj := &extv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.GetGiteaIngressName(),
			Namespace: proj.GetNamespaceName(),
		},
	}

	return syncer.New("GiteaIngress", proj, obj, func(existing runtime.Object) error {
		out := existing.(*extv1beta1.Ingress)
		out.Labels = GetGiteaPodLabels(proj)

		out.Spec.Rules = []extv1beta1.IngressRule{
			{
				Host: proj.GetGiteaDomain(),
				IngressRuleValue: extv1beta1.IngressRuleValue{
					HTTP: &extv1beta1.HTTPIngressRuleValue{
						Paths: []extv1beta1.HTTPIngressPath{
							{
								Path: "/",
								Backend: extv1beta1.IngressBackend{
									ServiceName: proj.GetGiteaServiceName(),
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
