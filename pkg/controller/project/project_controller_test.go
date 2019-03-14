/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package project

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const timeout = time.Second * 1

var _ = Describe("Project Namespace controller", func() {
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

	When("creating a new Project Namespace object", func() {
		var (
			expectedRequest  reconcile.Request
			project          *dashboardv1alpha1.Project
			projName         string
			projUserID       string
			orgName          string
			componentsLabels map[string]map[string]string
		)

		entries := []TableEntry{
			Entry("reconciles namespace when project has namespace reference", "proj-awesome", "namespace", &corev1.Namespace{}),
			Entry("reconciles namespace when project has namespace reference", "", "namespace", &corev1.Namespace{}),
		}

		BeforeEach(func() {
			projName = fmt.Sprintf("%d", rand.Int31())
			projUserID = fmt.Sprintf("user#%d", rand.Int31())
			orgName = fmt.Sprintf("%d", rand.Int31())

			expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      projName,
				Namespace: orgName,
			}}

			componentsLabels = map[string]map[string]string{
				"namespace": {
					"presslabs.com/project":        projName,
					"presslabs.com/organization":   orgName,
					"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
				},
			}
		})

		AfterEach(func() {
			Expect(c.Delete(context.TODO(), project)).To(Succeed())
			Eventually(func() codes.Code {
				p := dashboardv1alpha1.Project{}
				err := c.Get(context.TODO(), client.ObjectKey{
					Name:      projName,
					Namespace: orgName,
				}, &p)
				return status.Code(err)
			}).Should(Equal(codes.Unknown))
		})

		DescribeTable("the reconciler", func(projNamespace string, component string, obj runtime.Object) {
			project = &dashboardv1alpha1.Project{
				ObjectMeta: metav1.ObjectMeta{
					Name:      projName,
					Namespace: orgName,
					Labels: map[string]string{
						"presslabs.com/kind":         "project",
						"presslabs.com/project":      projName,
						"presslabs.com/organization": orgName,
					},
					Annotations: map[string]string{
						"presslabs.com/created-by": projUserID,
					},
				},
				Spec: dashboardv1alpha1.ProjectSpec{
					NamespaceName: projNamespace,
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

			// Reconciler will update the project with namespace reference if
			// it is empty. Get the updated project
			err := c.Get(context.TODO(), types.NamespacedName{Name: projName, Namespace: orgName}, project)
			Expect(err).To(Succeed())
			Expect(project.Spec.NamespaceName).Should(Not(BeEmpty()))

			key := types.NamespacedName{
				Name: project.Spec.NamespaceName,
			}
			Eventually(func() error { return c.Get(context.TODO(), key, obj) }, timeout).Should(Succeed())

			metaObj := obj.(metav1.Object)
			Expect(metaObj.GetLabels()).To(Equal(componentsLabels[component]))
		}, entries...)
	})
})
