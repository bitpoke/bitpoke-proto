/*
Copyright 2018 Pressinfra SRL

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
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

// NewPrometheusRoleBindingSyncer returns a new syncer.Interface for reconciling Prometheus ServiceAccount
func NewPrometheusRoleBindingSyncer(proj *projectns.ProjectNamespace, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(projectns.PrometheusRoleBinding)

	obj := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(projectns.PrometheusRoleBinding),
			Namespace: proj.Name,
		},
	}

	return syncer.NewObjectSyncer("PrometheusRoleBinding", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*rbacv1.RoleBinding)
		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		if out.CreationTimestamp.IsZero() {
			out.RoleRef = rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "dashboard.presslabs.com:project::prometheus",
			}
		}

		out.Subjects = []rbacv1.Subject{
			{
				Kind: rbacv1.ServiceAccountKind,
				Name: proj.ComponentName(projectns.PrometheusServiceAccount),
			},
		}

		return nil
	})
}
