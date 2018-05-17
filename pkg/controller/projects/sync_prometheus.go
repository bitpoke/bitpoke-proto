/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package projects

import (
	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"

	projectsv1 "github.com/presslabs/dashboard/pkg/apis/projects/v1alpha1"
)

func (c *Controller) syncPrometheus(proj *projectsv1.Project, ns *corev1.Namespace) (*monitoringv1.Prometheus, error) {
	pr := monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name: proj.ObjectMeta.Name,
			Labels: labels.Set{
				"dashboard.presslabs.com/project": proj.ObjectMeta.Name,
			},
		},
		Spec: monitoringv1.PrometheusSpec{
			ScrapeInterval:     "10s",
			EvaluationInterval: "30s",
		},
	}

	prom, err := c.PrometheusClient.MonitoringV1().Prometheuses(ns.ObjectMeta.Name).Create(&pr)

	if err == nil {
		glog.Infof("Created Prometheus for %s", proj.ObjectMeta.Name)

		return prom, err
	}

	if errors.IsAlreadyExists(err) {
		prom, err = c.PrometheusClient.MonitoringV1().Prometheuses(ns.ObjectMeta.Name).Update(&pr)
		if err == nil {
			glog.Infof("Updated Prometheus for %s", proj.ObjectMeta.Name)
		} else {
			glog.Errorf("Error while updating Prometheus for %s: %v", proj.ObjectMeta.Name, err)
		}

		return prom, err
	}

	glog.Errorf("Error while creating Prometheus for %s: %v", proj.ObjectMeta.Name, err)
	return prom, err
}
