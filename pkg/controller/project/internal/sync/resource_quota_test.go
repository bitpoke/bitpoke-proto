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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/controller/project/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

var defaultQuotaValues = corev1.ResourceList{
	corev1.ResourceRequestsCPU:    resource.MustParse("4"),
	corev1.ResourceRequestsMemory: resource.MustParse("15Gi"),
	corev1.ResourceLimitsCPU:      resource.MustParse("8"),
	corev1.ResourceLimitsMemory:   resource.MustParse("32Gi"),
	corev1.ResourcePods:           resource.MustParse("20"),
}

var _ = Describe("The ResourceQuotaSyncer transform func T", func() {
	Context("finds no existing ResourceQuota", func() {
		proj := project.New(&corev1.Namespace{})
		rq := &corev1.ResourceQuota{}
		rqSyncer := sync.NewResourceQuotaSyncer(proj, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)

		It("uses a default value", func() {
			err := rqSyncer.SyncFn(rq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rq.Spec.Hard).To(Equal(defaultQuotaValues))
		})
	})

	Context("finds existing ResourceQuota with overridden values", func() {
		proj := project.New(&corev1.Namespace{})
		var rq *corev1.ResourceQuota

		biggerResourceRequestsCPU := defaultQuotaValues[corev1.ResourceRequestsMemory].DeepCopy()
		biggerResourceRequestsCPU.Add(resource.MustParse("1"))

		smallerResourcePods := defaultQuotaValues[corev1.ResourcePods].DeepCopy()
		smallerResourcePods.Sub(resource.MustParse("5"))

		BeforeEach(func() {
			rq = &corev1.ResourceQuota{}

			rq.Spec.Hard = corev1.ResourceList{
				corev1.ResourceRequestsCPU:    biggerResourceRequestsCPU,
				corev1.ResourceRequestsMemory: defaultQuotaValues[corev1.ResourceRequestsMemory],
				corev1.ResourceLimitsCPU:      defaultQuotaValues[corev1.ResourceLimitsCPU],
				corev1.ResourceLimitsMemory:   defaultQuotaValues[corev1.ResourceLimitsMemory],
				corev1.ResourcePods:           smallerResourcePods,
			}

			rqSyncer := sync.NewResourceQuotaSyncer(proj, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)
			err := rqSyncer.SyncFn(rq)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("uses the bigger overridden values", func() {
			Expect(rq.Spec.Hard[corev1.ResourceRequestsCPU]).To(Equal(biggerResourceRequestsCPU))
		})

		It("uses the default values instead of smaller (or equal) values", func() {
			Expect(rq.Spec.Hard[corev1.ResourceRequestsMemory]).To(Equal(defaultQuotaValues[corev1.ResourceRequestsMemory]))
			Expect(rq.Spec.Hard[corev1.ResourceLimitsCPU]).To(Equal(defaultQuotaValues[corev1.ResourceLimitsCPU]))
			Expect(rq.Spec.Hard[corev1.ResourceLimitsMemory]).To(Equal(defaultQuotaValues[corev1.ResourceLimitsMemory]))
			Expect(rq.Spec.Hard[corev1.ResourcePods]).To(Equal(defaultQuotaValues[corev1.ResourcePods]))
		})
	})
})
