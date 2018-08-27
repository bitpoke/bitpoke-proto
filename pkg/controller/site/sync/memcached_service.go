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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const (
	memcachedServiceNameFmt = "%s-memcached"
	// MemcachedServiceFailed is the event reason for a failed Memcached Service reconcile
	MemcachedServiceFailed EventReason = "MemcachedFailed"
	// MemcachedServiceUpdated is the event reason for a successful Memcached Service reconcile
	MemcachedServiceUpdated EventReason = "MemcachedUpdated"
)

// memcachedServiceSyncer defines the Syncer for Memcached Service
type memcachedServiceSyncer struct {
	scheme   *runtime.Scheme
	wp       *wordpressv1alpha1.Wordpress
	key      types.NamespacedName
	existing *corev1.Service
}

// NewMemcachedServiceSyncer returns a new sync.Interface for reconciling Memcached Service
func NewMemcachedServiceSyncer(wp *wordpressv1alpha1.Wordpress, r *runtime.Scheme) Interface {
	return &memcachedServiceSyncer{
		scheme:   r,
		wp:       wp,
		existing: &corev1.Service{},
		key: types.NamespacedName{
			Namespace: wp.Namespace,
			Name:      wp.Name,
		},
	}
}

// GetKey returns the memcachedServiceSyncer key through which an existing object may be identified
func (s *memcachedServiceSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *memcachedServiceSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Memcached Service object
func (s *memcachedServiceSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.Service)

	out.ObjectMeta = metav1.ObjectMeta{
		Name:      fmt.Sprintf(memcachedServiceNameFmt, s.wp.ObjectMeta.Name),
		Labels:    dashboardv1alpha1.GetSiteLabels(s.wp, "memcached"),
		Namespace: s.wp.ObjectMeta.Namespace,
	}

	out.Spec.ClusterIP = "None"
	out.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "memcached",
			Port:       int32(80),
			TargetPort: intstr.FromInt(memcachedPort),
		},
	}
	out.Spec.Selector = dashboardv1alpha1.GetMemcachedSelector(s.wp)

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *memcachedServiceSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return MemcachedServiceUpdated
	}
	return MemcachedServiceFailed
}
