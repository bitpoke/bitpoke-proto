/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/gcloud/serviceaccount"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

var _ = Describe("GCloudServiceAccountSyncer", func() {
	var (
		p        *projectns.ProjectNamespace
		projName string
		orgName  string
		userID   string

		gcloudSASecret *corev1.Secret

		// k8s client
		cl client.Client
	)

	BeforeEach(func() {
		orgName = fmt.Sprintf("org-%d", rand.Int31())
		projName = fmt.Sprintf("proj-%d", rand.Int31())
		userID = fmt.Sprintf("user#%d", rand.Int31())

		cl = fake.NewFakeClient()

		p = projectns.New(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      projName,
				Namespace: orgName,
				Labels: map[string]string{
					"presslabs.com/organization": orgName,
					"presslabs.com/project":      projName,
				},
				Annotations: map[string]string{
					"presslabs.com/created-by": userID,
				},
			},
		})

		gcloudSASecret = &corev1.Secret{}
		gcloudSASyncer := NewGCloudServiceAccountSyncer(p, cl, scheme.Scheme).(*syncer.ObjectSyncer)
		err := gcloudSASyncer.SyncFn(gcloudSASecret)
		Expect(err).To(Succeed())
	})

	AfterEach(func() {
		Expect(serviceaccount.DeleteServiceAccount(string(gcloudSASecret.Data["service_account_email"]))).To(Succeed())
	})

	It("reconciles the GCloud Service Account Secret", func() {
		expectedLabels := map[string]string{
			"app.kubernetes.io/managed-by": "project-namespace-controller.dashboard.presslabs.com",
		}
		Expect(gcloudSASecret.GetLabels()).To(Equal(expectedLabels))
		Expect(gcloudSASecret.Data["google_application_credentials.json"]).NotTo(BeEmpty())
	})
})
