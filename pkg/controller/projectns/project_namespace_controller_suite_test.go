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
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	logf "github.com/presslabs/controller-util/log"
	"github.com/presslabs/dashboard/pkg/apis"
)

const (
	organizationDisplayName = "ACME Inc."
	organizationName        = "acme"
)

var cfg *rest.Config
var t *envtest.Environment

func TestProjectController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Project Namespace Controller Suite", []Reporter{envtest.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	var err error

	logf.SetLogger(logf.ZapLoggerTo(GinkgoWriter, true))

	t = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "config", "crds"),
			filepath.Join("..", "..", "..", "vendor/github.com/coreos/prometheus-operator/example/prometheus-operator-crd"),
			filepath.Join("..", "..", "..", "vendor/github.com/presslabs/wordpress-operator/config/crds"),
		},
	}
	apis.AddToScheme(scheme.Scheme)

	cfg, err = t.Start()
	Expect(err).NotTo(HaveOccurred())

	c, err := client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred())

	Expect(c.Create(context.TODO(), &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "dashboard.presslabs.com:project::member",
		},
	})).To(Succeed())

	Expect(c.Create(context.TODO(), &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "dashboard.presslabs.com:project::owner",
		},
	})).To(Succeed())

	Expect(c.Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: organizationName,
			Labels: map[string]string{
				"presslabs.com/organization": organizationName,
				"presslabs.com/kind":         "organization",
			},
			Annotations: map[string]string{
				"org.dashboard.presslabs.net/display-name": organizationDisplayName,
			},
		},
	})).To(Succeed())
})

var _ = AfterSuite(func() {
	t.Stop()
})

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner and
// writes the request to requests after Reconcile is finished.
func SetupTestReconcile(inner reconcile.Reconciler) (reconcile.Reconciler, chan reconcile.Request) {
	requests := make(chan reconcile.Request)
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)
		requests <- req
		return result, err
	})
	return fn, requests
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager) chan struct{} {
	stop := make(chan struct{})
	go func() {
		Expect(mgr.Start(stop)).To(Succeed())
	}()
	return stop
}
