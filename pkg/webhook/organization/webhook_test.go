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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	organization "github.com/presslabs/dashboard/pkg/internal/organization"
)

var _ = Describe("Organization webhook", func() {
	var (
		// stop channel for controller manager
		stop chan struct{}
		// controller k8s client
		c client.Client

		webhook *organizationValidation

		org *organization.Organization

		organizationName string
	)

	BeforeEach(func() {
		organizationName = fmt.Sprintf("organization%d", rand.Int31())

		org = organization.Wrap(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: organizationName,
			},
		})
		webhook = &organizationValidation{}

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
		err := webhook.validateOrganizationFn(org)

		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/organization\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/kind\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required annotation \"presslabs.com/created-by\" is missing")))
	})
	It("doesn't return error when metadata is set", func() {
		org.Namespace.SetLabels(map[string]string{
			"presslabs.com/organization": organizationName,
			"presslabs.com/kind":         "organization",
		})

		org.Namespace.SetAnnotations(map[string]string{
			"presslabs.com/created-by": "Andi",
		})

		Expect(webhook.validateOrganizationFn(org)).To(Succeed())
	})
})
