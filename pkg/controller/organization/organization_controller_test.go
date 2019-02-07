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

package organization

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/presslabs/dashboard/pkg/internal/testutil/gomega"
)

const timeout = time.Second * 1

var _ = Describe("Organization controller", func() {
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

		// create new k8s client
		c, err = client.New(cfg, client.Options{})
		Expect(err).To(Succeed())

		recFn, requests = SetupTestReconcile(newReconciler(mgr))
		Expect(add(mgr, recFn)).To(Succeed())

		stop = StartTestManager(mgr)
	})

	AfterEach(func() {
		time.Sleep(1 * time.Second)
		close(stop)
	})

	When("creating a new Organization object", func() {
		var (
			expectedRequest reconcile.Request
			org             *corev1.Namespace
			orgName         string
			orgNameLabel    string
			orgDisplayName  string
			orgCreatedBy    string
		)

		BeforeEach(func() {
			orgRand := rand.Int31()
			orgNameLabel = fmt.Sprintf("acme-%d", orgRand)
			orgName = fmt.Sprintf("org-%s", orgNameLabel)
			orgDisplayName = fmt.Sprintf("ACME %d Inc.", orgRand)
			orgCreatedBy = fmt.Sprintf("Dorel %d", rand.Int31())

			expectedRequest = reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: orgName,
				},
			}

			org = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: orgName,
					Labels: map[string]string{
						"presslabs.com/organization": orgNameLabel,
						"presslabs.com/kind":         "organization",
					},
					Annotations: map[string]string{
						"presslabs.com/display-name": orgDisplayName,
						"presslabs.com/created-by":   orgCreatedBy,
					},
				},
			}
			// Create the Organization in which the Project will live
			Expect(c.Create(context.TODO(), org)).To(Succeed())

			// Wait for initial reconciliation
			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))
			// TODO: redesign projects since you cannot have owner in another
			// namespace
			// Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))
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
			Expect(c.Delete(context.TODO(), org)).To(Succeed())
			Eventually(func() corev1.Namespace {
				ns := corev1.Namespace{}
				c.Get(context.TODO(), client.ObjectKey{Name: orgName}, &ns)
				return ns
			}).Should(BeInPhase(corev1.NamespaceTerminating))
		})

		It("reconciles the owner cluster role", func() {
			cr := &rbacv1.ClusterRole{}
			Eventually(func() error {
				return c.Get(
					context.TODO(),
					types.NamespacedName{
						Name: fmt.Sprintf("dashboard.presslabs.com:organization:%s:owner", orgNameLabel),
					},
					cr)
			}, timeout).Should(Succeed())

			Expect(cr.Labels).To(Equal(map[string]string{
				"presslabs.com/organization":   orgNameLabel,
				"app.kubernetes.io/managed-by": "organization-controller.dashboard.presslabs.com",
			}))
			Expect(cr.Rules).To(Equal([]rbacv1.PolicyRule{
				{
					Verbs:         []string{"delete", "update"},
					APIGroups:     []string{""},
					Resources:     []string{"namespaces"},
					ResourceNames: []string{orgName},
				},
			}))
		})
		It("reconciles the owners cluster role binding", func() {
			crb := &rbacv1.ClusterRoleBinding{}
			Eventually(func() error {
				return c.Get(
					context.TODO(),
					types.NamespacedName{
						Name: fmt.Sprintf("dashboard.presslabs.com:organization:%s:owners", orgNameLabel),
					},
					crb)
			}, timeout).Should(Succeed())
			Expect(crb.Labels).To(Equal(map[string]string{
				"presslabs.com/organization":   orgNameLabel,
				"app.kubernetes.io/managed-by": "organization-controller.dashboard.presslabs.com",
			}))
			Expect(crb.RoleRef).To(Equal(rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     fmt.Sprintf("dashboard.presslabs.com:organization:%s:owner", orgNameLabel),
			}))
		})
		It("reconciles the members role binding", func() {
			r := &rbacv1.RoleBinding{}
			Eventually(func() error {
				return c.Get(
					context.TODO(),
					types.NamespacedName{
						Name:      "members",
						Namespace: orgName,
					},
					r)
			}, timeout).Should(Succeed())
			Expect(r.Labels).To(Equal(map[string]string{
				"presslabs.com/organization":   orgNameLabel,
				"app.kubernetes.io/managed-by": "organization-controller.dashboard.presslabs.com",
			}))
			Expect(r.RoleRef).To(Equal(rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "dashboard.presslabs.com:organization::member",
			}))
		})
	})
})
