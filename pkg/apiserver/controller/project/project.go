/*
	Copyright 2019 Pressinfra SRL

	This file is subject to the terms and conditions defined in file LICENSE,
	which is part of this source code package.
*/

package project

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/gosimple/slug"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	projs "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/projects/v1"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/impersonate"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
	"github.com/presslabs/dashboard/pkg/internal/organization"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

type projectsServer struct {
	client client.Client
	cfg    *rest.Config
}

// Add creates a new Project Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	projs.RegisterProjectsServiceServer(server.GRPCServer,
		NewProjectsServiceServer(server.Manager.GetClient(), rest.CopyConfig(server.Manager.GetConfig())))
	return nil
}

func (s *projectsServer) CreateProject(ctx context.Context, r *projs.CreateProjectRequest) (*projs.Project, error) {
	userID := metadata.RequireUserID(ctx)

	var err error
	var name string
	if len(r.Project.Name) > 0 {
		if name, err = project.Resolve(r.Project.Name); err != nil {
			return nil, status.InvalidArgumentf("%s", err)
		}
	} else {
		name = slug.Make(r.Project.DisplayName)
	}
	if len(name) == 0 {
		return nil, status.InvalidArgumentf("project name cannot be empty")
	}

	parent := r.Parent
	if len(r.Parent) == 0 {
		parent = metadata.RequireOrganization(ctx)
	}

	orgName, err := organization.Resolve(parent)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}

	proj := project.New(&dashboardv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: organization.NamespaceName(orgName),
			Labels: map[string]string{
				"presslabs.com/kind":         "project",
				"presslabs.com/project":      name,
				"presslabs.com/organization": orgName,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": userID,
			},
		},
	})
	proj.UpdateDisplayName(r.Project.DisplayName)

	if err = s.client.Create(context.TODO(), proj.Unwrap()); err != nil {
		return nil, status.FromError(err)
	}

	return newProjectFromK8s(proj), nil
}

func (s *projectsServer) GetProject(ctx context.Context, r *projs.GetProjectRequest) (*projs.Project, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)

	name, err := project.Resolve(r.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	orgName := metadata.RequireOrganizationName(ctx)
	key := client.ObjectKey{
		Name:      name,
		Namespace: organization.NamespaceName(orgName),
	}

	var proj dashboardv1alpha1.Project
	if err := c.Get(ctx, key, &proj); err != nil {
		return nil, status.NotFoundf("project %s not found", r.Name).Because(err)
	}

	return newProjectFromK8s(project.New(&proj)), nil
}

func (s *projectsServer) UpdateProject(ctx context.Context, r *projs.UpdateProjectRequest) (*projs.Project, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)

	name, err := project.Resolve(r.Project.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	orgName := metadata.RequireOrganizationName(ctx)
	key := client.ObjectKey{
		Name:      name,
		Namespace: organization.NamespaceName(orgName),
	}

	var proj dashboardv1alpha1.Project
	if err = c.Get(ctx, key, &proj); err != nil {
		return nil, status.NotFound().Because(err)
	}

	project.New(&proj).UpdateDisplayName(r.Project.DisplayName)

	if err = c.Update(ctx, &proj); err != nil {
		return nil, status.NotFound().Because(err)
	}

	return newProjectFromK8s(project.New(&proj)), nil
}

func (s *projectsServer) DeleteProject(ctx context.Context, r *projs.DeleteProjectRequest) (*types.Empty, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)

	name, err := project.Resolve(r.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	orgName := metadata.RequireOrganizationName(ctx)
	key := client.ObjectKey{
		Name:      name,
		Namespace: organization.NamespaceName(orgName),
	}

	var proj dashboardv1alpha1.Project
	if err := c.Get(ctx, key, &proj); err != nil {
		return nil, status.NotFound().Because()
	}

	if err := c.Delete(ctx, &proj); err != nil {
		return nil, status.FromError(err)
	}

	return &types.Empty{}, nil
}

func (s *projectsServer) ListProjects(ctx context.Context, r *projs.ListProjectsRequest) (*projs.ListProjectsResponse, error) {
	userID := metadata.RequireUserID(ctx)
	orgName := metadata.RequireOrganizationName(ctx)

	projList := &dashboardv1alpha1.ProjectList{}
	resp := &projs.ListProjectsResponse{}

	listOptions := &client.ListOptions{
		Namespace: organization.NamespaceName(orgName),
		LabelSelector: labels.SelectorFromSet(
			labels.Set{
				"presslabs.com/kind": "project",
			},
		),
	}

	if err := s.client.List(context.TODO(), listOptions, projList); err != nil {
		return nil, status.FromError(err)
	}

	// TODO: implement pagination
	resp.Projects = []projs.Project{}
	for i := range projList.Items {
		if projList.Items[i].ObjectMeta.Annotations["presslabs.com/created-by"] == userID {
			resp.Projects = append(resp.Projects, *newProjectFromK8s(project.New(&projList.Items[i])))
		}
	}

	return resp, nil
}

// NewProjectsServiceServer creates a new gRPC server for projects
func NewProjectsServiceServer(client client.Client, cfg *rest.Config) projs.ProjectsServiceServer {
	return &projectsServer{
		client: client,
		cfg:    cfg,
	}
}

func newProjectFromK8s(p *project.Project) *projs.Project {
	return &projs.Project{
		Name:         project.FQName(p.Project.ObjectMeta.Labels["presslabs.com/project"]),
		Organization: p.Project.ObjectMeta.Labels["presslabs.com/organization"],
		DisplayName:  p.Annotations["presslabs.com/display-name"],
	}
}
