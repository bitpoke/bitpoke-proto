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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const (
	// MemcachedServiceMonitorFailed is the event reason for a failed ServiceMonitor reconcile
	MemcachedServiceMonitorFailed EventReason = "MemcachedServiceMonitorFailed"
	// MemcachedServiceMonitorUpdated is the event reason for a successful ServiceMonitor reconcile
	MemcachedServiceMonitorUpdated EventReason = "MemcachedServiceMonitorUpdated"
	memcachedServiceMonitorNameFmt             = "%s-memcached"
)

// memcachedServiceMonitorSyncer defines the Syncer for ServiceMonitor
type memcachedServiceMonitorSyncer struct {
	scheme   *runtime.Scheme
	wp       *wordpressv1alpha1.Wordpress
	key      types.NamespacedName
	existing *monitoringv1.ServiceMonitor
}

// NewMemcachedServiceMonitorSyncer returns a new sync.Interface for reconciling ServiceMonitor
func NewMemcachedServiceMonitorSyncer(wp *wordpressv1alpha1.Wordpress, r *runtime.Scheme) Interface {
	return &memcachedServiceMonitorSyncer{
		scheme:   r,
		wp:       wp,
		existing: &monitoringv1.ServiceMonitor{},
		key: types.NamespacedName{
			Name:      wp.Name,
			Namespace: wp.Namespace,
		},
	}
}

// GetInstance returns the memcachedServiceMonitorSyncer instance
func (s *memcachedServiceMonitorSyncer) GetInstance() runtime.Object { return s.wp }

// GetKey returns the memcachedServiceMonitorSyncer key through which an existing object may be identified
func (s *memcachedServiceMonitorSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *memcachedServiceMonitorSyncer) GetExistingObjectPlaceholder() runtime.Object {
	return s.existing
}

// T is the transform function used to reconcile the ServiceMonitor object
func (s *memcachedServiceMonitorSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*monitoringv1.ServiceMonitor)

	out.ObjectMeta = metav1.ObjectMeta{
		Name:      fmt.Sprintf(memcachedServiceMonitorNameFmt, s.wp.ObjectMeta.Name),
		Namespace: s.wp.ObjectMeta.Namespace,
		Labels:    dashboardv1alpha1.GetSiteLabels(s.wp, "memcached-service-monitor"),
	}

	out.Spec.Endpoints = []monitoringv1.Endpoint{
		{
			Port: "memcached",
		},
	}

	out.Spec.Selector = metav1.LabelSelector{
		MatchLabels: labels.Set{
			"app.kubernetes.io/app-instance": s.wp.Name,
			"app.kubernetes.io/component":    "memcached",
			"app.kubernetes.io/name":         "wordpress",
		},
	}

	err := controllerutil.SetControllerReference(s.wp, out, s.scheme)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *memcachedServiceMonitorSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return MemcachedServiceMonitorUpdated
	}
	return MemcachedServiceMonitorFailed
}
