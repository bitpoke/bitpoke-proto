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

package apiserver_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	pb "github.com/presslabs/dashboard/pkg/apiserver/projects/v1"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _ = Describe("API server", func() {
	var (
		// stop channel for controller manager
		stop chan struct{}
		// controller k8s client
		c client.Client

		conn       *grpc.ClientConn
		grpcClient pb.ProjectsClient
	)

	BeforeEach(func() {
		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())
		c = mgr.GetClient()

		Expect(apiserver.AddToManager(mgr)).To(Succeed())

		stop = StartTestManager(mgr)

		conn, err = grpc.Dial("localhost:6060", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))

		Expect(err).NotTo(HaveOccurred())
		grpcClient = pb.NewProjectsClient(conn)
	})

	AfterEach(func() {
		close(stop)
		conn.Close()
	})

	Describe("Projects Endpoints", func() {
		It("allows listing", func() {
			instance := &dashboardv1alpha1.Project{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "proj-x",
					Namespace: "default",
				},
			}
			Expect(c.Create(context.TODO(), instance)).To(Succeed())

			stream, err := grpcClient.ListProjects(context.TODO(), &pb.ListRequest{})
			Expect(err).NotTo(HaveOccurred())

			project, err := stream.Recv()
			Expect(err).NotTo(HaveOccurred())

			Expect(project.Name).To(Equal("proj-x"))
		})
	})
})
