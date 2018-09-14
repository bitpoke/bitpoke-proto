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
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/presslabs/controller-util/syncer"

	"github.com/presslabs/dashboard/pkg/controller/site/internal/sync"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

var _ = Describe("MemcachedStatefulSetSyncer", func() {
	var (
		wp        *wordpressv1alpha1.Wordpress
		memcached *appsv1.StatefulSet
		syncer    syncer.Interface
	)

	BeforeEach(func() {
		wp = &wordpressv1alpha1.Wordpress{}
		memcached = &appsv1.StatefulSet{}
		syncer = sync.NewMemcachedStatefulSetSyncer(wp)
	})

	DescribeTable("when Wordpress memcached.provisioner.presslabs.com/memory annotation",
		func(annotationValue, expectedValue string, shouldErr bool) {
			if len(annotationValue) > 0 {
				wp.ObjectMeta.Annotations = map[string]string{"memcached.provisioner.presslabs.com/memory": annotationValue}
			}
			err := syncer.SyncFn(memcached)
			if shouldErr {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				actual := memcached.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceMemory]
				Expect(actual).To(Equal(resource.MustParse(expectedValue)))
			}
		},
		Entry("is missing uses default value", "", sync.DefaultMemcachedMemory, false),
		Entry("is valid uses provided value", "10Gi", "10Gi", false),
		Entry("is not set uses default value", "invalid", "", true),
	)
})
