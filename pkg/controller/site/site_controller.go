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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/presslabs/dashboard/pkg/controller/site/sync"
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

	// Watch for changes to MysqlCluster
	err = c.Watch(&source.Kind{Type: &mysqlv1alpha1.MysqlCluster{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wordpressv1alpha1.Wordpress{},
	})
	if err != nil {
		return err
	}

	// Watch the Memcached StatefulSet created by Site
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wordpressv1alpha1.Wordpress{},
	})
	if err != nil {
		return err
	}

	// Wathc the serviceMonitor created by Site
	err = c.Watch(&source.Kind{Type: &monitoringv1.ServiceMonitor{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wordpressv1alpha1.Wordpress{},
	})
	if err != nil {
		return err
	}

	// Watch the Memcached Service created by Site
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wordpressv1alpha1.Wordpress{},
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

// Reconcile reads that state of the cluster for a Wordpress object and makes changes based on the state read
// and what is in the Wordpress.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wordpress.presslabs.org,resources=wordpress,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.presslabs.org,resources=mysqlcluster,verbs=get;list;watch;create;update;patch;delete
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

	syncers := []sync.Interface{
		sync.NewMemcachedStatefulSetSyncer(wp, r.scheme),
		sync.NewMemcachedServiceSyncer(wp, r.scheme),
		sync.NewMemcachedServiceMonitorSyncer(wp, r.scheme),
		sync.NewWordpressSyncer(wp, r.scheme),
		sync.NewMysqlClusterSyncer(wp, r.scheme),
		sync.NewMysqlServiceMonitorSyncer(wp, r.scheme),
	}

	for _, s := range syncers {
		key := s.GetKey()
		existing := s.GetExistingObjectPlaceholder()

		var op controllerutil.OperationType
		op, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, key, existing, s.T)
		reason := string(s.GetErrorEventReason(err))

		log.Info(string(op), "key", key, "kind", existing.GetObjectKind().GroupVersionKind().Kind)

		if err != nil {
			r.recorder.Eventf(s.GetInstance(), eventWarning, reason, "%T %s/%s failed syncing: %s", existing, key.Namespace, key.Name, err)
			return reconcile.Result{}, err
		}
		if op != controllerutil.OperationNoop {
			r.recorder.Eventf(s.GetInstance(), eventNormal, reason, "%T %s/%s %s successfully", existing, key.Namespace, key.Name, op)
		}
	}

	return reconcile.Result{}, err
}
