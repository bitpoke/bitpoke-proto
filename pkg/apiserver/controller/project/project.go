/*
	Copyright 2019 Pressinfra SRL

	This file is subject to the terms and conditions defined in file LICENSE,
	which is part of this source code package.
*/

package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/gosimple/slug"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// nolint: golint
	. "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/projects/v1"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/auth"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/header"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/impersonate"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

type projectsServer struct {
	client client.Client
	cfg    *rest.Config
}

const (
	projPrefix = "project/"
	orgPrefix  = "orgs/"
)

// resolveName resolves a fully-qualified project name to a k8s object name
func resolveName(path string) (string, error) {
	if !strings.HasPrefix(path, projPrefix) {
		return "", status.InvalidArgumentf("project fully-qualified name must be in form project/PROJECT-NAME, '%s' given", path)
	}
	name := path[len(projPrefix):]
	if len(name) == 0 {
		return "", status.InvalidArgumentf("project fully-qualified name must be in form project/PROJECT-NAME, '%s' given", path)
	}
	return name, nil
}

// resolveParent resolves a fully qualified parent name to a k8s object name
func resolveParent(path string) (string, error) {
	if !strings.HasPrefix(path, orgPrefix) {
		return "", status.InvalidArgumentf("parent organization fully-qualified name must be in form orgs/ORGANIZATION-NAME")
	}
	name := path[len(orgPrefix):]
	if len(name) == 0 {
		return "", status.InvalidArgumentf("parent organization fully-qualified name must be in form orgs/ORGANIZATION_NAME, '%s' given", path)
	}
	return name, nil
}

// Add creates a new Project Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	RegisterProjectsServiceServer(server.GRPCServer,
		NewProjectsServiceServer(server.Manager.GetClient(), rest.CopyConfig(server.Manager.GetConfig())))
	return nil
}

func (s *projectsServer) CreateProject(ctx context.Context, r *CreateProjectRequest) (*Project, error) {
	userID := auth.UserID(ctx)

	var err error
	var name string
	if len(r.Project.Name) > 0 {
		if name, err = resolveName(r.Project.Name); err != nil {
			return nil, err
		}
	} else {
		name = slug.Make(r.Project.DisplayName)
	}
	if len(name) == 0 {
		return nil, status.InvalidArgumentf("project name cannot be empty")
	}
	parent, err := resolveParent(r.Parent)
	if err != nil {
		return nil, err
	}
	ns := header.OrgFromContext(ctx)
	if len(ns) == 0 {
		return nil, status.InvalidArgumentf("organization cannot be empty")
	}

	proj := project.New(&dashboardv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"presslabs.com/kind":         "project",
				"presslabs.com/project":      name,
				"presslabs.com/organization": parent,
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

func (s *projectsServer) GetProject(ctx context.Context, r *GetProjectRequest) (*Project, error) {
	c, _, err := impersonate.ClientFromContext(ctx, s.cfg)
	if err != nil {
		return nil, status.FromError(err)
	}

	name, err := resolveName(r.Name)
	if err != nil {
		return nil, status.FromError(err)
	}
	ns := header.OrgFromContext(ctx)
	if ns == "" {
		return nil, status.InvalidArgumentf("organization id cannot be empty")
	}
	key := client.ObjectKey{
		Name:      name,
		Namespace: ns,
	}

	var proj dashboardv1alpha1.Project
	if err := c.Get(ctx, key, &proj); err != nil {
		return nil, status.NotFoundf("project %s not found", r.Name).Because(err)
	}

	return newProjectFromK8s(project.New(&proj)), nil
}

func (s *projectsServer) UpdateProject(ctx context.Context, r *UpdateProjectRequest) (*Project, error) {
	c, _, err := impersonate.ClientFromContext(ctx, s.cfg)
	if err != nil {
		return nil, status.FromError(err)
	}

	name, err := resolveName(r.Project.Name)
	if err != nil {
		return nil, status.FromError(err)
	}
	ns := header.OrgFromContext(ctx)
	if ns == "" {
		return nil, status.InvalidArgumentf("organization id cannot be empty")
	}
	key := client.ObjectKey{
		Name:      name,
		Namespace: ns,
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

func (s *projectsServer) DeleteProject(ctx context.Context, r *DeleteProjectRequest) (*types.Empty, error) {
	c, _, err := impersonate.ClientFromContext(ctx, s.cfg)
	if err != nil {
		return nil, status.FromError(err)
	}

	name, err := resolveName(r.Name)
	if err != nil {
		return nil, status.FromError(err)
	}
	ns := header.OrgFromContext(ctx)
	if ns == "" {
		return nil, status.InvalidArgumentf("organization id cannot be empty")
	}
	key := client.ObjectKey{
		Name:      name,
		Namespace: ns,
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

func (s *projectsServer) ListProjects(ctx context.Context, r *ListProjectsRequest) (*ListProjectsResponse, error) {
	userID := auth.UserID(ctx)

	projs := &dashboardv1alpha1.ProjectList{}
	resp := &ListProjectsResponse{}

	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set{
				"presslabs.com/kind": "project",
			},
		),
	}

	if err := s.client.List(context.TODO(), listOptions, projs); err != nil {
		return nil, status.FromError(err)
	}

	// TODO: implement pagination
	resp.Projects = []Project{}
	for i := range projs.Items {
		if projs.Items[i].ObjectMeta.Annotations["presslabs.com/created-by"] == userID {
			resp.Projects = append(resp.Projects, *newProjectFromK8s(project.New(&projs.Items[i])))
		}
	}

	return resp, nil
}

// NewProjectsServiceServer creates a new gRPC server for projects
func NewProjectsServiceServer(client client.Client, cfg *rest.Config) ProjectsServiceServer {
	return &projectsServer{
		client: client,
		cfg:    cfg,
	}
}

func newProjectFromK8s(p *project.Project) *Project {
	return &Project{
		Name:         fmt.Sprintf("%s%s", projPrefix, p.Project.ObjectMeta.Labels["presslabs.com/project"]),
		Organization: p.Project.ObjectMeta.Labels["presslabs.com/organization"],
		DisplayName:  p.Annotations["presslabs.com/display-name"],
	}
}
