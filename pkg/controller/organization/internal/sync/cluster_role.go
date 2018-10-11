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
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

// NewOwnerClusterRoleSyncer returns a new syncer.Interface for reconciling owner ClusterRole
func NewOwnerClusterRoleSyncer(org *organization.Organization, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := org.ComponentObject(organization.OwnerClusterRole)

	return syncer.NewObjectSyncer("OwnerClusterRole", org.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*rbacv1.ClusterRole)

		out.Labels = labels.Merge(labels.Merge(out.Labels, org.Labels()), controllerLabels)

		out.Rules = []rbacv1.PolicyRule{
			{
				Resources: []string{
					"namespaces",
				},
				ResourceNames: []string{
					org.Name,
				},
				Verbs: []string{
					"delete",
				},
				APIGroups: []string{""},
			},
		}

		return nil
	})
}
