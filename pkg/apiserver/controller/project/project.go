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

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gosimple/slug"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// nolint: golint
	. "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/projects/v1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

type projectsServer struct {
	client client.Client
}

// resolves an fully-qualified resource name to a k8s object name
func resolve(path string) (string, error) {
	prefix := fmt.Sprintf("proj/")
	if !strings.HasPrefix(path, prefix) {
		return "", fmt.Errorf("projects resources fully-qualified name must be in form proj/PROJECT-NAME")
	}
	name := path[len(prefix):]
	if len(name) == 0 {
		return "", fmt.Errorf("project name cannot be empty")
	}
	return name, nil
}

// Add creates a new Project Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	RegisterProjectsServiceServer(server.GRPCServer, NewProjectsServiceServer(server.Manager.GetClient()))
	return nil
}

func (s *projectsServer) CreateProject(ctx context.Context, r *CreateProjectRequest) (*Project, error) {
	cl := ctx.Value(middleware.AuthTokenContextKey)
	if cl == nil {
		return nil, status.Error(fmt.Errorf("no auth-token value in context"))
	}
	createdBy := cl.(middleware.Claims).Subject

	var name string
	var err error
	if len(r.Project.Name) > 0 {
		if name, err = resolve(r.Project.Name); err != nil {
			return nil, status.Error(err)
		}
	} else {
		name = slug.Make(r.Project.DisplayName)
	}

	if len(name) == 0 {
		return nil, status.Error(fmt.Errorf("project name cannot be empty"))
	}

	if len(r.Parent) <= 0 {
		return nil, status.Error(fmt.Errorf("parent cannot be empty"))
	}
	organization := r.Parent

	proj := project.New(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.NamespaceName(name),
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
	proj.UpdateDisplayName(r.Project.DisplayName)

	if err := s.client.Create(context.TODO(), proj.Unwrap()); err != nil {
		return nil, status.Error(err)
	}

	return newProjectFromK8s(proj), nil
}

func (s *projectsServer) GetProject(ctx context.Context, r *GetProjectRequest) (*Project, error) {
	var proj corev1.Namespace
	name, err := resolve(r.Name)
	if err != nil {
		return nil, status.Error(err)
	}

	key := client.ObjectKey{
		Name: project.NamespaceName(name),
	}

	if err := s.client.Get(ctx, key, &proj); err != nil {
		return nil, status.Error(err)
	}

	return newProjectFromK8s(project.New(&proj)), nil
}

func (s *projectsServer) UpdateProject(ctx context.Context, r *UpdateProjectRequest) (*Project, error) {
	var proj corev1.Namespace
	name, err := resolve(r.Project.Name)
	if err != nil {
		return nil, status.Error(err)
	}

	key := client.ObjectKey{
		Name: project.NamespaceName(name),
	}

	if err := s.client.Get(ctx, key, &proj); err != nil {
		return nil, status.Error(err)
	}

	project.New(&proj).UpdateDisplayName(r.Project.DisplayName)

	if err := s.client.Update(ctx, &proj); err != nil {
		return nil, status.Error(err)
	}

	return newProjectFromK8s(project.New(&proj)), nil
}

func (s *projectsServer) DeleteProject(ctx context.Context, r *DeleteProjectRequest) (*empty.Empty, error) {
	var proj corev1.Namespace
	name, err := resolve(r.Name)
	if err != nil {
		return nil, status.Error(err)
	}

	key := client.ObjectKey{
		Name: project.NamespaceName(name),
	}

	if err := s.client.Get(ctx, key, &proj); err != nil {
		return nil, status.Error(err)
	}

	if err := s.client.Delete(ctx, &proj); err != nil {
		return nil, status.Error(err)
	}

	return &empty.Empty{}, nil
}

func (s *projectsServer) ListProjects(ctx context.Context, r *ListProjectsRequest) (*ListProjectsResponse, error) {
	cl := ctx.Value(middleware.AuthTokenContextKey)
	if cl == nil {
		err := fmt.Errorf("no auth-token value in context")
		return nil, status.Error(err)
	}
	createdBy := cl.(middleware.Claims).Subject

	projs := &corev1.NamespaceList{}
	resp := &ListProjectsResponse{}

	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set{
				"presslabs.com/kind": "project",
			},
		),
	}

	if err := s.client.List(context.TODO(), listOptions, projs); err != nil {
		return nil, status.Error(err)
	}

	// TODO: implement pagination
	resp.Projects = []*Project{}
	for i := range projs.Items {
		if projs.Items[i].ObjectMeta.Annotations["presslabs.com/created-by"] == createdBy {
			resp.Projects = append(resp.Projects, newProjectFromK8s(project.New(&projs.Items[i])))
		}
	}

	return resp, nil
}

// NewProjectsServiceServer creates a new gRPC server for projects
func NewProjectsServiceServer(client client.Client) ProjectsServiceServer {
	return &projectsServer{
		client: client,
	}
}

func newProjectFromK8s(p *project.Project) *Project {
	return &Project{
		Name:         fmt.Sprintf("proj/%s", p.Namespace.ObjectMeta.Labels["presslabs.com/project"]),
		Organization: p.Namespace.ObjectMeta.Labels["presslabs.com/organization"],
		DisplayName:  p.Annotations["presslabs.com/display-name"],
	}
}
