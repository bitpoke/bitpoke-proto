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
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/controller/organization/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

var _ = Describe("The MemberRoleSyncer transform func T", func() {
	var org *organization.Organization
	var memberRoleBinding *rbacv1.RoleBinding
	var organizationName string

	BeforeEach(func() {
		orgRand := rand.Int31()
		organizationName = fmt.Sprintf("acme-%d", orgRand)

		org = organization.New(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: organizationName,
				Labels: map[string]string{
					"presslabs.com/organization": organizationName,
					"presslabs.com/kind":         "organization",
				},
			},
		})
		memberRoleBinding = &rbacv1.RoleBinding{}

		memberRoleBindingSyncer := sync.NewMemberRoleBindingSyncer(org, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)
		err := memberRoleBindingSyncer.SyncFn(memberRoleBinding)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("reconciles the member RoleBinding", func() {
		expectedSubjects := []rbacv1.Subject{
			{
				Kind: "User",
				Name: org.Annotations["created-by"],
			},
		}
		Expect(memberRoleBinding.Subjects).To(Equal(expectedSubjects))

		expectedLabels := map[string]string{
			"presslabs.com/organization":   organizationName,
			"app.kubernetes.io/managed-by": "organization-controller.dashboard.presslabs.com",
		}
		Expect(memberRoleBinding.Labels).To(Equal(expectedLabels))
	})
})
