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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const (
	// WordpressFailed is the event reason for a failed Wordpress reconcile
	WordpressFailed EventReason = "WordpressFailed"
	// WordpressUpdated is the event reason for a successful Wordpress reconcile
	WordpressUpdated EventReason = "WordpressUpdated"
)

// wordpressSyncer defines the Syncer for Wordpress
type wordpressSyncer struct {
	scheme   *runtime.Scheme
	wp       *wordpressv1alpha1.Wordpress
	key      types.NamespacedName
	existing *wordpressv1alpha1.Wordpress
}

// NewWordpressSyncer returns a new sync.Interface for reconciling Wordpress
func NewWordpressSyncer(wp *wordpressv1alpha1.Wordpress, r *runtime.Scheme) Interface {
	return &wordpressSyncer{
		scheme:   r,
		wp:       wp,
		existing: &wordpressv1alpha1.Wordpress{},
		key: types.NamespacedName{
			Namespace: wp.Namespace,
			Name:      wp.Name,
		},
	}
}

// GetInstance returns the wordpressSyncer instance (wordpressSyncer.wp)
func (s *wordpressSyncer) GetInstance() runtime.Object { return s.wp }

// GetKey returns the wordpressSyncer key through which an existing object may be identified
func (s *wordpressSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *wordpressSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Wordpress object
func (s *wordpressSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*wordpressv1alpha1.Wordpress)

	out.Spec.Env = []corev1.EnvVar{
		{
			Name:  "MEMCACHED_DISCOVERY_SERVICE",
			Value: fmt.Sprintf("%s-memcached.%s", s.wp.ObjectMeta.Name, s.wp.ObjectMeta.Namespace),
		},
	}

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *wordpressSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return WordpressUpdated
	}
	return WordpressFailed
}
