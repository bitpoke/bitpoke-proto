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
	"reflect"
	"strconv"

	"github.com/imdario/mergo"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/mergo/transformers"
	"github.com/presslabs/controller-util/syncer"
	site "github.com/presslabs/dashboard/pkg/internal/site"
)

const (
	// DefaultMemcachedMemory is the default value of memory for Memcached StatefulSet
	DefaultMemcachedMemory   = "512Mi"
	memcachedCPU             = "100m"
	memcachedImage           = "docker.io/library/memcached:1.5.9-alpine"
	memcachedImagePullPolicy = "IfNotPresent"
	memcachedExporterPort    = 9150
	memcachedPort            = 11211
	memcachedExporterImage   = "quay.io/prometheus/memcached-exporter:v0.4.1"
)

var (
	resMemcachedCPU         = resource.MustParse(memcachedCPU)
	memcachedReplicas int32 = 1
)

// NewMemcachedStatefulSetSyncer returns a new syncer.Interface for reconciling Memcached StatefulSet
func NewMemcachedStatefulSetSyncer(wp *site.Site, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := wp.ComponentLabels(site.MemcachedService)

	obj := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wp.ComponentName(site.MemcachedStatefulSet),
			Namespace: wp.Namespace,
		},
	}

	return syncer.NewObjectSyncer("MemcachedStatefulSet", wp.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*appsv1.StatefulSet)

		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		out.Spec.Replicas = &memcachedReplicas

		selector := metav1.SetAsLabelSelector(objLabels)
		if !reflect.DeepEqual(selector, out.Spec.Selector) {
			if out.ObjectMeta.CreationTimestamp.IsZero() {
				out.Spec.Selector = selector
			} else {
				return fmt.Errorf("statefullset selector is immutable")
			}
		}

		if out.Spec.ServiceName != wp.ComponentName(site.MemcachedService) {
			if out.ObjectMeta.CreationTimestamp.IsZero() {
				out.Spec.ServiceName = wp.ComponentName(site.MemcachedService)
			} else {
				return fmt.Errorf("statefullset service is immutable")
			}
		}

		out.Spec.Template.ObjectMeta.Labels = objLabels

		spec, err := getMemcachedPodSpec(wp)
		if err != nil {
			return err
		}

		err = mergo.Merge(&out.Spec.Template.Spec, spec, mergo.WithTransformers(transformers.PodSpec))
		if err != nil {
			return err
		}
		return nil
	})
}

func getMemcachedPodSpec(wp *site.Site) (corev1.PodSpec, error) {
	spec := corev1.PodSpec{}
	memcachedMemory, e := wp.ObjectMeta.Annotations["memcached.provisioner.presslabs.com/memory"]
	if !e {
		memcachedMemory = DefaultMemcachedMemory
	}

	resMemcachedMemory, err := resource.ParseQuantity(memcachedMemory)
	if err != nil {
		return spec, err
	}

	memcachedMemoryInt64, ok := resMemcachedMemory.AsInt64()
	if !ok {
		return spec, fmt.Errorf("Cannot convert %s into int64", memcachedMemory)
	}
	// make conversion: 12Mi (12 * 2^20)
	memcachedMemoryArg := memcachedMemoryInt64 / 1024 / 1024

	spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:            "memcached",
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
				Env: []corev1.EnvVar{
					{
						Name:  "MY_POD_NAMESPACE",
						Value: "default",
					},
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          "memcached",
						Protocol:      corev1.ProtocolTCP,
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
				Name:            "memcached-exporter",
				Image:           memcachedExporterImage,
				ImagePullPolicy: memcachedImagePullPolicy,
				Ports: []corev1.ContainerPort{
					{
						Name:          "prometheus",
						Protocol:      corev1.ProtocolTCP,
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
	}

	return spec, nil
}
