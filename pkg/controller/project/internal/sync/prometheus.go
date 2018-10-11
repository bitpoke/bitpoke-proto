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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

const (
	prometheusBaseImage       = "quay.io/prometheus/prometheus"
	prometheusVersion         = "v2.3.2"
	defaultScrapeInterval     = "10s"
	defaultEvaluationInterval = "30s"
)

// NewPrometheusSyncer returns a new syncer.Interface for reconciling Prometheus
func NewPrometheusSyncer(proj *project.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(project.Prometheus)

	obj := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(project.Prometheus),
			Namespace: proj.ComponentName(project.Namespace),
		},
	}

	return syncer.NewObjectSyncer("Prometheus", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*monitoringv1.Prometheus)
		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		out.Spec = monitoringv1.PrometheusSpec{
			ScrapeInterval:     defaultScrapeInterval,
			EvaluationInterval: defaultEvaluationInterval,
			Version:            prometheusVersion,
			BaseImage:          prometheusBaseImage,
		}

		return nil
	})
}
