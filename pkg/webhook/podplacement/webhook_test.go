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

package podplacement_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	crwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/presslabs/dashboard/pkg/webhook"
	"github.com/presslabs/dashboard/pkg/webhook/podplacement"
)

var _ = Describe("Webhook server", func() {
	var (
		// stop channel for controller manager
		stop chan struct{}
		// controller k8s client
		c client.Client
		//
		server *crwebhook.Server
	)

	BeforeEach(func() {

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())
		c = mgr.GetClient()

		server, err = webhook.NewServer(mgr)
		Expect(err).NotTo(HaveOccurred())

		Expect(podplacement.AddToServer(mgr, server)).To(Succeed())

		stop = StartTestManager(mgr)

		time.Sleep(time.Second)
	})

	AfterEach(func() {
		close(stop)
	})

	When("registrd", func() {
		BeforeEach(func() {
			whc := &admissionregistrationv1beta1.ValidatingWebhookConfiguration{}
			key := types.NamespacedName{
				Name: server.Name,
			}
			c.Get(context.TODO(), key, whc)
			fmt.Printf("%v\n", whc)
		})
		It("---", func() {

		})

	})
})
