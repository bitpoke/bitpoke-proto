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
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/presslabs/dashboard/pkg/internal/testutil/gomega"
)

const timeout = time.Second * 1

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

		// create new k8s client
		c, err = client.New(cfg, client.Options{})
		Expect(err).To(Succeed())

		recFn, requests = SetupTestReconcile(newReconciler(mgr))
		Expect(add(mgr, recFn)).To(Succeed())

		stop = StartTestManager(mgr)
	})

	AfterEach(func() {
		close(stop)
	})

	When("creating a new Project object", func() {
		var (
			expectedRequest  reconcile.Request
			project          *corev1.Namespace
			projectNamespace string
			projectName      string
			projectCreatedBy string
			componentsLabels map[string]map[string]string
		)

		entries := []TableEntry{
			Entry("reconciles limit range", "default", "presslabs-dashboard", &corev1.LimitRange{}),
			Entry("reconciles resource quota", "default", "presslabs-dashboard", &corev1.ResourceQuota{}),
			Entry("reconciles prometheus service account", "prometheus", "prometheus", &corev1.ServiceAccount{}),
			Entry("reconciles prometheus role binding", "prometheus", "prometheus%.0s", &rbacv1.RoleBinding{}),
			Entry("reconciles prometheus", "prometheus-instance", "prometheus%.0s", &monitoringv1.Prometheus{}),
			Entry("reconciles gitea deployment", "gitea-deployment", "gitea%.0s", &appsv1.Deployment{}),
			Entry("reconciles gitea service", "gitea", "gitea%.0s", &corev1.Service{}),
			Entry("reconciles gitea ingress", "gitea", "gitea%.0s", &extv1beta1.Ingress{}),
			Entry("reconciles gitea pvc", "gitea", "gitea%.0s", &corev1.PersistentVolumeClaim{}),
			Entry("reconciles gitea secret", "gitea", "gitea-conf%.0s", &corev1.Secret{}),
			Entry("reconciles member role binding", "member", "member", &rbacv1.RoleBinding{}),
			Entry("reconciles owner role binding", "owner", "owner", &rbacv1.RoleBinding{}),
			Entry("reconciles memcached service monitor", "prometheus", "memcached", &monitoringv1.ServiceMonitor{}),
			Entry("reconciles mysql service monitor", "prometheus", "mysql", &monitoringv1.ServiceMonitor{}),
			Entry("reconciles wordpress service monitor", "prometheus", "wordpress", &monitoringv1.ServiceMonitor{}),
		}

		BeforeEach(func() {
			projectName = fmt.Sprintf("awesome-%d", rand.Int31())
			projectNamespace = fmt.Sprintf("proj-%s", projectName)
			projectCreatedBy = fmt.Sprintf("Dorel %d", rand.Int31())

			expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: projectNamespace}}

			project = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: projectNamespace,
					Labels: map[string]string{
						"presslabs.com/organization": organizationName,
						"presslabs.com/project":      projectName,
						"presslabs.com/kind":         "project",
					},
					Annotations: map[string]string{
						"presslabs.com/created-by": projectCreatedBy,
					},
				},
			}

			componentsLabels = map[string]map[string]string{
				"default": {
					"presslabs.com/project":        projectName,
					"presslabs.com/organization":   organizationName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
				},
				"prometheus-instance": {
					"presslabs.com/project":        projectName,
					"presslabs.com/organization":   organizationName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
					"app.kubernetes.io/name":       "prometheus",
					"app.kubernetes.io/version":    "v2.3.2",
				},
				"prometheus": {
					"presslabs.com/project":        projectName,
					"presslabs.com/organization":   organizationName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
					"app.kubernetes.io/name":       "prometheus",
				},
				"gitea": {
					"presslabs.com/project":        projectName,
					"presslabs.com/organization":   organizationName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
					"app.kubernetes.io/name":       "gitea",
					"app.kubernetes.io/component":  "web",
				},
				"gitea-deployment": {
					"presslabs.com/project":        projectName,
					"presslabs.com/organization":   organizationName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
					"app.kubernetes.io/name":       "gitea",
					"app.kubernetes.io/component":  "web",
					"app.kubernetes.io/version":    "1.5.2",
				},
				"member": {
					"presslabs.com/project":        projectName,
					"presslabs.com/organization":   organizationName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
					"presslabs.com/kind":           "project-member-list",
				},
				"owner": {
					"presslabs.com/project":        projectName,
					"presslabs.com/organization":   organizationName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
					"presslabs.com/kind":           "project-owner-list",
				},
			}

			// Create the Project object and expect the Reconcile and Namespace to be created
			Expect(c.Create(context.TODO(), project)).To(Succeed())

			// Wait for initial reconciliation
			Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))
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
			Expect(c.Delete(context.TODO(), project)).To(Succeed())
			Eventually(func() corev1.Namespace {
				ns := corev1.Namespace{}
				c.Get(context.TODO(), client.ObjectKey{Name: projectNamespace}, &ns)
				return ns
			}).Should(BeInPhase(corev1.NamespaceTerminating))
		})

		DescribeTable("the reconciler", func(component string, nameFmt string, obj runtime.Object) {
			name := nameFmt
			if strings.Count(nameFmt, "%") > 0 {
				name = fmt.Sprintf(nameFmt, project.Name)
			}
			key := types.NamespacedName{
				Name:      name,
				Namespace: projectNamespace,
			}
			Eventually(func() error { return c.Get(context.TODO(), key, obj) }, timeout).Should(Succeed())

			metaObj := obj.(metav1.Object)
			Expect(metaObj.GetLabels()).To(Equal(componentsLabels[component]))
		}, entries...)
	})
})
