/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package organization

import (
	"context"
	"fmt"
	"strings"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/gosimple/slug"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	// nolint: golint
	. "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/organizations/v1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

type organizationsService struct {
	client client.Client
}

// resolves an fully-qualified resource name to a k8s object name
func resolve(path string) (string, error) {
	if !strings.HasPrefix(path, "orgs/") {
		return "", fmt.Errorf("organization resources fully-qualified name must be in form orgs/ORGANIZATION-NAME")
	}
	name := path[5:]
	if len(name) == 0 {
		return "", fmt.Errorf("organization name cannot be empty")
	}
	return name, nil
}

// Add creates a new Organization Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	RegisterOrganizationsServiceServer(server.GRPCServer, NewOrganizationsServiceServer(server.Manager.GetClient()))

	err := server.Manager.GetFieldIndexer().IndexField(&rbacv1.ClusterRoleBinding{}, "subject.user", func(in runtime.Object) []string {
		crb := in.(*rbacv1.ClusterRoleBinding)
		var users []string
		for _, sub := range crb.Subjects {
			if sub.APIGroup == "rbac.authorization.k8s.io" && sub.Kind == "User" {
				users = append(users, sub.Name)
			}
		}
		return users
	})
	if err != nil {
		return err
	}

	err = server.Manager.GetFieldIndexer().IndexField(&rbacv1.RoleBinding{}, "subject.user", func(in runtime.Object) []string {
		rb := in.(*rbacv1.RoleBinding)
		var users []string
		for _, sub := range rb.Subjects {
			if sub.APIGroup == "rbac.authorization.k8s.io" && sub.Kind == "User" {
				users = append(users, sub.Name)
			}
		}
		return users
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *organizationsService) CreateOrganization(ctx context.Context, r *CreateOrganizationRequest) (*Organization, error) {
	cl := ctx.Value(middleware.AuthTokenContextKey)
	var name string
	var err error
	if len(r.Organization.Name) > 0 {
		if name, err = resolve(r.Organization.Name); err != nil {
			return nil, status.Error(err)
		}
	} else {
		name = slug.Make(r.Organization.DisplayName)
	}

	if len(name) == 0 {
		return nil, status.Error(fmt.Errorf("organization name cannot be empty"))
	}

	if cl == nil {
		return nil, status.Error(fmt.Errorf("no auth-token value in context"))
	}
	createdBy := cl.(middleware.Claims).Subject

	org := organization.New(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: organization.NamespaceName(name),
			Labels: map[string]string{
				"presslabs.com/kind":         "organization",
				"presslabs.com/organization": name,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": createdBy,
			},
		},
	})
	org.UpdateDisplayName(r.Organization.DisplayName)

	if err := s.client.Create(context.TODO(), org.Unwrap()); err != nil {
		return nil, status.Error(err)
	}

	return newOrganizationFromK8s(org), nil
}

func (s *organizationsService) GetOrganization(ctx context.Context, r *GetOrganizationRequest) (*Organization, error) {
	var org corev1.Namespace
	name, err := resolve(r.Name)
	if err != nil {
		return nil, status.Error(err)
	}

	key := client.ObjectKey{
		Name: organization.NamespaceName(name),
	}

	if err := s.client.Get(ctx, key, &org); err != nil {
		return nil, status.Error(err)
	}

	return newOrganizationFromK8s(organization.New(&org)), nil
}

func (s *organizationsService) UpdateOrganization(ctx context.Context, r *UpdateOrganizationRequest) (*Organization, error) {
	var org corev1.Namespace
	name, err := resolve(r.Organization.Name)
	if err != nil {
		return nil, status.Error(err)
	}

	key := client.ObjectKey{
		Name: organization.NamespaceName(name),
	}

	if err := s.client.Get(ctx, key, &org); err != nil {
		return nil, status.Error(err)
	}

	organization.New(&org).UpdateDisplayName(r.Organization.DisplayName)

	if err := s.client.Update(ctx, &org); err != nil {
		return nil, status.Error(err)
	}

	return newOrganizationFromK8s(organization.New(&org)), nil
}

func (s *organizationsService) DeleteOrganization(ctx context.Context, r *DeleteOrganizationRequest) (*empty.Empty, error) {
	var org corev1.Namespace
	name, err := resolve(r.Name)
	if err != nil {
		return nil, status.Error(err)
	}

	key := client.ObjectKey{
		Name: organization.NamespaceName(name),
	}

	if err := s.client.Get(ctx, key, &org); err != nil {
		return nil, status.Error(err)
	}

	if err := s.client.Delete(ctx, &org); err != nil {
		return nil, status.Error(err)
	}

	return &empty.Empty{}, nil
}

func (s *organizationsService) ListOrganizations(ctx context.Context, r *ListOrganizationsRequest) (*ListOrganizationsResponse, error) {
	cl := ctx.Value(middleware.AuthTokenContextKey)
	if cl == nil {
		return nil, status.Error(fmt.Errorf("no auth-token value in context"))
	}
	userID := cl.(middleware.Claims).Subject

	memberLists := &rbacv1.RoleBindingList{}
	opts := client.MatchingField("subject.user", userID)
	if err := s.client.List(context.TODO(), opts, memberLists); err != nil {
		return nil, status.Error(err)
	}

	var orgNames []string
	for _, list := range memberLists.Items {
		if list.Name == "members" && len(list.Labels) > 0 && len(list.Labels["presslabs.com/organization"]) > 0 {
			orgNames = append(orgNames, list.Labels["presslabs.com/organization"])
		}
	}

	orgs := &corev1.NamespaceList{}
	opts = &client.ListOptions{}
	err := opts.SetLabelSelector(fmt.Sprintf("presslabs.com/kind=organization, presslabs.com/organization in (%s)", strings.Join(orgNames, ", ")))
	if err != nil {
		return nil, status.Error(err)
	}
	resp := &ListOrganizationsResponse{}

	if err := s.client.List(context.TODO(), opts, orgs); err != nil {
		return nil, status.Error(err)
	}

	// TODO: implement pagination
	resp.Organizations = []*Organization{} //make([]*Organization, len(orgs.Items))
	for i := range orgs.Items {
		resp.Organizations = append(resp.Organizations, newOrganizationFromK8s(organization.New(&orgs.Items[i])))
	}

	return resp, nil
}

// NewOrganizationsServiceServer creates a new gRPC organizations service server
func NewOrganizationsServiceServer(client client.Client) OrganizationsServiceServer {
	return &organizationsService{
		client: client,
	}
}

func newOrganizationFromK8s(o *organization.Organization) *Organization {
	return &Organization{
		Name:        fmt.Sprintf("orgs/%s", o.Namespace.ObjectMeta.Labels["presslabs.com/organization"]),
		DisplayName: o.Annotations["presslabs.com/display-name"],
	}
}
