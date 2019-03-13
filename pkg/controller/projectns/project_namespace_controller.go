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

package projectns

import (
	"context"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
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
	"github.com/presslabs/dashboard/pkg/controller/projectns/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/predicate"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

var log = logf.Log.WithName("project-namespace-controller")

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
		recorder: mgr.GetRecorder("project-namespace-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("project-namespace-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to the Project Namespace
	err = c.Watch(
		&source.Kind{Type: &corev1.Namespace{}},
		&handler.EnqueueRequestForObject{},
		predicate.NewKindPredicate("project"),
		predicate.ResourceNotDeleted,
	)
	if err != nil {
		return err
	}

	subresources := []runtime.Object{
		&corev1.ResourceQuota{},
		&corev1.LimitRange{},
		&corev1.Service{},
		&corev1.PersistentVolumeClaim{},
		&appsv1.Deployment{},
		&extv1beta1.Ingress{},
		&rbacv1.RoleBinding{},
		&monitoringv1.ServiceMonitor{},
		&monitoringv1.Prometheus{},
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

var _ reconcile.Reconciler = &ReconcileProject{}

// ReconcileProject reconciles a Project Namespace object
type ReconcileProject struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a Project object and makes changes based on the state read
// and what is in the Project.Spec
// +kubebuilder:rbac:groups=,resources=services;persistentvolumeclaims;resourcequotas;namespaces;limitranges;events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dashboard.presslabs.com,resources=projects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheuses;servicemonitors,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileProject) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Project instance
	proj := projectns.New(&corev1.Namespace{})
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
		log.Info("skip reconcile for invalid project namespace", "obj", proj.Unwrap(), "error", err)
		return reconcile.Result{}, nil
	}

	if !proj.DeletionTimestamp.IsZero() {
		log.Info("skip reconcile for deleted project namespace", "obj", proj.Unwrap())
		return reconcile.Result{}, nil
	}

	syncers := []syncer.Interface{
		sync.NewLimitRangeSyncer(proj, r.Client, r.scheme),
		sync.NewResourceQuotaSyncer(proj, r.Client, r.scheme),
		sync.NewGiteaSecretSyncer(proj, r.Client, r.scheme),
		sync.NewGiteaPVCSyncer(proj, r.Client, r.scheme),
		sync.NewGiteaDeploymentSyncer(proj, r.Client, r.scheme),
		sync.NewGiteaServiceSyncer(proj, r.Client, r.scheme),
		sync.NewGiteaIngressSyncer(proj, r.Client, r.scheme),
		sync.NewPrometheusServiceAccountSyncer(proj, r.Client, r.scheme),
		sync.NewPrometheusRoleBindingSyncer(proj, r.Client, r.scheme),
		sync.NewPrometheusSyncer(proj, r.Client, r.scheme),
		sync.NewMemberRoleBindingSyncer(proj, r.Client, r.scheme),
		sync.NewOwnerRoleBindingSyncer(proj, r.Client, r.scheme),
		sync.NewMemcachedServiceMonitorSyncer(proj, r.Client, r.scheme),
		sync.NewMysqlServiceMonitorSyncer(proj, r.Client, r.scheme),
		sync.NewWordpressServiceMonitorSyncer(proj, r.Client, r.scheme),
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
