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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/presslabs/controller-util/syncer"

	"github.com/presslabs/dashboard/pkg/controller/site/internal/sync"
	mysqlv1alpha1 "github.com/presslabs/mysql-operator/pkg/apis/mysql/v1alpha1"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

var _ = Describe("MysqlClusterSyncer", func() {
	var (
		wp     *wordpressv1alpha1.Wordpress
		mysql  *mysqlv1alpha1.MysqlCluster
		syncer syncer.Interface
	)

	BeforeEach(func() {
		wp = &wordpressv1alpha1.Wordpress{}
		mysql = &mysqlv1alpha1.MysqlCluster{}
		syncer = sync.NewMysqlClusterSyncer(wp)
	})

	DescribeTable("when Wordpress mysql.provisioner.presslabs.com/storage annotation",
		func(annotationValue, expectedValue string, shouldErr bool) {
			if len(annotationValue) > 0 {
				wp.ObjectMeta.Annotations = map[string]string{"mysql.provisioner.presslabs.com/storage": annotationValue}
			}
			err := syncer.SyncFn(mysql)
			if shouldErr {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				actual := mysql.Spec.VolumeSpec.PersistentVolumeClaimSpec.Resources.Requests[corev1.ResourceStorage]
				Expect(actual).To(Equal(resource.MustParse(expectedValue)))
			}
		},
		Entry("is missing uses default value", "", sync.DefaultMysqlVolumeStorage, false),
		Entry("is valid uses provided value", "10Gi", "10Gi", false),
		Entry("is not set uses default value", "invalid", "", true),
	)

	DescribeTable("when Wordpress mysql.provisioner.presslabs.com/memory annotation",
		func(annotationValue, expectedValue string, shouldErr bool) {
			if len(annotationValue) > 0 {
				wp.ObjectMeta.Annotations = map[string]string{"mysql.provisioner.presslabs.com/memory": annotationValue}
			}
			err := syncer.SyncFn(mysql)
			if shouldErr {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				actual := mysql.Spec.PodSpec.Resources.Requests[corev1.ResourceMemory]
				Expect(actual).To(Equal(resource.MustParse(expectedValue)))
			}
		},
		Entry("is missing uses default value", "", sync.DefaultMysqlPodMemory, false),
		Entry("is valid uses provided value", "10Gi", "10Gi", false),
		Entry("is not set uses default value", "invalid", "", true),
	)

	DescribeTable("when Wordpress mysql.provisioner.presslabs.com/cpu annotation",
		func(annotationValue, expectedValue string, shouldErr bool) {
			if len(annotationValue) > 0 {
				wp.ObjectMeta.Annotations = map[string]string{"mysql.provisioner.presslabs.com/cpu": annotationValue}
			}
			err := syncer.SyncFn(mysql)
			if shouldErr {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				actual := mysql.Spec.PodSpec.Resources.Requests[corev1.ResourceCPU]
				Expect(actual).To(Equal(resource.MustParse(expectedValue)))
			}
		},
		Entry("is missing uses default value", "", sync.DefaultMysqlPodCPU, false),
		Entry("is valid uses provided value", "10Gi", "10Gi", false),
		Entry("is not set uses default value", "invalid", "", true),
	)
})
