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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/site"
	mysqlv1alpha1 "github.com/presslabs/mysql-operator/pkg/apis/mysql/v1alpha1"
)

const (
	// DefaultMysqlVolumeStorage is the default value of storage for MysqlCluster
	DefaultMysqlVolumeStorage = "8Gi"
	// DefaultMysqlPodMemory is the default value of memory for MysqlCluster
	DefaultMysqlPodMemory = "512Mi"
	// DefaultMysqlPodCPU is the default value of CPU for MysqlCluster
	DefaultMysqlPodCPU = "200m"
)

// NewMysqlClusterSyncer returns a new syncer.Interface for reconciling MysqlCluster
func NewMysqlClusterSyncer(wp *site.Site, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := wp.ComponentLabels(site.MysqlCluster)

	obj := &mysqlv1alpha1.MysqlCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wp.ComponentName(site.MysqlCluster),
			Namespace: wp.Namespace,
		},
	}

	return syncer.NewObjectSyncer("MysqlCluster", wp.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*mysqlv1alpha1.MysqlCluster)

		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		out.Spec.SecretName = wp.ComponentName(site.MysqlClusterSecret)

		volumeStorage, exists := wp.ObjectMeta.Annotations["mysql.provisioner.presslabs.com/storage"]
		if !exists {
			volumeStorage = DefaultMysqlVolumeStorage
		}

		resVolumeStorage, err := resource.ParseQuantity(volumeStorage)
		if err != nil {
			return err
		}

		memory, exists := wp.ObjectMeta.Annotations["mysql.provisioner.presslabs.com/memory"]
		if !exists {
			memory = DefaultMysqlPodMemory
		}
		resPodMemory, err := resource.ParseQuantity(memory)
		if err != nil {
			return err
		}

		cpu, exists := wp.ObjectMeta.Annotations["mysql.provisioner.presslabs.com/cpu"]
		if !exists {
			cpu = DefaultMysqlPodCPU
		}
		resPodCPU, err := resource.ParseQuantity(cpu)
		if err != nil {
			return err
		}

		out.Spec.PodSpec.Resources = corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceMemory: resPodMemory,
				corev1.ResourceCPU:    resPodCPU,
			},
		}

		if len(out.Spec.VolumeSpec.PersistentVolumeClaimSpec.Resources.Requests) == 0 {
			out.Spec.VolumeSpec.PersistentVolumeClaimSpec.Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
		}
		out.Spec.VolumeSpec.PersistentVolumeClaimSpec.Resources.Requests[corev1.ResourceStorage] = resVolumeStorage

		return nil
	})
}
