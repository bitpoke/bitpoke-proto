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
	"golang.org/x/net/publicsuffix"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/rand"

	sites "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/sites/v1"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/impersonate"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/site"
	"github.com/presslabs/dashboard/pkg/apiserver/status"
	"github.com/presslabs/dashboard/pkg/internal/project"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
	dashboardsite "github.com/presslabs/dashboard/pkg/internal/site"
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

// generateNameFromDomain generates site name from primary domain if no one is given
func generateNameFromDomain(domain string) string {
	etld, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		panic(err)
	}
	label := strings.Split(etld, ".")[0]

	gen := rand.NewStringGenerator("abcdefghijklmnopqrstuvwxyz0123456789")
	randomString, err := gen(5)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s-%s", label, randomString)
}

// Add creates a new Site Controller and adds it to the API Server
func Add(server *apiserver.APIServer) error {
	sites.RegisterSitesServiceServer(server.GRPCServer,
		NewSitesServiceServer(server.Manager.GetClient(), rest.CopyConfig(server.Manager.GetConfig())))
	return nil
}

func (s *sitesService) CreateSite(ctx context.Context, r *sites.CreateSiteRequest) (*sites.Site, error) {
	c, userID := impersonate.ClientFromContext(ctx, s.cfg)
	org := metadata.RequireOrganizationName(ctx)

	var siteName, projName string

	proj, err := project.Resolve(r.Parent)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}

	projNs, err := projectns.Lookup(s.client, proj, org)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, status.NotFoundf("project not found")
		}
		return nil, status.InternalError()
	}

	if len(r.Site.PrimaryDomain) == 0 {
		return nil, status.InvalidArgumentf("primary domain cannot be empty")
	}

	if len(r.Name) == 0 {
		r.Name = dashboardsite.FQName(proj, generateNameFromDomain(r.PrimaryDomain))
	}
	if !strings.HasPrefix(r.Name, r.Parent) {
		return nil, status.InvalidArgumentf("parent and project are not matching")
	}
	siteName, projName, err = site.Resolve(r.Name)
	if err != nil {
		return nil, err
	}
	if proj != projName {
		return nil, status.InvalidArgumentf("parent and project are not matching")
	}

	wp := dashboardsite.New(&wordpressv1alpha1.Wordpress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      siteName,
			Namespace: projNs.Name,
			Labels: map[string]string{
				"presslabs.com/kind":         "site",
				"presslabs.com/site":         siteName,
				"presslabs.com/project":      proj,
				"presslabs.com/organization": org,
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

	return newSiteFromK8s(wp, []*sites.Endpoint{}), nil
}

func (s *sitesService) GetSite(ctx context.Context, r *sites.GetSiteRequest) (*sites.Site, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)
	org := metadata.RequireOrganizationName(ctx)

	key, err := site.ResolveToObjectKey(s.client, r.Name, org)
	if err != nil {
		return nil, err
	}

	wp := dashboardsite.New(&wordpressv1alpha1.Wordpress{})
	if err = c.Get(ctx, *key, wp.Unwrap()); err != nil {
		return nil, status.NotFoundf("site %s not found", r.Name).Because(err)
	}

	endp, err := getEndpoints(s.client, wp)
	if err != nil {
		return nil, err
	}

	return newSiteFromK8s(wp, endp), nil
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
func updateWordpressImage(wp *wordpressv1alpha1.Wordpress, image string, fieldMask types.FieldMask) error { // nolint:unparam
	if len(fieldMask.Paths) == 0 || containsString(fieldMask.Paths, "site.wordpress_image") {
		wp.Spec.Image = image
	}
	return nil
}

func (s *sitesService) UpdateSite(ctx context.Context, r *sites.UpdateSiteRequest) (*sites.Site, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)
	org := metadata.RequireOrganizationName(ctx)

	wp := dashboardsite.New(&wordpressv1alpha1.Wordpress{})
	key, err := site.ResolveToObjectKey(s.client, r.Name, org)
	if err != nil {
		return nil, err
	}

	// get the site
	if err = c.Get(ctx, *key, wp.Unwrap()); err != nil {
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

	endp, err := getEndpoints(s.client, wp)
	if err != nil {
		return nil, err
	}

	return newSiteFromK8s(wp, endp), nil
}

func (s *sitesService) DeleteSite(ctx context.Context, r *sites.DeleteSiteRequest) (*types.Empty, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)
	org := metadata.RequireOrganizationName(ctx)

	key, err := site.ResolveToObjectKey(s.client, r.Name, org)
	if err != nil {
		return nil, err
	}

	wp := dashboardsite.New(&wordpressv1alpha1.Wordpress{})
	if err := c.Get(ctx, *key, wp.Unwrap()); err != nil {
		return nil, status.NotFound().Because(err)
	}

	if err := c.Delete(ctx, wp.Unwrap()); err != nil {
		return nil, status.FromError(err)
	}

	return &types.Empty{}, nil
}

func (s *sitesService) ListSites(ctx context.Context, r *sites.ListSitesRequest) (*sites.ListSitesResponse, error) {
	c, _ := impersonate.ClientFromContext(ctx, s.cfg)
	org := metadata.RequireOrganizationName(ctx)

	wpList := &wordpressv1alpha1.WordpressList{}
	resp := &sites.ListSitesResponse{}

	proj, err := project.Resolve(r.Parent)
	if err != nil {
		return nil, status.InvalidArgumentf("%s", err)
	}

	projNs, err := projectns.Lookup(s.client, proj, org)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, status.NotFoundf("project not found")
		}
		return nil, status.InternalError()
	}

	listOptions := &client.ListOptions{
		Namespace: projNs.Name,
	}

	if err := c.List(context.TODO(), listOptions, wpList); err != nil {
		return nil, status.FromError(err)
	}

	// TODO: implement pagination
	resp.Sites = make([]sites.Site, len(wpList.Items))
	for i := range wpList.Items {
		endp, err := getEndpoints(s.client, dashboardsite.New(&wpList.Items[i]))
		if err != nil {
			return nil, err
		}
		resp.Sites[i] = *newSiteFromK8s(dashboardsite.New(&wpList.Items[i]), endp)
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

func getEndpoints(cl client.Client, s *dashboardsite.Site) ([]*sites.Endpoint, error) {
	ingr := extv1beta1.Ingress{}
	key := client.ObjectKey{
		Name:      s.Unwrap().Name,
		Namespace: s.Unwrap().Namespace,
	}
	if err := cl.Get(context.TODO(), key, &ingr); err != nil {
		return nil, status.FromError(err)
	}

	endp := make([]*sites.Endpoint, len(ingr.Status.LoadBalancer.Ingress))
	for i := range ingr.Status.LoadBalancer.Ingress {
		endp[i].Ip = ingr.Status.LoadBalancer.Ingress[i].IP
		endp[i].Host = ingr.Status.LoadBalancer.Ingress[i].Hostname
	}

	return endp, nil
}

func newSiteFromK8s(s *dashboardsite.Site, endpoints []*sites.Endpoint) *sites.Site {
	return &sites.Site{
		Name:           dashboardsite.FQName(s.ObjectMeta.Labels["presslabs.com/project"], s.Name),
		PrimaryDomain:  string(s.Spec.Domains[0]),
		WordpressImage: s.Spec.Image,
		Endpoints:      endpoints,
	}
}
