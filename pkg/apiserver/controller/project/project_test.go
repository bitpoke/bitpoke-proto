/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package project

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gosimple/slug"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	projv1 "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/projects/v1"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/auth"
	"github.com/presslabs/dashboard/pkg/controller"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
	. "github.com/presslabs/dashboard/pkg/internal/testutil/gomega"
)

const (
	ctxTimeout    = time.Second * 3
	updateTimeout = time.Second
	deleteTimeout = time.Second
)

// createProject is a helper func that creates a project
func createProject(name, displayName, createdBy, organization string) *projectns.ProjectNamespace {
	proj := projectns.New(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectns.NamespaceName(name),
			Labels: map[string]string{
				"presslabs.com/kind":         "project",
				"presslabs.com/project":      name,
				"presslabs.com/organization": organization,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": createdBy,
			},
		},
	})
	proj.UpdateDisplayName(displayName)

	return proj
}

// getProjectNamespaceFn is a helper func that returns a project namespace
func getProjectNamespaceFn(ctx context.Context, c client.Client, key client.ObjectKey) func() corev1.Namespace {
	return func() corev1.Namespace {
		var p corev1.Namespace
		c.Get(ctx, key, &p)
		return p
	}
}

func expectProperNamespace(c client.Client, name, displayName, createdBy, organization string) {
	var ns corev1.Namespace
	key := client.ObjectKey{
		Name: projectns.NamespaceName(name),
	}
	Expect(c.Get(context.TODO(), key, &ns)).To(Succeed())
	Expect(ns.Name).To(Equal(fmt.Sprintf("proj-%s", name)))
	Expect(ns.Labels).To(HaveKeyWithValue("presslabs.com/kind", "project"))
	Expect(ns.Labels).To(HaveKeyWithValue("presslabs.com/project", name))
	Expect(ns.Annotations).To(HaveKeyWithValue("presslabs.com/display-name", displayName))
	Expect(ns.Annotations).To(HaveKeyWithValue("presslabs.com/created-by", createdBy))
	Expect(ns.Labels).To(HaveKeyWithValue("presslabs.com/organization", organization))
}

