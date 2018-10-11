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
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mysqlv1alpha1 "github.com/presslabs/mysql-operator/pkg/apis/mysql/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const timeout = time.Second * 2

var _ = Describe("Site controller", func() {
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
		close(stop)
	})

	When("creating a new Wordpress resource", func() {
		var (
			wp              *wordpressv1alpha1.Wordpress
			expectedRequest reconcile.Request
		)

		entries := []TableEntry{
			Entry("reconciles memcached statefulset", "%s-memcached", &appsv1.StatefulSet{}),
			Entry("reconciles memcached service", "%s-memcached", &corev1.Service{}),
			Entry("reconciles memcached service monitor", "%s-memcached", &monitoringv1.ServiceMonitor{}),
			Entry("reconciles mysql cluster", "%s", &mysqlv1alpha1.MysqlCluster{}),
			Entry("reconciles mysql service monitor", "%s-mysql", &monitoringv1.ServiceMonitor{}),
			Entry("reconciles mysql cluster secret", "%s-mysql", &corev1.Secret{}),
			Entry("reconciles wordpress service monitor", "%s-wordpress", &monitoringv1.ServiceMonitor{}),
		}

		BeforeEach(func() {
			name := fmt.Sprintf("wp-%d", rand.Int31())
			namespace := "default"

			wp = &wordpressv1alpha1.Wordpress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: wordpressv1alpha1.WordpressSpec{
					Runtime: "runtime-example",
					Domains: []wordpressv1alpha1.Domain{
						"domain.com",
					},
				},
			}

			expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}

			// create Wordpress resource
			Expect(c.Create(context.TODO(), wp)).To(Succeed())

			// Wait for initial reconciliation
			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))
			// Wait for a second reconciliation triggered by components being created
			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))
			// TODO: find out why sometimes we get extra reconciliation requests and remove this loop
			done := time.After(100 * time.Millisecond)
		drain:
			for {
				select {
				case <-requests:
					continue
				case <-done:
					break drain
				}
			}
			// We need to make sure that the controller does not create infinite loops
			Consistently(requests).ShouldNot(Receive(Equal(expectedRequest)))
		})

		AfterEach(func() {
			// cleanup Wordpress resource
			Expect(c.Delete(context.TODO(), wp)).To(Succeed())

			// GC created objects
			for _, e := range entries {
				obj := e.Parameters[1].(runtime.Object)
				nameFmt := e.Parameters[0].(string)
				mo := obj.(metav1.Object)
				mo.SetName(fmt.Sprintf(nameFmt, wp.Name))
				mo.SetNamespace(wp.Namespace)
				c.Delete(context.TODO(), obj)
			}
		})

		DescribeTable("the reconciler", func(nameFmt string, obj runtime.Object) {
			key := types.NamespacedName{
				Name:      fmt.Sprintf(nameFmt, wp.Name),
				Namespace: wp.Namespace,
			}
			Eventually(func() error { return c.Get(context.TODO(), key, obj) }, timeout).Should(Succeed())
		}, entries...)
	})
})
