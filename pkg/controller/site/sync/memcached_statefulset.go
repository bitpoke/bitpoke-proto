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
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/presslabs/controller-util/syncer"

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

// NewMemcachedStatefulSetSyncer returns a new syncer.Interface for reconciling Memcached StatefulSet
func NewMemcachedStatefulSetSyncer(wp *wordpressv1alpha1.Wordpress) syncer.Interface {
	obj := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(memcachedStatefulSetNameFmt, wp.Name),
			Namespace: wp.Namespace,
		},
	}

	return syncer.New("MemcachedStatefulSet", wp, obj, func(existing runtime.Object) error {
		out := existing.(*appsv1.StatefulSet)

		replicas := memcachedReplicas

		memcachedMemory, e := wp.ObjectMeta.Annotations["memcached.provisioner.presslabs.com/memory"]
		if !e {
			memcachedMemory = DefaultMemcachedMemory
		}

		resMemcachedMemory, err := resource.ParseQuantity(memcachedMemory)
		if err != nil {
			return err
		}
		intVal, ok := resMemcachedMemory.AsInt64()
		if !ok {
			return fmt.Errorf("Cannot convert %s into int64", memcachedMemory)
		}
		// make conversion: 12Mi (12 * 2^20)
		memcachedMemoryArg := intVal / 1024 / 1024

		out.ObjectMeta.Labels = dashboardv1alpha1.GetSiteLabels(wp, "memcached")

		out.Spec.ServiceName = fmt.Sprintf(memcachedServiceNameFmt, wp.ObjectMeta.Name)
		out.Spec.Replicas = &replicas
		out.Spec.Selector = metav1.SetAsLabelSelector(dashboardv1alpha1.GetMemcachedSelector(wp))
		out.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: dashboardv1alpha1.GetSiteLabels(wp, "memcached"),
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            fmt.Sprintf(memcachedStatefulSetNameFmt, wp.ObjectMeta.Name),
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

		return nil
	})
}
