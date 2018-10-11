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
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/presslabs/controller-util/syncer"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/controller/project/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

var _ = Describe("The GiteaIngressSyncer transform func T", func() {
	var giteaIngress *extv1beta1.Ingress
	var proj *project.Project

	BeforeEach(func() {
		proj = project.New(&dashboardv1alpha1.Project{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("proj-%d", rand.Int31()),
				Namespace: fmt.Sprintf("org-%d", rand.Int31()),
			},
		})
		giteaIngress = &extv1beta1.Ingress{}

		giteaIngressSyncer := sync.NewGiteaIngressSyncer(proj, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)
		err := giteaIngressSyncer.SyncFn(giteaIngress)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("uses the right domain", func() {
		Expect(giteaIngress.Spec.Rules[0].Host).To(Equal(fmt.Sprintf("%s.git.presslabs.net", proj.Name)))
	})
})
