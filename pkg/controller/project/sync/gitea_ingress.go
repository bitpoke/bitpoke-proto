/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"fmt"

	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
)

const (
	GiteaIngressFailed  EventReason = "GiteaIngressFailed"
	GiteaIngressUpdated EventReason = "GiteaIngressUpdated"
)

type GiteaIngressSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *extv1beta1.Ingress
}

func NewGiteaIngressSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) *GiteaIngressSyncer {
	return &GiteaIngressSyncer{
		scheme:   r,
		existing: &extv1beta1.Ingress{},
		p:        p,
		key:      p.GetGiteaIngressKey(),
	}
}

func (s *GiteaIngressSyncer) GetKey() types.NamespacedName                 { return s.key }
func (s *GiteaIngressSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func GetProjectGiteaPath(p *dashboardv1alpha1.Project) string {
	return fmt.Sprintf("%s.%s", p.GetProjectNamespacedName(), options.GitBaseDomainURL.URL.Path)
}

func (s *GiteaIngressSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*extv1beta1.Ingress)
	out.Labels = GetGiteaPodLabels(s.p)

	out.Spec.Rules = []extv1beta1.IngressRule{
		{
			Host: GetProjectGiteaPath(s.p),
			IngressRuleValue: extv1beta1.IngressRuleValue{
				HTTP: &extv1beta1.HTTPIngressRuleValue{
					Paths: []extv1beta1.HTTPIngressPath{
						extv1beta1.HTTPIngressPath{
							Path: "/",
							Backend: extv1beta1.IngressBackend{
								ServiceName: s.p.GetGiteaServiceName(),
								ServicePort: intstr.FromString("http"),
							},
						},
					},
				},
			},
		},
	}
	return out, nil
}

func (s *GiteaIngressSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaIngressUpdated
	} else {
		return GiteaIngressFailed
	}
}
