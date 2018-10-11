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
	"github.com/presslabs/dashboard/pkg/internal/site"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

// NewWordpressSyncer returns a new syncer.Interface for reconciling Wordpress
func NewWordpressSyncer(wp *site.Site, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &wordpressv1alpha1.Wordpress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wp.Name,
			Namespace: wp.Namespace,
		},
	}

	return syncer.NewObjectSyncer("Wordpress", wp.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*wordpressv1alpha1.Wordpress)

		out.Labels = labels.Merge(labels.Merge(out.Labels, wp.Labels()), controllerLabels)

		out.Spec.Env = []corev1.EnvVar{
			{
				Name:  "MEMCACHED_DISCOVERY_SERVICE",
				Value: fmt.Sprintf("%s.%s", wp.ComponentName(site.MemcachedService), wp.Namespace),
			},
			{
				Name:  "WORDPRESS_DB_HOST",
				Value: fmt.Sprintf("%s-mysql-master.%s", wp.Name, wp.Namespace),
			},
			{
				Name: "WORDPRESS_DB_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: wp.ComponentName(site.MysqlClusterSecret),
						},
						Key: "USER",
					},
				},
			},
			{
				Name: "WORDPRESS_DB_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: wp.ComponentName(site.MysqlClusterSecret),
						},
						Key: "PASSWORD",
					},
				},
			},
			{
				Name: "WORDPRESS_DB_NAME",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: wp.ComponentName(site.MysqlClusterSecret),
						},
						Key: "DATABASE",
					},
				},
			},
		}

		return nil
	})
}
