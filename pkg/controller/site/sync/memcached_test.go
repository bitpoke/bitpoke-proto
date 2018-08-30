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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/presslabs/dashboard/pkg/controller/site/sync"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

var _ = Describe("MemcachedStatefulSetSyncer", func() {
	When("Wordpress has no memory annotation", func() {
		It("uses a default value", func() {
			wp := &wordpressv1alpha1.Wordpress{}
			statefulSet := &appsv1.StatefulSet{}
			syncer := sync.NewMemcachedStatefulSetSyncer(wp, rts)
			newStatefulSet, err := syncer.T(statefulSet)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newStatefulSet.(*appsv1.StatefulSet).Spec.Template.Spec.Containers[0].Resources.Requests["memory"]).To(Equal(resource.MustParse(sync.DefaultMemcachedMemory)))
		})
	})
	When("Wordpress has a valid annotation", func() {
		It("successfully sets allocated resources", func() {
			memcachedMemory := "64Mi"
			m := make(map[string]string)
			m["memcached.provisioner.presslabs.com/memory"] = memcachedMemory
			wp := &wordpressv1alpha1.Wordpress{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: m,
				},
			}
			statefulSet := &appsv1.StatefulSet{}
			syncer := sync.NewMemcachedStatefulSetSyncer(wp, rts)
			newStatefulSet, err := syncer.T(statefulSet)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newStatefulSet.(*appsv1.StatefulSet).Spec.Template.Spec.Containers[0].Resources.Requests["memory"]).To(Equal(resource.MustParse(memcachedMemory)))
		})
	})
	When("Wordpress has an invalid annotation", func() {
		It("returns error and doesn't change the statefullset", func() {
			m := make(map[string]string)
			m["memcached.provisioner.presslabs.com/memory"] = "invalid"
			wp := &wordpressv1alpha1.Wordpress{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: m,
				},
			}
			statefulSet := &appsv1.StatefulSet{}
			syncer := sync.NewMemcachedStatefulSetSyncer(wp, rts)
			_, err := syncer.T(statefulSet)
			Expect(err).Should(HaveOccurred())
		})
	})
})
