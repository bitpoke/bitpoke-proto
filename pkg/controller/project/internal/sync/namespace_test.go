/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"fmt"
	mathrand "math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/presslabs/controller-util/syncer"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

var _ = Describe("NamespaceSyncer", func() {
	var (
		p             *project.Project
		ns            *corev1.Namespace
		projName      string
		orgName       string
		userID        string
		projNamespace string
		displayName   string
	)

	BeforeEach(func() {
		orgName = fmt.Sprintf("%d", mathrand.Int31())
		projName = fmt.Sprintf("%d", mathrand.Int31())
		projNamespace = fmt.Sprintf("proj-%d", mathrand.Int31())
		userID = fmt.Sprintf("user#%d", mathrand.Int31())
		displayName = fmt.Sprintf("Awesome Project %s", projName)

		p = project.New(&dashboardv1alpha1.Project{
			ObjectMeta: metav1.ObjectMeta{
				Name:      projName,
				Namespace: orgName,
				Labels: map[string]string{
					"presslabs.com/organization": orgName,
					"presslabs.com/project":      projName,
				},
				Annotations: map[string]string{
					"presslabs.com/created-by":   userID,
					"presslabs.com/display-name": displayName,
				},
			},
			Spec: dashboardv1alpha1.ProjectSpec{
				NamespaceName: projNamespace,
			},
		})
		ns = &corev1.Namespace{}

		projSyncer := NewNamespaceSyncer(p, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)
		err := projSyncer.SyncFn(ns)
		Expect(err).To(Succeed())
	})

	It("reconciles the Namespace", func() {
		expectedLabels := map[string]string{
			"presslabs.com/kind":           "project",
			"app.kubernetes.io/managed-by": "project-controller.dashboard.presslabs.com",
			"presslabs.com/organization":   orgName,
			"presslabs.com/project":        projName,
		}
		Expect(ns.GetLabels()).To(Equal(expectedLabels))

		expectedAnnotations := map[string]string{
			"presslabs.com/created-by":   userID,
			"presslabs.com/display-name": displayName,
		}
		Expect(ns.Annotations).To(Equal(expectedAnnotations))
	})
})
