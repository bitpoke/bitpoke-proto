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

// NewMemberRoleBindingSyncer returns a new syncer.Interface for reconciling member RoleBinding
func NewMemberRoleBindingSyncer(org *organization.Organization, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := org.ComponentObject(organization.MemberRoleBinding)

	return syncer.NewObjectSyncer("MemberRoleBinding", org.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*rbacv1.RoleBinding)

		out.Labels = labels.Merge(labels.Merge(out.Labels, org.Labels()), controllerLabels)

		// only add the organization creator if the subjects list is empty, meaning it's just been created
		if out.ObjectMeta.CreationTimestamp.IsZero() {
			member := rbacv1.Subject{
				Kind: "User",
				Name: org.Annotations["presslabs.com/created-by"],
			}
			out.Subjects = append(out.Subjects, member)
		}
		out.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "dashboard.presslabs.com:organization::member",
		}

		return nil
	})
}
