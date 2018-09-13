/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	prometheusName            = "prometheus"
	prometheusBaseImage       = "quay.io/prometheus/prometheus"
	prometheusVersion         = "v2.3.2"
	defaultScrapeInterval     = "10s"
	defaultEvaluationInterval = "30s"
)

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

// NewPrometheusSyncer returns a new syncer.Interface for reconciling Prometheus
func NewPrometheusSyncer(proj *dashboardv1alpha1.Project) syncer.Interface {
	obj := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.GetPrometheusName(),
			Namespace: proj.GetNamespaceName(),
		},
	}

	return syncer.New("Prometheus", proj, obj, func(existing runtime.Object) error {
		out := existing.(*monitoringv1.Prometheus)
		out.Labels = GetPrometheusLabels(proj)

		out.Spec = monitoringv1.PrometheusSpec{
			ScrapeInterval:     defaultScrapeInterval,
			EvaluationInterval: defaultEvaluationInterval,
			ServiceMonitorSelector: &metav1.LabelSelector{
				MatchLabels: proj.GetProjectLabel(),
			},
			Version:   prometheusVersion,
			BaseImage: prometheusBaseImage,
		}

		return nil
	})
}
