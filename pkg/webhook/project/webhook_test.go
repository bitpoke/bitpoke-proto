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
	"context"
	"fmt"
	"math/rand"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

var _ = Describe("Project webhook", func() {
	var (
		// stop channel for controller manager
		stop chan struct{}
		// controller k8s client
		c client.Client

		webhook *projectValidation

		project *dashboardv1alpha1.Project

		organizationName string
	)

	BeforeEach(func() {
		organizationName = fmt.Sprintf("organization%d", rand.Int31())
		projectName := fmt.Sprintf("project%d", rand.Int31())

		project = &dashboardv1alpha1.Project{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: fmt.Sprintf("org-%s", organizationName),
				Name:      projectName,
			},
		}
		webhook = &projectValidation{}

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())
		c = mgr.GetClient()
		webhook.InjectClient(c)
		webhook.InjectDecoder(mgr.GetAdmissionDecoder())

		stop = StartTestManager(mgr)
	})

	AfterEach(func() {
		close(stop)
	})

	It("returns error when metadata is missing", func() {
		err := webhook.validateProjectFn(context.TODO(), project)

		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/organization\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/project\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required annotation \"presslabs.com/created-by\" is missing")))
	})
	It("creates project namespace", func() {
		project.SetLabels(map[string]string{
			"presslabs.com/organization": organizationName,
			"presslabs.com/project":      project.Name,
		})

		project.SetAnnotations(map[string]string{
			"presslabs.com/created-by": "Andi",
		})

		Expect(webhook.validateProjectFn(context.TODO(), project)).To(Succeed())

		Expect(c.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("proj-%s", project.Name)}, &corev1.Namespace{})).To(Succeed())
	})
	It("returns error when project namespace is taken", func() {
		Expect(c.Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("proj-%s", project.Name),
			},
		})).To(Succeed())
		project.SetLabels(map[string]string{
			"presslabs.com/organization": organizationName,
			"presslabs.com/project":      project.Name,
		})

		project.SetAnnotations(map[string]string{
			"presslabs.com/created-by": "Andi",
		})

		Expect(webhook.validateProjectFn(context.TODO(), project)).To(MatchError(NewStatusError(http.StatusBadRequest, fmt.Errorf("project \"%s\" is not available", project.Name))))
	})
})
