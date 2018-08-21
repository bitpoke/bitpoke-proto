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

package project

import (
	"context"
	"fmt"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/controller/project/sync"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("project-controller")

// Add creates a new Project Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileProject{
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder("project-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("project-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Project
	err = c.Watch(&source.Kind{Type: &dashboardv1alpha1.Project{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch the Namespace created by Project
	err = c.Watch(&source.Kind{Type: &corev1.Namespace{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &dashboardv1alpha1.Project{},
	})
	if err != nil {
		return err
	}

	// Watch the Prometheus created by Project
	err = c.Watch(&source.Kind{Type: &monitoringv1.Prometheus{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &dashboardv1alpha1.Project{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileProject{}

// ReconcileProject reconciles a Project object
type ReconcileProject struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

const (
	eventNormal  = "Normal"
	eventWarning = "Warning"
)

// Reconcile reads that state of the cluster for a Project object and makes changes based on the state read
// and what is in the Project.Spec
// +kubebuilder:rbac:groups=dashboard.presslabs.com,resources=projects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheuses,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileProject) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Project instance
	instance := &dashboardv1alpha1.Project{}
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

	syncers := []sync.Interface{
		sync.NewNamespaceSyncer(instance, r.scheme),
		sync.NewPrometheusSyncer(instance, r.scheme),
	}

	for _, s := range syncers {
		key := s.GetKey()
		existing := s.GetExistingObjectPlaceholder()

		var op controllerutil.OperationType
		op, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, key, existing, s.T)
		reason := string(s.GetErrorEventReason(err))

		log.Info(fmt.Sprintf("%T %s/%s %s", existing, key.Namespace, key.Name, op))

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
