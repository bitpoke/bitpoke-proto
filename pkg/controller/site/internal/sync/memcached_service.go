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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

func memcachedServiceName(wp *wordpressv1alpha1.Wordpress) string {
	return fmt.Sprintf("%s-memcached", wp.Name)
}

// NewMemcachedServiceSyncer returns a new syncer.Interface for reconciling Memcached Service
func NewMemcachedServiceSyncer(wp *wordpressv1alpha1.Wordpress, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcachedServiceName(wp),
			Namespace: wp.Namespace,
		},
	}

	return syncer.NewObjectSyncer("MemcachedService", wp, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Service)

		siteLabels := getSiteLabels(wp, "memcached")
		out.ObjectMeta.Labels = labels.Merge(out.ObjectMeta.Labels, siteLabels)

		out.Spec.ClusterIP = "None"

		if !labels.Equals(siteLabels, out.Spec.Selector) {
			if out.ObjectMeta.CreationTimestamp.IsZero() {
				out.Spec.Selector = siteLabels
			} else {
				return fmt.Errorf("service selector is immutable")
			}
		}

		if len(out.Spec.Ports) != 1 {
			out.Spec.Ports = make([]corev1.ServicePort, 1)
		}
		out.Spec.Ports[0].Name = "memcached"
		out.Spec.Ports[0].Port = memcachedPort

		return nil
	})
}
