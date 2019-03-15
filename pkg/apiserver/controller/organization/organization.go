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

	"github.com/gogo/protobuf/types"
	"github.com/gosimple/slug"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// nolint: golint
	. "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/organizations/v1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/impersonate"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

type organizationsService struct {
	client client.Client
	cfg    *rest.Config
}

// Add creates a new Organization Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	RegisterOrganizationsServiceServer(server.GRPCServer,
		NewOrganizationsServiceServer(server.Manager.GetClient(), rest.CopyConfig(server.Manager.GetConfig())))

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
	userID := metadata.RequireUserID(ctx)

	var err error
	var name string

	if len(r.Organization.Name) > 0 {
		if name, err = organization.Resolve(r.Organization.Name); err != nil {
			return nil, status.InvalidArgumentf("%s", err)
		}
	} else {
		name = slug.Make(r.Organization.DisplayName)
	}
	if len(name) == 0 {
		return nil, status.FromError(fmt.Errorf("organization name cannot be empty"))
	}

	org := organization.New(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: organization.NamespaceName(name),
			Labels: map[string]string{
				"presslabs.com/kind":         "organization",
				"presslabs.com/organization": name,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": userID,
			},
		},
	})
	org.UpdateDisplayName(r.Organization.DisplayName)

	if err := s.client.Create(context.TODO(), org.Unwrap()); err != nil {
		return nil, status.FromError(err)
	}

	return newOrganizationFromK8s(org), nil
}

func (s *organizationsService) GetOrganization(ctx context.Context, r *GetOrganizationRequest) (*Organization, error) {
	c, _, err := impersonate.ClientFromContext(ctx, s.cfg)
	if err != nil {
		return nil, status.FromError(err)
	}

	name, err := organization.Resolve(r.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	key := client.ObjectKey{
		Name: organization.NamespaceName(name),
	}

	var org corev1.Namespace
	if err := c.Get(ctx, key, &org); err != nil {
		return nil, status.NotFoundf("organization %s not found", r.Name).Because(err)
	}

	return newOrganizationFromK8s(organization.New(&org)), nil
}

func (s *organizationsService) UpdateOrganization(ctx context.Context, r *UpdateOrganizationRequest) (*Organization, error) {
	c, _, err := impersonate.ClientFromContext(ctx, s.cfg)
	if err != nil {
		return nil, status.FromError(err)
	}

	name, err := organization.Resolve(r.Organization.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	key := client.ObjectKey{
		Name: organization.NamespaceName(name),
	}

	var org corev1.Namespace
	if err = c.Get(ctx, key, &org); err != nil {
		return nil, status.NotFound().Because(err)
	}

	organization.New(&org).UpdateDisplayName(r.Organization.DisplayName)

	if err = c.Update(ctx, &org); err != nil {
		return nil, status.FromError(err)
	}

	return newOrganizationFromK8s(organization.New(&org)), nil
}

func (s *organizationsService) DeleteOrganization(ctx context.Context, r *DeleteOrganizationRequest) (*types.Empty, error) {
	c, _, err := impersonate.ClientFromContext(ctx, s.cfg)
	if err != nil {
		return nil, status.FromError(err)
	}

	name, err := organization.Resolve(r.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	key := client.ObjectKey{
		Name: organization.NamespaceName(name),
	}

	var org corev1.Namespace
	if err := c.Get(ctx, key, &org); err != nil {
		return nil, status.NotFound().Because(err)
	}

	if err := c.Delete(ctx, &org); err != nil {
		return nil, status.FromError(err)
	}

	return &types.Empty{}, nil
}

func (s *organizationsService) ListOrganizations(ctx context.Context, r *ListOrganizationsRequest) (*ListOrganizationsResponse, error) {
	var err error
	userID := metadata.RequireUserID(ctx)

	memberLists := &rbacv1.RoleBindingList{}
	opts := client.MatchingField("subject.user", userID)
	if err = s.client.List(context.TODO(), opts, memberLists); err != nil {
		return nil, status.FromError(err)
	}

	var orgNames []string
	for _, list := range memberLists.Items {
		if list.Name == "members" && len(list.Labels) > 0 && len(list.Labels["presslabs.com/organization"]) > 0 {
			orgNames = append(orgNames, list.Labels["presslabs.com/organization"])
		}
	}

	orgs := &corev1.NamespaceList{}
	opts = &client.ListOptions{}
	err = opts.SetLabelSelector(fmt.Sprintf("presslabs.com/kind=organization, presslabs.com/organization in (%s)", strings.Join(orgNames, ", ")))
	if err != nil {
		return nil, status.FromError(err)
	}
	resp := &ListOrganizationsResponse{}

	if err := s.client.List(context.TODO(), opts, orgs); err != nil {
		return nil, status.FromError(err)
	}

	// TODO: implement pagination
	resp.Organizations = []Organization{} //make([]*Organization, len(orgs.Items))
	for i := range orgs.Items {
		resp.Organizations = append(resp.Organizations, *newOrganizationFromK8s(organization.New(&orgs.Items[i])))
	}

	return resp, nil
}

// NewOrganizationsServiceServer creates a new gRPC organizations service server
func NewOrganizationsServiceServer(client client.Client, cfg *rest.Config) OrganizationsServiceServer {
	return &organizationsService{
		client: client,
		cfg:    cfg,
	}
}

func newOrganizationFromK8s(o *organization.Organization) *Organization {
	return &Organization{
		Name:        organization.FQName(o.Namespace.ObjectMeta.Labels["presslabs.com/organization"]),
		DisplayName: o.Annotations["presslabs.com/display-name"],
	}
}
