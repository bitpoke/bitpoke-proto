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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	When("creating the pod", func() {
		var (
			pod *corev1.Pod
		)

		name := "testpod"
		namespace := "default"
		key := types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}

		BeforeEach(func() {

			pod = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container-name",
							Image: "container-image",
						},
					},
				},
			}

			// create pod
			Expect(c.Create(context.TODO(), pod)).To(Succeed())
		})
		It("should been mutated", func() {
			Expect(c.Get(context.TODO(), key, pod)).To(Succeed())
			Expect(pod.Annotations["example-mutating-admission-webhook"]).To(Equal("foo"))
		})

	})
})
