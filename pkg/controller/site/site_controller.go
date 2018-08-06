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

package site

import (
	"context"
	"log"

	"github.com/presslabs/dashboard/pkg/controller/site/sync"
	wpapiv1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Site Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this dashboard.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSite{
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder("site-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("site-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Wordpress
	err = c.Watch(&source.Kind{Type: &wpapiv1.Wordpress{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by Site - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wpapiv1.Wordpress{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSite{}

// ReconcileSite reconciles a Wordpress object
type ReconcileSite struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

const (
	eventNormal  = "Normal"
	eventWarning = "Warning"
)

type syncer struct {
	objType runtime.Object
	sync    sync.Interface
}

// Reconcile reads that state of the cluster for a Wordpress object and makes changes based on the state read
// and what is in the Wordpress.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wordpress.presslabs.org,resources=wordpress,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileSite) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Site instance
	instance := &wpapiv1.Wordpress{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	syncers := []sync.Interface{}

	for _, s := range syncers {
		key := s.GetKey()
		existing := s.GetExistingObjectPlaceholder()

		op, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, key, existing, s.T)
		reason := string(s.GetErrorEventReason(err))

		log.Printf("%T %s/%s %s", existing, key.Namespace, key.Name, op)

		if err != nil {
			r.recorder.Eventf(instance, eventWarning, reason, "%T %s/%s failed syncing: %s", existing, key.Namespace, key.Name, err)
			return reconcile.Result{}, err
		}
		if op != controllerutil.OperationNoop {
			r.recorder.Eventf(instance, eventNormal, reason, "%T %s/%s %s successfully", existing, key.Namespace, key.Name, op)
		}
	}

	return reconcile.Result{}, err
}
