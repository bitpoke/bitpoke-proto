/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
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
	"github.com/presslabs/dashboard/pkg/controller/projectns/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

var _ = Describe("The MemberRoleSyncer transform func T", func() {
	var (
		proj              *projectns.ProjectNamespace
		memberRoleBinding *rbacv1.RoleBinding
		organizationName  string
		projectName       string
	)

	BeforeEach(func() {
		projRand := rand.Int31()
		organizationName = fmt.Sprintf("org-%d", projRand)
		projectName = fmt.Sprintf("proj-%d", projRand)

		proj = projectns.New(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: projectName,
				Labels: map[string]string{
					"presslabs.com/project":      projectName,
					"presslabs.com/organization": organizationName,
				},
			},
		})

		memberRoleBinding = &rbacv1.RoleBinding{}
		memberRoleBindingSyncer := sync.NewMemberRoleBindingSyncer(proj, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)
		err := memberRoleBindingSyncer.SyncFn(memberRoleBinding)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("reconciles the member RoleBinding", func() {
		expectedSubjects := []rbacv1.Subject{
			{
				Kind: "User",
				Name: proj.Annotations["created-by"],
			},
		}
		Expect(memberRoleBinding.Subjects).To(Equal(expectedSubjects))

		expectedLabels := map[string]string{
			"presslabs.com/project":        projectName,
			"presslabs.com/organization":   organizationName,
			"app.kubernetes.io/managed-by": "project-namespace-controller.dashboard.presslabs.com",
			"presslabs.com/kind":           "project-member-list",
		}
		Expect(memberRoleBinding.Labels).To(Equal(expectedLabels))
	})
})
