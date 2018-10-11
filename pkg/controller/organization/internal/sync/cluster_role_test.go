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

var _ = Describe("The OwnerClusterRoleSyncer transform func T", func() {
	var org *organization.Organization
	var ownerClusterRole *rbacv1.ClusterRole
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
		ownerClusterRole = &rbacv1.ClusterRole{}

		ownerClusterRoleSyncer := sync.NewOwnerClusterRoleSyncer(org, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)
		err := ownerClusterRoleSyncer.SyncFn(ownerClusterRole)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("reconciles the owner ClusterRole", func() {
		expectedRules := []rbacv1.PolicyRule{
			{
				Resources: []string{
					"namespaces",
				},
				ResourceNames: []string{
					org.Name,
				},
				Verbs: []string{
					"delete",
				},
				APIGroups: []string{""},
			},
		}
		Expect(ownerClusterRole.Rules).To(Equal(expectedRules))

		expectedLabels := map[string]string{
			"presslabs.com/organization":   organizationName,
			"app.kubernetes.io/managed-by": "organization-controller.dashboard.presslabs.com",
		}
		Expect(ownerClusterRole.Labels).To(Equal(expectedLabels))
	})
})
