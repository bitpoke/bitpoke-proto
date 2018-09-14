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
	prometheusBaseImage       = "quay.io/prometheus/prometheus"
	prometheusVersion         = "v2.3.2"
	defaultScrapeInterval     = "10s"
	defaultEvaluationInterval = "30s"
)

// prometheusSelector returns a set of labels that can be used to identify Prometheus
// related resources
func prometheusSelector(project *dashboardv1alpha1.Project) labels.Set {
	prometheusLabels := labels.Set{
		"app.kubernetes.io/name": prometheusName(project),
	}
	return labels.Merge(getDefaultLabels(project), prometheusLabels)
}

// prometheusLabels returns a set of labels that should be applied on Prometheus
// related objects that are managed by the project controller
func prometheusLabels(project *dashboardv1alpha1.Project) labels.Set {
	prometheusLabels := labels.Set{
		"app.kubernetes.io/version": prometheusVersion,
	}
	return labels.Merge(prometheusSelector(project), prometheusLabels)
}

// prometheusName returns the name of the Prometheus resource
func prometheusName(project *dashboardv1alpha1.Project) string {
	return "prometheus"
}

// NewPrometheusSyncer returns a new syncer.Interface for reconciling Prometheus
func NewPrometheusSyncer(proj *dashboardv1alpha1.Project) syncer.Interface {
	obj := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusName(proj),
			Namespace: getNamespaceName(proj),
		},
	}

	return syncer.New("Prometheus", proj, obj, func(existing runtime.Object) error {
		out := existing.(*monitoringv1.Prometheus)
		out.Labels = prometheusLabels(proj)

		out.Spec = monitoringv1.PrometheusSpec{
			ScrapeInterval:     defaultScrapeInterval,
			EvaluationInterval: defaultEvaluationInterval,
			ServiceMonitorSelector: &metav1.LabelSelector{
				MatchLabels: getProjectLabel(proj),
			},
			Version:   prometheusVersion,
			BaseImage: prometheusBaseImage,
		}

		return nil
	})
}
