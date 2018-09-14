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

	"github.com/presslabs/controller-util/syncer"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

// NewWordpressSyncer returns a new syncer.Interface for reconciling Wordpress
func NewWordpressSyncer(wp *wordpressv1alpha1.Wordpress) syncer.Interface {
	obj := &wordpressv1alpha1.Wordpress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wp.Name,
			Namespace: wp.Namespace,
		},
	}

	return syncer.New("Wordpress", wp, obj, func(existing runtime.Object) error {
		out := existing.(*wordpressv1alpha1.Wordpress)

		out.Spec.Env = []corev1.EnvVar{
			{
				Name:  "MEMCACHED_DISCOVERY_SERVICE",
				Value: fmt.Sprintf("%s.%s", memcachedServiceName(wp), wp.Namespace),
			},
		}

		return nil
	})
}
