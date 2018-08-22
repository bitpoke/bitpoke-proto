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
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const timeout = time.Second * 5

var _ = Describe("Project controller", func() {
	var (
		// channel for incoming reconcile requests
		requests chan reconcile.Request
		// stop channel for controller manager
		stop chan struct{}
		// controller k8s client
		c client.Client
	)

	BeforeEach(func() {
		var recFn reconcile.Reconciler

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())
		c = mgr.GetClient()

		recFn, requests = SetupTestReconcile(newReconciler(mgr))
		Expect(add(mgr, recFn)).To(Succeed())

		stop = StartTestManager(mgr)
	})

	AfterEach(func() {
		time.Sleep(1 * time.Second)
		close(stop)
	})

	Describe("when creating a new Project object", func() {
		var expectedRequest reconcile.Request
		var instance *dashboardv1alpha1.Project

		BeforeEach(func() {
			name := fmt.Sprintf("proj-%d", rand.Int31())

			expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: name}}
			instance = &dashboardv1alpha1.Project{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "default",
				},
			}
		})

		It("reconciles the namespace", func() {
			// Create the Wordpress object and expect the Reconcile and Namespace to be created
			Expect(c.Create(context.TODO(), instance)).To(Succeed())
			defer c.Delete(context.TODO(), instance)

			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))

			ns := &corev1.Namespace{}
			Eventually(func() error { return c.Get(context.TODO(), instance.GetNamespaceKey(), ns) }, timeout).Should(Succeed())
			Expect(ns.Labels).To(Equal(map[string]string{
				"project.dashboard.presslabs.com/project": instance.Name,
				"app.kubernetes.io/deploy-manager":        "project-controller.dashboard.presslabs.com",
			}))
		})

		It("reconciles the Prometheus", func() {
			Expect(c.Create(context.TODO(), instance)).To(Succeed())
			defer c.Delete(context.TODO(), instance)

			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))

			prom := &monitoringv1.Prometheus{}
			Eventually(func() error { return c.Get(context.TODO(), instance.GetPrometheusKey(), prom) }, timeout).Should(Succeed())
			Expect(prom.Labels).To(Equal(map[string]string{
				"project.dashboard.presslabs.com/project": instance.Name,
				"app.kubernetes.io/deploy-manager":        "project-controller.dashboard.presslabs.com",
				"app.kubernetes.io/name":                  "prometheus",
				"app.kubernetes.io/version":               "v2.3.2",
			}))
		})

		It("reconciles the ResourceQuota", func() {
			Expect(c.Create(context.TODO(), instance)).To(Succeed())
			defer c.Delete(context.TODO(), instance)

			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))

			rq := &corev1.ResourceQuota{}
			Eventually(func() error { return c.Get(context.TODO(), instance.GetResourceQuotaKey(), rq) }, timeout).Should(Succeed())
		})
	})
})
