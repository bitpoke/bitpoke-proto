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
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const timeout = time.Second * 5

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
		time.Sleep(1 * time.Second)
		close(stop)
	})

	Describe("when creating a new Wordpress resource", func() {
		var expectedRequest reconcile.Request
		var wp *wordpressv1alpha1.Wordpress
		var ssKey types.NamespacedName
		var wpKey types.NamespacedName

		BeforeEach(func() {
			name := fmt.Sprintf("wp-%d", rand.Int31())
			namespace := "default"

			expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
			wp = &wordpressv1alpha1.Wordpress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: wordpressv1alpha1.WordpressSpec{
					Runtime: "runtime-example",
					Domains: []wordpressv1alpha1.Domain{
						"domain-example",
					},
				},
			}
			ssKey = types.NamespacedName{
				Name:      fmt.Sprintf("%s-memcached", name),
				Namespace: namespace,
			}
			wpKey = types.NamespacedName{
				Name:      name,
				Namespace: namespace,
			}
		})

		It("reconciles the memcached statefulset", func() {
			// Create the Wordpress object and expect the Reconcile and StatefulSet to be created
			Expect(c.Create(context.TODO(), wp)).To(Succeed())
			defer c.Delete(context.TODO(), wp)

			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))

			statefulSet := &appsv1.StatefulSet{}
			Eventually(func() error { return c.Get(context.TODO(), ssKey, statefulSet) }, timeout).Should(Succeed())
		})

		It("reconciles the memcached service", func() {
			// Create the Wordpress object and expect the Reconcile and Service to be created
			Expect(c.Create(context.TODO(), wp)).To(Succeed())
			defer c.Delete(context.TODO(), wp)

			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))

			service := &corev1.Service{}
			Eventually(func() error { return c.Get(context.TODO(), ssKey, service) }, timeout).Should(Succeed())
		})

		It("reconciles the wordpress resource", func() {
			// Create the Wordpress object and expcet the Reconcile to be created
			Expect(c.Create(context.TODO(), wp)).To(Succeed())
			defer c.Delete(context.TODO(), wp)

			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))

			wordpress := &wordpressv1alpha1.Wordpress{}
			Eventually(func() error { return c.Get(context.TODO(), wpKey, wordpress) }, timeout).Should(Succeed())
		})
	})
})
