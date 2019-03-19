/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"context"
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
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

var _ = Describe("SMTPSecretSyncer", func() {
	var (
		p        *projectns.ProjectNamespace
		projName string
		orgName  string
		userID   string

		defSMTPSecret *corev1.Secret
		smtpSecret    *corev1.Secret
		defSMTPHost   []byte
		defSMTPPort   []byte
		defSMTPTLS    []byte

		// k8s client
		cl client.Client
	)

	BeforeEach(func() {
		orgName = fmt.Sprintf("org-%d", rand.Int31())
		projName = fmt.Sprintf("proj-%d", rand.Int31())
		userID = fmt.Sprintf("user#%d", rand.Int31())

		defSMTPHost = []byte("localhost")
		defSMTPPort = []byte(string(587))
		defSMTPTLS = []byte(string("yes"))

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

		defSMTPSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      options.SMTPSecret,
				Namespace: "default",
			},
			Data: map[string][]byte{
				"SMTP_HOST": defSMTPHost,
				"SMTP_PORT": defSMTPPort,
				"SMTP_TLS":  defSMTPTLS,
			},
		}
		Expect(cl.Create(context.TODO(), defSMTPSecret)).To(Succeed())

		smtpSecret = &corev1.Secret{}
		smtpSecretSyncer := NewSMTPSecretSyncer(p, cl, scheme.Scheme).(*syncer.ObjectSyncer)
		err := smtpSecretSyncer.SyncFn(smtpSecret)
		Expect(err).To(Succeed())
	})

	AfterEach(func() {
		Expect(cl.Delete(context.TODO(), defSMTPSecret)).To(Succeed())
	})

	It("reconciles the SMTP Secret", func() {
		expectedLabels := map[string]string{
			"presslabs.com/kind":                "smtp",
			"app.kubernetes.io/managed-by":      "project-namespace-controller.dashboard.presslabs.com",
			"dashboard.presslabs.com/reconcile": "true",
		}
		Expect(smtpSecret.GetLabels()).To(Equal(expectedLabels))
		Expect(smtpSecret.Data).To(Equal(defSMTPSecret.Data))
	})
})
