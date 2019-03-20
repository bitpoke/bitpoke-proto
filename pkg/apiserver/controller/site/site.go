/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package site

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sites "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/sites/v1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/impersonate"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
	"github.com/presslabs/dashboard/pkg/internal/project"
	"github.com/presslabs/dashboard/pkg/internal/site"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

type sitesService struct {
	client client.Client
	cfg    *rest.Config
}

// containsString returns true if a string is present in a iteratee.
func containsString(s []string, e string) bool {
	for _, ss := range s {
		if ss == e {
			return true
		}
	}
	return false
}

// Add creates a new Site Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	sites.RegisterSitesServiceServer(server.GRPCServer,
		NewSitesServiceServer(server.Manager.GetClient(), rest.CopyConfig(server.Manager.GetConfig())))
	return nil
}

func (s *sitesService) CreateSite(ctx context.Context, r *sites.CreateSiteRequest) (*sites.Site, error) {
	c, userID := impersonate.ClientFromContext(ctx, s.cfg)

	name, proj, err := site.Resolve(r.Site.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	parent, err := project.Resolve(r.Parent)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	if strings.Compare(parent, proj) != 0 {
		return nil, status.InvalidArgumentf("parent and project are not matching")
	}
	if len(r.Site.PrimaryDomain) == 0 {
		return nil, status.InvalidArgumentf("primary domain cannot be empty")
	}

	wp := site.New(&wordpressv1alpha1.Wordpress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: parent,
			Labels: map[string]string{
				"presslabs.com/kind":    "site",
				"presslabs.com/site":    name,
				"presslabs.com/project": proj,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": userID,
			},
		},
		Spec: wordpressv1alpha1.WordpressSpec{
			Domains: []wordpressv1alpha1.Domain{wordpressv1alpha1.Domain(r.Site.PrimaryDomain)},
			Image:   r.Site.WordpressImage,
		},
	})

	if err := c.Create(context.TODO(), wp.Unwrap()); err != nil {
		return nil, status.FromError(err)
	}

	return newSiteFromK8s(wp), nil
}

func (s *sitesService) GetSite(ctx context.Context, r *sites.GetSiteRequest) (*sites.Site, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)

	key, err := site.ResolveToObjectKey(r.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}

	wp := site.New(&wordpressv1alpha1.Wordpress{})
	if err := c.Get(ctx, key, wp.Unwrap()); err != nil {
		return nil, status.NotFoundf("site %s not found", r.Name).Because(err)
	}

	return newSiteFromK8s(site.New(wp.Unwrap())), nil
}

// updatePrimaryDomain updates the primary domain
func updatePrimaryDomain(wp *wordpressv1alpha1.Wordpress, domain string, fieldMask types.FieldMask) error {
	if len(fieldMask.Paths) == 0 || containsString(fieldMask.Paths, "site.primary_domain") {
		if len(domain) == 0 {
			return fmt.Errorf("primary domain cannot be empty")
		} else if len(domain) > 0 {
			wp.Spec.Domains[0] = wordpressv1alpha1.Domain(domain)
		}
	}
	return nil
}

// updateWordpressImage updates the wordpress image
func updateWordpressImage(wp *wordpressv1alpha1.Wordpress, image string, fieldMask types.FieldMask) error {
	if len(fieldMask.Paths) == 0 || containsString(fieldMask.Paths, "site.wordpress_image") {
		wp.Spec.Image = image
	}
	return nil
}

func (s *sitesService) UpdateSite(ctx context.Context, r *sites.UpdateSiteRequest) (*sites.Site, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)

	wp := site.New(&wordpressv1alpha1.Wordpress{})
	key, err := site.ResolveToObjectKey(r.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}

	// get the site
	if err = c.Get(ctx, key, wp.Unwrap()); err != nil {
		return nil, status.NotFoundf("site %s not found", r.Site.Name).Because(err)
	}

	// update primary domain and wordpress image
	if err = updatePrimaryDomain(wp.Unwrap(), r.PrimaryDomain, r.FieldMask); err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	if err = updateWordpressImage(wp.Unwrap(), r.WordpressImage, r.FieldMask); err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}
	if err = c.Update(ctx, wp.Unwrap()); err != nil {
		return nil, status.FromError(err)
	}

	return newSiteFromK8s(site.New(wp.Unwrap())), nil
}

func (s *sitesService) DeleteSite(ctx context.Context, r *sites.DeleteSiteRequest) (*types.Empty, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)

	key, err := site.ResolveToObjectKey(r.Name)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}

	wp := site.New(&wordpressv1alpha1.Wordpress{})
	if err := c.Get(ctx, key, wp.Unwrap()); err != nil {
		return nil, status.NotFound().Because(err)
	}

	if err := c.Delete(ctx, wp.Unwrap()); err != nil {
		return nil, status.FromError(err)
	}

	return &types.Empty{}, nil
}

func (s *sitesService) ListSites(ctx context.Context, r *sites.ListSitesRequest) (*sites.ListSitesResponse, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)

	wpList := &wordpressv1alpha1.WordpressList{}
	resp := &sites.ListSitesResponse{}

	if len(r.Parent) <= 0 {
		return nil, status.InvalidArgumentf("parent cannot be empty")
	}

	ns, err := project.Resolve(r.Parent)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}

	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set{
				"presslabs.com/kind": "site",
			},
		),
		Namespace: ns,
	}

	if err := c.List(context.TODO(), listOptions, wpList); err != nil {
		return nil, status.FromError(err)
	}

	// TODO: implement pagination
	resp.Sites = make([]sites.Site, len(wpList.Items))
	for i := range wpList.Items {
		resp.Sites[i] = *newSiteFromK8s(site.New(&wpList.Items[i]))
	}

	return resp, nil
}

// NewSitesServiceServer creates a new gRPC sites service server
func NewSitesServiceServer(client client.Client, cfg *rest.Config) sites.SitesServiceServer {
	return &sitesService{
		client: client,
		cfg:    cfg,
	}
}

func newSiteFromK8s(s *site.Site) *sites.Site {
	return &sites.Site{
		Name:           site.FQName(s.ObjectMeta.Labels["presslabs.com/project"], s.ObjectMeta.Labels["presslabs.com/site"]),
		PrimaryDomain:  string(s.Spec.Domains[0]),
		WordpressImage: s.Spec.Image,
	}
}
