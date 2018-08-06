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

package sync_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/presslabs/dashboard/pkg/controller/site/sync"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

func TestSiteWordpress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Site Memcached Suite", []Reporter{envtest.NewlineReporter{}})
}

var _ = Describe("WordpressSyncer", func() {
	BeforeEach(func() {
		appsv1.AddToScheme(rts)
	})
	When("Wordpress has no MEMCACHED_DISCOVERY_SERVICE envvar", func() {
		It("successfully sets an envvar named MEMCACHED_DISCOVERY_SERVICE", func() {
			wp := &wordpressv1alpha1.Wordpress{}
			wpRes := &wordpressv1alpha1.Wordpress{}
			syncer := sync.NewWordpressSyncer(wp, rts)
			newWpRes, err := syncer.T(wpRes)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newWpRes.(*wordpressv1alpha1.Wordpress).Spec.Env[0].Name).To(Equal("MEMCACHED_DISCOVERY_SERVICE"))
			Expect(newWpRes.(*wordpressv1alpha1.Wordpress).Spec.Env[0].Value).To(Equal(fmt.Sprintf("%s-memcached.%s", wp.ObjectMeta.Name, wp.ObjectMeta.Namespace)))
		})
	})
})
