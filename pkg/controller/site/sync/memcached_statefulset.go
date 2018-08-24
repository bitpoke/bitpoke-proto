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
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const (
	// DefaultMemcachedMemory is the default value of memory for Memcached StatefulSet
	DefaultMemcachedMemory            = "512Mi"
	memcachedCPU                      = "100m"
	memcachedReplicas           int32 = 1
	memcachedImage                    = "docker.io/library/memcached:1.5.9-alpine"
	memcachedImagePullPolicy          = "IfNotPresent"
	memcachedExporterPort             = 9150
	memcachedPort                     = 11211
	memcachedExporterImage            = "quay.io/prometheus/memcached-exporter:v0.4.1"
	memcachedStatefulSetNameFmt       = "%s-memcached"
	// MemcachedStatefulSetFailed is the event reason for a failed Memcached StatefulSet reconcile
	MemcachedStatefulSetFailed EventReason = "MemcachedFailed"
	// MemcachedStatefulSetUpdated is the event reason for a successful Memcached StatefulSet reconcile
	MemcachedStatefulSetUpdated EventReason = "MemcachedUpdated"
)

var (
	resMemcachedCPU = resource.MustParse(memcachedCPU)
)

// memcachedServiceSyncer defines the Syncer for Memcached Service
type memcachedStatefulSetSyncer struct {
	scheme   *runtime.Scheme
	wp       *wordpressv1alpha1.Wordpress
	key      types.NamespacedName
	existing *appsv1.StatefulSet
}

// NewMemcachedStatefulSetSyncer returns a new sync.Interface for reconciling Memcached StatefulSet
func NewMemcachedStatefulSetSyncer(wp *wordpressv1alpha1.Wordpress, r *runtime.Scheme) Interface {
	return &memcachedStatefulSetSyncer{
		scheme:   r,
		wp:       wp,
		existing: &appsv1.StatefulSet{},
		key: types.NamespacedName{
			Namespace: wp.Namespace,
			Name:      wp.Name,
		},
	}
}

// GetKey returns the memcachedStatefulSetSyncer key through which an existing object may be identified
func (s *memcachedStatefulSetSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *memcachedStatefulSetSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Memcached StatefulSet object
func (s *memcachedStatefulSetSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*appsv1.StatefulSet)

	replicas := int32(memcachedReplicas)

	memcachedMemory, e := s.wp.ObjectMeta.Annotations["memcached.provisioner.presslabs.com/memory"]
	if !e {
		memcachedMemory = DefaultMemcachedMemory
	}

	resMemcachedMemory, err := resource.ParseQuantity(memcachedMemory)
	if err != nil {
		return nil, err
	}
	intVal, ok := resMemcachedMemory.AsInt64()
	if !ok {
		return nil, fmt.Errorf("Cannot convert %s into int64", memcachedMemory)
	}
	// make conversion: 12Mi (12 * 2^20)
	memcachedMemoryArg := intVal / 1024 / 1024

	out.ObjectMeta = metav1.ObjectMeta{
		Name:      fmt.Sprintf(memcachedStatefulSetNameFmt, s.wp.ObjectMeta.Name),
		Labels:    dashboardv1alpha1.GetSiteLabels(s.wp, "memcached"),
		Namespace: s.wp.ObjectMeta.Namespace,
	}

	out.Spec.ServiceName = fmt.Sprintf(memcachedServiceNameFmt, s.wp.ObjectMeta.Name)
	out.Spec.Replicas = &replicas
	out.Spec.Selector = metav1.SetAsLabelSelector(dashboardv1alpha1.GetMemcachedSelector(s.wp))
	out.Spec.Template = corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: dashboardv1alpha1.GetSiteLabels(s.wp, "memcached"),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            fmt.Sprintf(memcachedStatefulSetNameFmt, s.wp.ObjectMeta.Name),
					Image:           memcachedImage,
					ImagePullPolicy: memcachedImagePullPolicy,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"memory": resMemcachedMemory,
							"cpu":    resMemcachedCPU,
						},
						Limits: corev1.ResourceList{
							"memory": resMemcachedMemory,
							"cpu":    resMemcachedCPU,
						},
					},
					Command: []string{"memcached"},
					Args:    []string{"-m", strconv.FormatInt(memcachedMemoryArg, 10)},
					Ports: []corev1.ContainerPort{
						{
							Name:          "memcached",
							ContainerPort: memcachedPort,
						},
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromString("memcached"),
							},
						},
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromString("memcached"),
							},
						},
					},
				},
				{
					Name:            fmt.Sprintf("memcached-exporter"),
					Image:           memcachedExporterImage,
					ImagePullPolicy: memcachedImagePullPolicy,
					Ports: []corev1.ContainerPort{
						{
							HostPort:      memcachedExporterPort,
							ContainerPort: memcachedExporterPort,
						},
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/",
								Port: intstr.FromInt(memcachedExporterPort),
							},
						},
					},
				},
			},
		},
	}

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *memcachedStatefulSetSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return MemcachedStatefulSetUpdated
	}
	return MemcachedStatefulSetFailed
}
