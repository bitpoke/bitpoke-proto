/*
Copyright 2018 Pressinfra SRL.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sync

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

// NewMemcachedServiceMonitorSyncer returns a new syncer.Interface for reconciling Memcached ServiceMonitor
func NewMemcachedServiceMonitorSyncer(proj *projectns.ProjectNamespace, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(projectns.MemcachedServiceMonitor)

	obj := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(projectns.MemcachedServiceMonitor),
			Namespace: proj.Name,
		},
	}

	return syncer.NewObjectSyncer("MemcachedServiceMonitor", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*monitoringv1.ServiceMonitor)

		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		out.Spec.Endpoints = []monitoringv1.Endpoint{
			{
				Port: "prometheus",
				// MetricRelabelConfigs: []*monitoringv1.RelabelConfig{
				// 	{
				// 		SourceLabels: []string{"__meta_kubernetes_service_label_app_kubernetes_io_instance"},
				// 		TargetLabel:  "site",
				// 	},
				// },
			},
		}

		out.Spec.Selector = metav1.LabelSelector{MatchLabels: map[string]string{
			"app.kubernetes.io/name": "memcached",
		}}

		return nil
	})
}
