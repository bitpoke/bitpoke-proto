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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	mysqlv1alpha1 "github.com/presslabs/mysql-operator/pkg/apis/mysql/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const (
	mysqlClustereNameFmt = "%s-mysql"
	// MysqlClusterFailed is the event reason for a failed MysqlCluster reconcile
	MysqlClusterFailed EventReason = "MysqlClusterFailed"
	// MysqlClusterUpdated is the event reason for a successful MysqlCluster reconcile
	MysqlClusterUpdated EventReason = "MysqlClusterUpdated"
	// DefaultMysqlVolumeStorage is the default value of storage for MysqlCluster
	DefaultMysqlVolumeStorage = "8Gi"
	// DefaultMysqlPodMemory is the default value of memory for MysqlCluster
	DefaultMysqlPodMemory = "512Mi"
	// DefaultMysqlPodCPU is the default value of CPU for MysqlCluster
	DefaultMysqlPodCPU = "200m"
)

// mysqlClusterSyncer defines the Syncer for MysqlCluster
type mysqlClusterSyncer struct {
	scheme   *runtime.Scheme
	wp       *wordpressv1alpha1.Wordpress
	key      types.NamespacedName
	existing *mysqlv1alpha1.MysqlCluster
}

// NewMysqlClusterSyncer returns a new sync.Interface for reconciling Wordpress
func NewMysqlClusterSyncer(wp *wordpressv1alpha1.Wordpress, r *runtime.Scheme) Interface {
	return &mysqlClusterSyncer{
		scheme:   r,
		wp:       wp,
		existing: &mysqlv1alpha1.MysqlCluster{},
		key: types.NamespacedName{
			Namespace: wp.Namespace,
			Name:      wp.Name,
		},
	}
}

// GetInstance returns the mysqlClusterSyncer instance (mysqlClusterSyncer.mysql)
func (s *mysqlClusterSyncer) GetInstance() runtime.Object { return s.wp }

// GetKey returns the mysqlClusterSyncer key through which an existing object may be identified
func (s *mysqlClusterSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *mysqlClusterSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the MysqlCluster object
func (s *mysqlClusterSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*mysqlv1alpha1.MysqlCluster)

	volumeStorage, exists := s.wp.ObjectMeta.Annotations["mysql.provisioner.presslabs.com/storage"]
	if !exists {
		volumeStorage = DefaultMysqlVolumeStorage
	}
	resVolumeStorage, err := resource.ParseQuantity(volumeStorage)
	if err != nil {
		return nil, err
	}

	memory, exists := s.wp.ObjectMeta.Annotations["mysql.provisioner.presslabs.com/memory"]
	if !exists {
		memory = DefaultMysqlPodMemory
	}
	resPodMemory, err := resource.ParseQuantity(memory)
	if err != nil {
		return nil, err
	}

	cpu, exists := s.wp.ObjectMeta.Annotations["mysql.provisioner.presslabs.com/cpu"]
	if !exists {
		cpu = DefaultMysqlPodCPU
	}
	resPodCPU, err := resource.ParseQuantity(cpu)
	if err != nil {
		return nil, err
	}

	out.ObjectMeta = metav1.ObjectMeta{
		Name:      fmt.Sprintf(mysqlClustereNameFmt, s.wp.ObjectMeta.Name),
		Namespace: s.wp.ObjectMeta.Namespace,
	}

	out.Spec.PodSpec.Resources = corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceMemory: resPodMemory,
			corev1.ResourceCPU:    resPodCPU,
		},
	}

	out.Spec.VolumeSpec.PersistentVolumeClaimSpec = corev1.PersistentVolumeClaimSpec{
		Resources: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: resVolumeStorage},
		},
	}

	return out, nil
}

// GetErrorEventReason returns a reason for changes in the object state
func (s *mysqlClusterSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return MysqlClusterUpdated
	}
	return MysqlClusterFailed
}
