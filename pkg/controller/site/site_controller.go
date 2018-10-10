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

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	"github.com/presslabs/dashboard/pkg/controller/site/internal/sync"
	mysqlv1alpha1 "github.com/presslabs/mysql-operator/pkg/apis/mysql/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

var log = logf.Log.WithName("site-controller")

// Add creates a new Site Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
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
	err = c.Watch(&source.Kind{Type: &wordpressv1alpha1.Wordpress{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	subresources := []runtime.Object{
		&corev1.Service{},
		&corev1.Secret{},
		&appsv1.StatefulSet{},
		&mysqlv1alpha1.MysqlCluster{},
		&monitoringv1.ServiceMonitor{},
	}

	for _, subresource := range subresources {
		err = c.Watch(&source.Kind{Type: subresource}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &wordpressv1alpha1.Wordpress{},
		})
		if err != nil {
			return err
		}
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

// Reconcile reads that state of the cluster for a Wordpress object and makes changes based on the state read
// and what is in the Wordpress.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=,resources=services;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wordpress.presslabs.org,resources=wordpresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.presslabs.org,resources=mysqlclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileSite) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Site instance
	wp := &wordpressv1alpha1.Wordpress{}

	err := r.Get(context.TODO(), request.NamespacedName, wp)
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
		sync.NewMemcachedStatefulSetSyncer(wp, r.Client, r.scheme),
		sync.NewMemcachedServiceSyncer(wp, r.Client, r.scheme),
		sync.NewMemcachedServiceMonitorSyncer(wp, r.Client, r.scheme),
		sync.NewWordpressSyncer(wp, r.Client, r.scheme),
		sync.NewWordpressServiceMonitorSyncer(wp, r.Client, r.scheme),
		sync.NewMysqlClusterSyncer(wp, r.Client, r.scheme),
		sync.NewMysqlClusterSecretSyncer(wp, r.Client, r.scheme),
		sync.NewMysqlServiceMonitorSyncer(wp, r.Client, r.scheme),
	}

	return reconcile.Result{}, r.sync(syncers)
}

func (r *ReconcileSite) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			log.Error(err, "unable to sync")
			return err
		}
	}
	return nil
}
