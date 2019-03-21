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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gosimple/slug"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	. "github.com/presslabs/dashboard/pkg/internal/testutil/gomega"

	projv1 "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/projects/v1"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/controller"
	"github.com/presslabs/dashboard/pkg/internal/organization"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

const (
	ctxTimeout    = time.Second * 3
	updateTimeout = time.Second
)

// createProject is a helper func that creates a project
func createProject(name, displayName, createdBy, org string) *project.Project {
	proj := project.New(&dashboardv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: organization.NamespaceName(org),
			Labels: map[string]string{
				"presslabs.com/kind":         "project",
				"presslabs.com/project":      name,
				"presslabs.com/organization": org,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": createdBy,
			},
		},
	})
	proj.UpdateDisplayName(displayName)

	return proj
}

// getProjectFn is a helper func that returns a project
func getProjectFn(ctx context.Context, c client.Client, key client.ObjectKey) func() dashboardv1alpha1.Project {
	return func() dashboardv1alpha1.Project {
		var p dashboardv1alpha1.Project
		c.Get(ctx, key, &p)
		return p
	}
}

func expectProperProject(c client.Client, name, displayName, createdBy, org string) {
	var p dashboardv1alpha1.Project
	key := client.ObjectKey{
		Name:      name,
		Namespace: organization.NamespaceName(org),
	}
	Expect(c.Get(context.TODO(), key, &p)).To(Succeed())
	Expect(p.Name).To(Equal(fmt.Sprintf("%s", name)))
	Expect(p.Labels).To(HaveKeyWithValue("presslabs.com/kind", "project"))
	Expect(p.Labels).To(HaveKeyWithValue("presslabs.com/project", name))
	Expect(p.Annotations).To(HaveKeyWithValue("presslabs.com/display-name", displayName))
	Expect(p.Annotations).To(HaveKeyWithValue("presslabs.com/created-by", createdBy))
	Expect(p.Labels).To(HaveKeyWithValue("presslabs.com/organization", org))
}

