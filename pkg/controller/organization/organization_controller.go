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

package organization

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/controller/organization/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/organization"
	"github.com/presslabs/dashboard/pkg/internal/predicate"
)

var log = logf.Log.WithName("organization-controller")

// Add creates a new Organization Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOrganization{
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder("organization-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("organization-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Organization Namespace
	err = c.Watch(
		&source.Kind{Type: &corev1.Namespace{}},
		&handler.EnqueueRequestForObject{},
		predicate.NewKindPredicate("organization"),
	)
	if err != nil {
		return err
	}

	subresources := []runtime.Object{
		&rbacv1.ClusterRole{},
		&rbacv1.ClusterRoleBinding{},
		&rbacv1.RoleBinding{},
	}

	for _, subresource := range subresources {
		err = c.Watch(&source.Kind{Type: subresource}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &corev1.Namespace{},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileOrganization{}

// ReconcileOrganization reconciles a Organization object
type ReconcileOrganization struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a Organization object and makes changes based on the state read
// and what is in the Organization.Spec
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings;rolebindings,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileOrganization) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Organization instance
	org := organization.New(&corev1.Namespace{})
	err := r.Get(context.TODO(), request.NamespacedName, org.Unwrap())
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	syncers := []syncer.Interface{
		sync.NewOwnerClusterRoleSyncer(org, r.Client, r.scheme),
		sync.NewOwnerClusterRoleBindingSyncer(org, r.Client, r.scheme),
		sync.NewMemberRoleBindingSyncer(org, r.Client, r.scheme),
	}

	return reconcile.Result{}, r.sync(syncers)
}

func (r *ReconcileOrganization) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			log.Error(err, "unable to sync")
			return err
		}
	}
	return nil
}
