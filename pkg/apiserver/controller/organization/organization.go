/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package organization

import (
	"context"
	"fmt"

	empty "github.com/golang/protobuf/ptypes/empty"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	// nolint: golint
	. "github.com/presslabs/dashboard/pkg/api/organizations/v1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

type organizationsService struct {
	client client.Client
}

// Add creates a new Organization Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	RegisterOrganizationsServiceServer(server.GRPCServer, NewOrganizationsServiceServer(server.Client))
	return nil
}

func (s *organizationsService) CreateOrganization(ctx context.Context, r *CreateOrganizationRequest) (*Organization, error) {
	cl := ctx.Value(middleware.AuthTokenContextKey)
	if cl == nil {
		return nil, fmt.Errorf("No auth-token value in context")
	}
	createdBy := cl.(middleware.Claims).Subject

	org := organization.New(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: organization.NamespaceName(r.Organization.Name),
			Labels: map[string]string{
				"presslabs.com/kind":         "organization",
				"presslabs.com/organization": r.Organization.Name,
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
	key := client.ObjectKey{
		Name: organization.NamespaceName(r.Name),
	}

	if err := s.client.Get(ctx, key, &org); err != nil {
		return nil, status.Error(err)
	}

	return newOrganizationFromK8s(organization.New(&org)), nil
}

func (s *organizationsService) UpdateOrganization(ctx context.Context, r *UpdateOrganizationRequest) (*Organization, error) {
	var org corev1.Namespace
	key := client.ObjectKey{
		Name: organization.NamespaceName(r.Organization.Name),
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
	key := client.ObjectKey{
		Name: organization.NamespaceName(r.Name),
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
		return nil, fmt.Errorf("No auth-token value in context")
	}
	createdBy := cl.(middleware.Claims).Subject

	orgs := &corev1.NamespaceList{}
	resp := &ListOrganizationsResponse{}

	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set{
				"presslabs.com/kind": "organization",
			},
		),
	}

	if err := s.client.List(context.TODO(), listOptions, orgs); err != nil {
		return nil, status.Error(err)
	}

	// TODO: implement pagination
	resp.Organizations = []*Organization{} //make([]*Organization, len(orgs.Items))
	for i := range orgs.Items {
		if orgs.Items[i].ObjectMeta.Annotations["presslabs.com/created-by"] == createdBy {
			resp.Organizations = append(resp.Organizations, newOrganizationFromK8s(organization.New(&orgs.Items[i])))
		}
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
		Name:        o.Namespace.ObjectMeta.Labels["presslabs.com/organization"],
		DisplayName: o.Annotations["presslabs.com/display-name"],
	}
}
