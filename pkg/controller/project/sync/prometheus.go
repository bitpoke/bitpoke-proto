/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	defaultScrapeInterval     = "10s"
	defaultEvaluationInterval = "30s"
)

const (
	EventReasonPrometheusFailed  EventReason = "PrometheusFailed"
	EventReasonPrometheusUpdated EventReason = "PrometheusUpdated"
)

type PrometheusSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *monitoringv1.Prometheus
}

func NewPrometheusSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) *PrometheusSyncer {
	return &PrometheusSyncer{
		scheme:   r,
		existing: &monitoringv1.Prometheus{},
		p:        p,
		key:      p.GetPrometheusKey(),
	}
}

func (s *PrometheusSyncer) GetKey() types.NamespacedName                 { return s.key }
func (s *PrometheusSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func (s *PrometheusSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*monitoringv1.Prometheus)
	out.Labels = s.p.GetPrometheusLabels()

	out.Spec = monitoringv1.PrometheusSpec{
		ScrapeInterval:     defaultScrapeInterval,
		EvaluationInterval: defaultEvaluationInterval,
		ServiceMonitorSelector: &metav1.LabelSelector{
			MatchLabels: s.p.GetProjectLabel(),
		},
	}

	return out, nil
}

func (s *PrometheusSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return EventReasonPrometheusUpdated
	} else {
		return EventReasonPrometheusFailed
	}
}
