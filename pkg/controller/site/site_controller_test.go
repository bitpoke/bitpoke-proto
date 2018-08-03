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

package site

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const timeout = time.Second * 5

var _ = Describe("Project controller", func() {
	var (
		// channel for incoming reconcile requests
		// requests chan reconcile.Request
		// stop channel for controller manager
		stop chan struct{}
		// controller k8s client
		c    client.Client
		proj dashboardv1alpha1.Project
	)

	BeforeEach(func() {
		var recFn reconcile.Reconciler

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())
		c = mgr.GetClient()

		recFn, _ = SetupTestReconcile(newReconciler(mgr)) //recfn, requests
		Expect(add(mgr, recFn)).To(Succeed())

		stop = StartTestManager(mgr)

		proj = dashboardv1alpha1.Project{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("proj-%d", rand.Int31()),
				Namespace: "default",
			},
		}
		defer c.Delete(context.TODO(), &proj)
	})

	AfterEach(func() {
		time.Sleep(1 * time.Second)
		close(stop)
	})

	Describe("when creating a new Site object", func() {
		// var expectedRequest reconcile.Request
		// var site *wpapiv1.Wordpress
		//
		// BeforeEach(func() {
		// 	name := fmt.Sprintf("wp-%d", rand.Int31())
		//
		// 	expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: proj.GetNamespaceName()}}
		// 	site = &wpapiv1.Wordpress{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name:      name,
		// 			Namespace: proj.GetNamespaceName(),
		// 		},
		// 	}
		// })
	})
})
