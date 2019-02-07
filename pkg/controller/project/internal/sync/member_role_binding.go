/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

// NewMemberRoleBindingSyncer returns a new syncer.Interface for reconciling
// member RoleBinding
func NewMemberRoleBindingSyncer(proj *project.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(project.MemberRoleBinding),
			Namespace: proj.ComponentName(project.Namespace),
		},
	}

	return syncer.NewObjectSyncer("MemberRoleBinding", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*rbacv1.RoleBinding)

		out.Labels = labels.Merge(labels.Merge(out.Labels, proj.Labels()), controllerLabels)
		out.Labels["presslabs.com/kind"] = "project-member-list"

		// only add the project creator if the subjects list is empty, meaning it's just been created
		if out.ObjectMeta.CreationTimestamp.IsZero() {
			member := rbacv1.Subject{
				Kind: "User",
				Name: proj.Annotations["presslabs.com/created-by"],
			}
			out.Subjects = append(out.Subjects, member)
		}
		out.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "dashboard.presslabs.com:project::member",
		}

		return nil
	})
}
