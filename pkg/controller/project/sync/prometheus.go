/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	prometheusName            = "prometheus"
	prometheusBaseImage       = "quay.io/prometheus/prometheus"
	prometheusVersion         = "v2.3.2"
	defaultScrapeInterval     = "10s"
	defaultEvaluationInterval = "30s"
)

const (
	// EventReasonPrometheusFailed is the event reason for a failed Prometheus reconcile
	EventReasonPrometheusFailed EventReason = "PrometheusFailed"
	// EventReasonPrometheusUpdated is the event reason for a successful Prometheus reconcile
	EventReasonPrometheusUpdated EventReason = "PrometheusUpdated"
)

// prometheusSyncer defines the Syncer for Prometheus
type prometheusSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *monitoringv1.Prometheus
}

// GetPrometheusSelector returns a set of labels that can be used to identify Prometheus
// related resources
func GetPrometheusSelector(project *dashboardv1alpha1.Project) labels.Set {
	prometheusLabels := labels.Set{
		"app.kubernetes.io/name": prometheusName,
	}
	return labels.Merge(project.GetDefaultLabels(), prometheusLabels)
}

// GetPrometheusLabels returns a set of labels that should be applied on Prometheus
// related objects that are managed by the project controller
func GetPrometheusLabels(project *dashboardv1alpha1.Project) labels.Set {
	prometheusLabels := labels.Set{
		"app.kubernetes.io/version": prometheusVersion,
	}
	return labels.Merge(GetPrometheusSelector(project), prometheusLabels)
}

// NewPrometheusSyncer returns a new sync.Interface for reconciling Prometheus
func NewPrometheusSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &prometheusSyncer{
		scheme:   r,
		existing: &monitoringv1.Prometheus{},
		p:        p,
		key:      p.GetPrometheusKey(),
	}
}

// GetKey returns the prometheusSyncer key through which an existing object may be identified
func (s *prometheusSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *prometheusSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Prometheus object
func (s *prometheusSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*monitoringv1.Prometheus)
	out.Labels = GetPrometheusLabels(s.p)

	out.Spec = monitoringv1.PrometheusSpec{
		ScrapeInterval:     defaultScrapeInterval,
		EvaluationInterval: defaultEvaluationInterval,
		ServiceMonitorSelector: &metav1.LabelSelector{
			MatchLabels: s.p.GetProjectLabel(),
		},
		Version:   prometheusVersion,
		BaseImage: prometheusBaseImage,
	}

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *prometheusSyncer) GetErrorEventReason(err error) EventReason {
	if err != nil {
		return EventReasonPrometheusFailed
	}
	return EventReasonPrometheusUpdated
}
