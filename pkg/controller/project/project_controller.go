/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package project

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	logf "github.com/presslabs/controller-util/log"
	"github.com/presslabs/controller-util/syncer"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/controller/project/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/predicate"
	"github.com/presslabs/dashboard/pkg/internal/project"
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

	// Watch for changes to the Project
	err = c.Watch(
		&source.Kind{Type: &dashboardv1alpha1.Project{}},
		&handler.EnqueueRequestForObject{},
		predicate.NewKindPredicate("project"),
		predicate.ResourceNotDeleted,
	)
	if err != nil {
		return err
	}

	subresources := []runtime.Object{
		&corev1.Namespace{},
	}

	for _, subresource := range subresources {
		err = c.Watch(&source.Kind{Type: subresource}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &dashboardv1alpha1.Project{},
		})
		if err != nil {
			return err
		}
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

// Reconcile reads that state of the cluster for a Project object and makes changes based on the state read
// and what is in the Project.Spec
// +kubebuilder:rbac:groups=,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dashboard.presslabs.com,resources=projects,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileProject) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Project instance
	proj := project.New(&dashboardv1alpha1.Project{})
	err := r.Get(context.TODO(), request.NamespacedName, proj.Unwrap())
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if err := proj.ValidateMetadata(); err != nil {
		log.Info("skip reconcile for invalid project", "obj", proj.Unwrap(), "error", err)
		return reconcile.Result{}, nil
	}

	if !proj.DeletionTimestamp.IsZero() {
		log.Info("skip reconcile for deleted project", "obj", proj.Unwrap())
		return reconcile.Result{}, nil
	}

	syncers := []syncer.Interface{
		sync.NewNamespaceSyncer(proj, r.Client, r.scheme),
	}

	return reconcile.Result{}, r.sync(syncers)
}

func (r *ReconcileProject) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			log.Error(err, "unable to sync")
			return err
		}
	}
	return nil
}
