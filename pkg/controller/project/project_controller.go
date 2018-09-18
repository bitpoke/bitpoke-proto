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

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
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

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/controller/project/internal/sync"
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

	subresources := []runtime.Object{
		&corev1.Namespace{},
		&corev1.ResourceQuota{},
		&corev1.Service{},
		&corev1.PersistentVolumeClaim{},
		&appsv1.Deployment{},
		&extv1beta1.Ingress{},
		&monitoringv1.Prometheus{},
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
// +kubebuilder:rbac:groups=,resources=services;persistentvolumeclaims;resourcequotas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dashboard.presslabs.com,resources=projects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheuses,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileProject) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Project instance
	project := &dashboardv1alpha1.Project{}
	err := r.Get(context.TODO(), request.NamespacedName, project)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	syncers := []syncer.Interface{
		sync.NewNamespaceSyncer(project, r.Client, r.scheme),
		sync.NewResourceQuotaSyncer(project, r.Client, r.scheme),
		sync.NewGiteaSecretSyncer(project, r.Client, r.scheme),
		sync.NewGiteaPVCSyncer(project, r.Client, r.scheme),
		sync.NewGiteaDeploymentSyncer(project, r.Client, r.scheme),
		sync.NewGiteaServiceSyncer(project, r.Client, r.scheme),
		sync.NewGiteaIngressSyncer(project, r.Client, r.scheme),
		sync.NewPrometheusSyncer(project, r.Client, r.scheme),
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