var _ = Describe("API server", func() {
	var (
		// stop channel for apiserver
		stop chan struct{}
		// controller k8s client
		c client.Client
		// client connection to an RPC server
		conn *grpc.ClientConn
		// projClient
		projClient projv1.ProjectsServiceClient
	)

	var (
		id, autoID     string
		name, autoName string
		displayName    string
		createdBy      string
		organization   string
		parent         string
	)

	BeforeEach(func() {
		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).To(Succeed())

		server := SetupAPIServer(mgr)
		// add ourselves to the server
		Add(server)

		// create new k8s client
		c, err = client.New(cfg, client.Options{})
		Expect(err).To(Succeed())

		// Add controllers for testing side effects
		Expect(controller.AddToManager(mgr)).To(Succeed())

		stop = StartTestManager(mgr)

		conn, err = grpc.Dial(server.GetGRPCAddr(), grpc.WithInsecure(), grpc.WithBlock(),
			grpc.WithTimeout(ctxTimeout))
		Expect(err).To(Succeed())

		projClient = projv1.NewProjectsServiceClient(conn)

		name = fmt.Sprintf("%d", rand.Int31())
		id = fmt.Sprintf("project/%s", name)
		displayName = fmt.Sprintf("Project %s", name)
		autoName = slug.Make(displayName)
		autoID = fmt.Sprintf("project/%s", autoName)
		createdBy = fmt.Sprintf("user#%s", name)
		auth.FakeSubject = createdBy
		organization = fmt.Sprintf("%d", rand.Int31())
		parent = fmt.Sprintf("orgs/%s", organization)
	})

	AfterEach(func() {
		// close the gRPC client connection
		conn.Close()
		// stop the manager and API server
		close(stop)
	})

	Describe("at Create request", func() {
		It("returns AlreadyExists error when project already exists", func() {
			proj := createProject(name, displayName, createdBy, organization)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}

			_, err := projClient.CreateProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.AlreadyExists))
		})

		It("returns error when no parent is given", func() {
			req := projv1.CreateProjectRequest{
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}
			_, err := projClient.CreateProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when parent is not fully-qualified", func() {
			req := projv1.CreateProjectRequest{
				Parent: "not-fully-qualified-parent",
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}
			_, err := projClient.CreateProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when no name is given", func() {
			req := projv1.CreateProjectRequest{
				Parent:  parent,
				Project: projv1.Project{},
			}
			_, err := projClient.CreateProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when name is not fully qualified", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name: "not-fully-qualified-name",
				},
			}
			_, err := projClient.CreateProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when name is empty", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name: "project/",
				},
			}
			_, err := projClient.CreateProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("creates project when no project name is given", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					DisplayName: displayName,
				},
			}

			resp, err := projClient.CreateProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(autoID))
			expectProperNamespace(c, slug.Make(displayName), displayName, createdBy, organization)
		})

		It("creates project when project name is given", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}

			resp, err := projClient.CreateProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			expectProperNamespace(c, name, displayName, createdBy, organization)
		})

		It("fills display_name when no one is given", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name: id,
				},
			}
			resp, err := projClient.CreateProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			expectProperNamespace(c, name, name, createdBy, organization)
		})
	})

	Describe("at Get request", func() {
		It("returns the project", func() {
			proj := createProject(name, displayName, createdBy, organization)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			req := projv1.GetProjectRequest{
				Name: id,
			}

			resp, err := projClient.GetProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.DisplayName).To(Equal(displayName))
			Expect(resp.Organization).To(Equal(organization))
		})

		It("returns NotFound when organization does not exist", func() {
			req := projv1.GetProjectRequest{
				Name: id,
			}
			_, err := projClient.GetProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

	})

	Describe("at Delete request", func() {
		It("deletes existing project", func() {
			proj := createProject(name, displayName, createdBy, organization)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			req := projv1.DeleteProjectRequest{
				Name: id,
			}

			_, err := projClient.DeleteProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			key := client.ObjectKey{
				Name: projectns.NamespaceName(name),
			}
			Eventually(getProjectNamespaceFn(context.TODO(), c, key), deleteTimeout).Should(
				BeInPhase(corev1.NamespaceTerminating))
		})

		It("returns NotFound when project does not exists", func() {
			req := projv1.DeleteProjectRequest{
				Name: id,
			}
			_, err := projClient.DeleteProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})
	})

	Describe("at update request", func() {
		BeforeEach(func() {
			proj := createProject(name, displayName, createdBy, organization)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
		})

		It("updates dispplay_name of existing project", func() {
			newDisplayName := "The New Display Name"
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name:        id,
					DisplayName: newDisplayName,
				},
			}

			resp, err := projClient.UpdateProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			Expect(resp.DisplayName).To(Equal(newDisplayName))

			key := client.ObjectKey{
				Name: projectns.NamespaceName(name),
			}
			Eventually(getProjectNamespaceFn(context.TODO(), c, key), updateTimeout).Should(
				HaveAnnotation("presslabs.com/display-name", newDisplayName))
		})

		It("sets display_name to default when no one is given", func() {
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name: id,
				},
			}

			resp, err := projClient.UpdateProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			Expect(resp.DisplayName).To(Equal(name))

			key := client.ObjectKey{
				Name: projectns.NamespaceName(name),
			}
			Eventually(getProjectNamespaceFn(context.TODO(), c, key), updateTimeout).Should(
				HaveAnnotation("presslabs.com/display-name", name))
		})

		It("returns NotFound when project does not exist", func() {
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name: "project/inexistent",
				},
			}
			_, err := projClient.UpdateProject(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("does not update the organization", func() {
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name:         id,
					Organization: "the-new-organization",
				},
			}
			resp, err := projClient.UpdateProject(context.TODO(), &req)
			Expect(err).To(Succeed())
			Expect(resp.Organization).To(Equal(organization))

			key := client.ObjectKey{
				Name: projectns.NamespaceName(name),
			}
			Eventually(getProjectNamespaceFn(context.TODO(), c, key), updateTimeout).Should(
				HaveLabel("presslabs.com/organization", organization))
		})
	})

	Describe("at list request", func() {
		var projsCount = 3
		BeforeEach(func() {
			for i := 1; i <= projsCount; i++ {
				_name := fmt.Sprintf("%s-%02d", name, i)
				_displayName := fmt.Sprintf("%s %02d Inc.", name, i)
				_organization := fmt.Sprintf("%s-%02d", organization, i)
				proj := createProject(_name, _displayName, createdBy, _organization)
				Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			}
			proj := createProject(name, displayName, "user#anoter", organization)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
		})

		It("returns only my orgnanizations", func() {
			req := projv1.ListProjectsRequest{}
			Eventually(func() ([]projv1.Project, error) {
				resp, err := projClient.ListProjects(context.TODO(), &req)
				return resp.Projects, err
			}).Should(HaveLen(projsCount))
		})
	})
})
