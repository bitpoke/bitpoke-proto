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

package apis

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

var (
	// PrometheusSchemeGroupVersion is group version used to register Prometheus related objects
	PrometheusSchemeGroupVersion = schema.GroupVersion{Group: "monitoring.coreos.com", Version: "v1"}

	// PrometheusSchemeBuilder is used to add go types to the GroupVersionKind scheme
	PrometheusSchemeBuilder = &scheme.Builder{GroupVersion: PrometheusSchemeGroupVersion}
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	PrometheusSchemeBuilder.Register(&monitoringv1.Prometheus{}, &monitoringv1.PrometheusList{}, &monitoringv1.ServiceMonitor{}, &monitoringv1.ServiceMonitorList{},
		&monitoringv1.Alertmanager{})
	AddToSchemes = append(AddToSchemes, PrometheusSchemeBuilder.AddToScheme)
}
