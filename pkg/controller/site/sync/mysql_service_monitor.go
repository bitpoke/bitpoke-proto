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

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const (
	mysqlServiceMonitorNameFmt = "%s-mysql"
)

// NewMysqlServiceMonitorSyncer returns a new syncer.Interface for reconciling Mysql ServiceMonitor
func NewMysqlServiceMonitorSyncer(wp *wordpressv1alpha1.Wordpress) syncer.Interface {
	obj := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(mysqlServiceMonitorNameFmt, wp.Name),
			Namespace: wp.Namespace,
		},
	}

	return syncer.New("MysqlServiceMonitor", wp, obj, func(existing runtime.Object) error {
		out := existing.(*monitoringv1.ServiceMonitor)

		out.ObjectMeta.Labels = dashboardv1alpha1.GetSiteLabels(wp, "mysql-service-monitor")

		out.Spec.Endpoints = []monitoringv1.Endpoint{
			{
				Port: "prometheus",
			},
		}

		out.Spec.Selector = metav1.LabelSelector{
			MatchLabels: labels.Set{
				"app.kubernetes.io/app-instance": wp.Name,
				"app.kubernetes.io/component":    "mysql",
				"app.kubernetes.io/name":         "wordpress",
			},
		}

		return nil
	})
}
