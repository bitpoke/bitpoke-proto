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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/presslabs/dashboard/pkg/internal/site"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

var _ = Describe("Organization webhook", func() {
	var (
		// stop channel for controller manager
		stop chan struct{}

		webhook          *siteValidation
		wp               *site.Site
		siteName         string
		projectName      string
		organizationName string
	)

	BeforeEach(func() {
		siteName = fmt.Sprintf("site%d", rand.Int31())
		projectName = fmt.Sprintf("project%d", rand.Int31())
		organizationName = fmt.Sprintf("organization%d", rand.Int31())

		wp = site.New(&wordpressv1alpha1.Wordpress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      siteName,
				Namespace: fmt.Sprintf("proj-%s", projectName),
			},
		})
		webhook = &siteValidation{}

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())

		webhook.InjectClient(mgr.GetClient())
		webhook.InjectDecoder(mgr.GetAdmissionDecoder())

		stop = StartTestManager(mgr)
	})

	AfterEach(func() {
		close(stop)
	})

	It("returns error when metadata is missing", func() {
		err := webhook.validateSiteFn(wp)

		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/organization\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/project\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/site\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required annotation \"presslabs.com/created-by\" is missing")))
	})
	It("doesn't return error when metadata is set", func() {
		wp.SetLabels(map[string]string{
			"presslabs.com/organization": organizationName,
			"presslabs.com/project":      projectName,
			"presslabs.com/site":         siteName,
		})

		wp.SetAnnotations(map[string]string{
			"presslabs.com/created-by": "Andi",
		})

		Expect(webhook.validateSiteFn(wp)).To(Succeed())
	})
})
