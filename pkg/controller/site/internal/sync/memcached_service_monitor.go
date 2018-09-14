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
	"fmt"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/presslabs/controller-util/syncer"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

func memcachedServiceMonitorName(wp *wordpressv1alpha1.Wordpress) string {
	return fmt.Sprintf("%s-memcached", wp.Name)
}

// NewMemcachedServiceMonitorSyncer returns a new syncer.Interface for reconciling Memcached ServiceMonitor
func NewMemcachedServiceMonitorSyncer(wp *wordpressv1alpha1.Wordpress) syncer.Interface {
	obj := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcachedServiceMonitorName(wp),
			Namespace: wp.Namespace,
		},
	}

	return syncer.New("MemcachedServiceMonitor", wp, obj, func(existing runtime.Object) error {
		out := existing.(*monitoringv1.ServiceMonitor)

		out.ObjectMeta.Labels = getSiteLabels(wp, "memcached-service-monitor")

		out.Spec.Endpoints = []monitoringv1.Endpoint{
			{
				Port: "memcached",
			},
		}

		out.Spec.Selector = metav1.LabelSelector{
			MatchLabels: labels.Set{
				"app.kubernetes.io/app-instance": wp.Name,
				"app.kubernetes.io/component":    "memcached",
				"app.kubernetes.io/name":         "wordpress",
			},
		}

		return nil
	})
}
