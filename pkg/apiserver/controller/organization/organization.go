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
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/presslabs/dashboard/pkg/apiserver/errors"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	// nolint: golint
	. "github.com/presslabs/dashboard/pkg/api/organizations/v1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

type organizationsServer struct {
	client client.Client
}

// Add creates a new Organization Controller and adds it to the API Server
func Add(mgr manager.Manager, auth grpc_auth.AuthFunc, grpcAddr, httpAddr string) error {
	apiServer, err := apiserver.AddToServer(mgr, auth, grpcAddr, httpAddr)
	if err != nil {
		return err
	}
	RegisterOrganizationsServiceServer(apiServer.GRPCServer, NewOrganizationsServer(mgr.GetClient()))
	return nil
}

func (s *organizationsServer) CreateOrganization(ctx context.Context, r *CreateOrganizationRequest) (*Organization, error) {
	cl := ctx.Value(middleware.AuthTokenContextKey)
	if cl == nil {
		return nil, fmt.Errorf("No auth-token value in context")
	}
	createdBy := cl.(middleware.Claims).Subject

	org := organization.New(r.Organization.Name, r.Organization.DisplayName, createdBy)

	if err := s.client.Create(context.TODO(), org.Unwrap()); err != nil {
		return nil, errors.NewApiserverError(err)
	}

	return newOrganizationFromK8s(org), nil
}

func (s *organizationsServer) GetOrganization(ctx context.Context, r *GetOrganizationRequest) (*Organization, error) {
	var org corev1.Namespace
	key := client.ObjectKey{
		Name: organization.NamespaceName(r.Name),
	}

	if err := s.client.Get(ctx, key, &org); err != nil {
		return nil, errors.NewApiserverError(err)
	}

	return newOrganizationFromK8s(organization.Wrap(&org)), nil
}

func (s *organizationsServer) UpdateOrganization(ctx context.Context, r *UpdateOrganizationRequest) (*Organization, error) {
	var org corev1.Namespace
	key := client.ObjectKey{
		Name: organization.NamespaceName(r.Organization.Name),
	}

	if err := s.client.Get(ctx, key, &org); err != nil {
		return nil, errors.NewApiserverError(err)
	}

	organization.Wrap(&org).UpdateDisplayName(r.Organization.DisplayName)

	if err := s.client.Update(ctx, &org); err != nil {
		return nil, errors.NewApiserverError(err)
	}

	return newOrganizationFromK8s(organization.Wrap(&org)), nil
}

func (s *organizationsServer) DeleteOrganization(ctx context.Context, r *DeleteOrganizationRequest) (*empty.Empty, error) {
	var org corev1.Namespace
	key := client.ObjectKey{
		Name: organization.NamespaceName(r.Name),
	}

	if err := s.client.Get(ctx, key, &org); err != nil {
		return nil, errors.NewApiserverError(err)
	}

	if err := s.client.Delete(ctx, &org); err != nil {
		return nil, errors.NewApiserverError(err)
	}

	return &empty.Empty{}, nil
}

func (s *organizationsServer) ListOrganizations(ctx context.Context, r *ListOrganizationsRequest) (*ListOrganizationsResponse, error) {
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
		return nil, errors.NewApiserverError(err)
	}

	// TODO: implement pagination
	resp.Organizations = []*Organization{} //make([]*Organization, len(orgs.Items))
	for i := range orgs.Items {
		if orgs.Items[i].ObjectMeta.Annotations["presslabs.com/created-by"] == createdBy {
			resp.Organizations = append(resp.Organizations, newOrganizationFromK8s(organization.Wrap(&orgs.Items[i])))
		}
	}

	return resp, nil
}

// NewOrganizationsServer creates a new gRPC server for organizations
func NewOrganizationsServer(client client.Client) OrganizationsServiceServer {
	return &organizationsServer{
		client: client,
	}
}

func newOrganizationFromK8s(o *organization.Organization) *Organization {
	return &Organization{
		Name:        o.Namespace.ObjectMeta.Labels["presslabs.com/organization"],
		DisplayName: o.Annotations["presslabs.com/display-name"],
	}
}