var _ = Describe("API project controller", func() {
	var (
		// stop channel for apiserver
		stop chan struct{}
		// controller k8s client
		c client.Client
		// client connection to an RPC server
		conn *grpc.ClientConn
		// projClient
		projClient projv1.ProjectsServiceClient
		// context for requests
		orgCtx context.Context
	)

	var (
		id, autoID     string
		name, autoName string
		displayName    string
		createdBy      string
		org            string
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
		metadata.FakeSubject = createdBy
		org = fmt.Sprintf("%d", rand.Int31())
		parent = fmt.Sprintf("orgs/%s", org)

		orgCtx = metadata.AddOrgInContext(context.Background(), parent)
	})

	AfterEach(func() {
		// close the gRPC client connection
		conn.Close()
		// stop the manager and API server
		close(stop)
	})

	Describe("at Create request", func() {
		It("returns AlreadyExists error when project already exists", func() {
			proj := createProject(name, displayName, createdBy, org)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}

			_, err := projClient.CreateProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.AlreadyExists))
		})

		It("returns error when no parent is given", func() {
			req := projv1.CreateProjectRequest{
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}
			_, err := projClient.CreateProject(context.Background(), &req)
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
			_, err := projClient.CreateProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when no name is given", func() {
			req := projv1.CreateProjectRequest{
				Parent:  parent,
				Project: projv1.Project{},
			}
			_, err := projClient.CreateProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when name is not fully qualified", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name: "not-fully-qualified-name",
				},
			}
			_, err := projClient.CreateProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when name is empty", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name: "project/",
				},
			}
			_, err := projClient.CreateProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("creates project when no project name is given", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					DisplayName: displayName,
				},
			}

			resp, err := projClient.CreateProject(orgCtx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(autoID))
			expectProperProject(c, slug.Make(displayName), displayName, createdBy, org)
		})

		It("creates project when project name is given", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}

			resp, err := projClient.CreateProject(orgCtx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			expectProperProject(c, name, displayName, createdBy, org)
		})

		It("fills display_name when no one is given", func() {
			req := projv1.CreateProjectRequest{
				Parent: parent,
				Project: projv1.Project{
					Name: id,
				},
			}
			resp, err := projClient.CreateProject(orgCtx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			expectProperProject(c, name, name, createdBy, org)
		})

		It("returns error when no parent is given and no organization is set in metadata", func() {
			req := projv1.CreateProjectRequest{
				Project: projv1.Project{
					Name:        id,
					DisplayName: displayName,
				},
			}

			_, err := projClient.CreateProject(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("at Get request", func() {
		It("returns the project", func() {
			proj := createProject(name, displayName, createdBy, org)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			req := projv1.GetProjectRequest{
				Name: id,
			}

			resp, err := projClient.GetProject(orgCtx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.DisplayName).To(Equal(displayName))
			Expect(resp.Organization).To(Equal(org))
		})

		It("returns NotFound when organization does not exist", func() {
			req := projv1.GetProjectRequest{
				Name: id,
			}
			_, err := projClient.GetProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("returns error when header does not have organization id", func() {
			req := projv1.GetProjectRequest{
				Name: id,
			}
			_, err := projClient.GetProject(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("at Delete request", func() {
		It("deletes existing project", func() {
			proj := createProject(name, displayName, createdBy, org)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			req := projv1.DeleteProjectRequest{
				Name: id,
			}

			_, err := projClient.DeleteProject(orgCtx, &req)
			Expect(err).To(Succeed())
			key := client.ObjectKey{
				Name:      name,
				Namespace: organization.NamespaceName(org),
			}
			var p dashboardv1alpha1.Project
			err = c.Get(orgCtx, key, &p)
			Expect(status.Code(err)).To(Equal(codes.Unknown))
		})

		It("returns NotFound when project does not exists", func() {
			req := projv1.DeleteProjectRequest{
				Name: id,
			}
			_, err := projClient.DeleteProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("returns error when header does not have organization id", func() {
			req := projv1.DeleteProjectRequest{
				Name: id,
			}
			_, err := projClient.DeleteProject(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("at update request", func() {
		BeforeEach(func() {
			proj := createProject(name, displayName, createdBy, org)
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

			resp, err := projClient.UpdateProject(orgCtx, &req)
			Expect(err).To(Succeed())
			Expect(resp.DisplayName).To(Equal(newDisplayName))

			key := client.ObjectKey{
				Name:      name,
				Namespace: organization.NamespaceName(org),
			}
			Eventually(getProjectFn(context.TODO(), c, key), updateTimeout).Should(
				HaveAnnotation("presslabs.com/display-name", newDisplayName))
		})

		It("sets display_name to default when no one is given", func() {
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name: id,
				},
			}

			resp, err := projClient.UpdateProject(orgCtx, &req)
			Expect(err).To(Succeed())
			Expect(resp.DisplayName).To(Equal(name))

			key := client.ObjectKey{
				Name:      name,
				Namespace: organization.NamespaceName(org),
			}
			Eventually(getProjectFn(context.TODO(), c, key), updateTimeout).Should(
				HaveAnnotation("presslabs.com/display-name", name))
		})

		It("returns NotFound when project does not exist", func() {
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name: "project/inexistent",
				},
			}
			_, err := projClient.UpdateProject(orgCtx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("does not update the organization", func() {
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name:         id,
					Organization: "the-new-organization",
				},
			}
			resp, err := projClient.UpdateProject(orgCtx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Organization).To(Equal(org))

			key := client.ObjectKey{
				Name:      name,
				Namespace: organization.NamespaceName(org),
			}
			Eventually(getProjectFn(context.TODO(), c, key), updateTimeout).Should(
				HaveLabel("presslabs.com/organization", org))
		})

		It("returns error when header does not have organization id", func() {
			req := projv1.UpdateProjectRequest{
				Project: projv1.Project{
					Name: id,
				},
			}
			_, err := projClient.UpdateProject(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("at list request", func() {
		var projsCount = 3
		BeforeEach(func() {
			for i := 1; i <= projsCount; i++ {
				_name := fmt.Sprintf("%s-%02d", name, i)
				_displayName := fmt.Sprintf("%s %02d Inc.", name, i)
				proj := createProject(_name, _displayName, createdBy, org)
				Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
			}
			proj := createProject(name, displayName, "user#anoter", org)
			Expect(c.Create(context.TODO(), proj.Unwrap())).To(Succeed())
		})

		It("returns only my orgnanizations", func() {
			req := projv1.ListProjectsRequest{}
			Eventually(func() ([]projv1.Project, error) {
				resp, err := projClient.ListProjects(orgCtx, &req)
				return resp.Projects, err
			}).Should(HaveLen(projsCount))
		})

		It("returns error when header doesn not have organization id", func() {
			req := projv1.ListProjectsRequest{}
			_, err := projClient.ListProjects(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})
})
