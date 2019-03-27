/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package site

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	. "github.com/presslabs/dashboard/pkg/internal/testutil/gomega"

	sites "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/sites/v1"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/controller"
	"github.com/presslabs/dashboard/pkg/internal/site"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const (
	ctxTimeout = time.Second * 3
)

// createSite is a helper func that creates a site
func createSite(name, userID, project, image string, domains []wordpressv1alpha1.Domain) *site.Site {
	wp := site.New(&wordpressv1alpha1.Wordpress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: project,
			Labels: map[string]string{
				"presslabs.com/kind":    "site",
				"presslabs.com/site":    name,
				"presslabs.com/project": project,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": userID,
			},
		},
		Spec: wordpressv1alpha1.WordpressSpec{
			Domains: domains,
			Image:   image,
		},
	})

	return wp
}

func expectProperWordpress(c client.Client, name, userID, project, image string, domains []wordpressv1alpha1.Domain) {
	wp := &wordpressv1alpha1.Wordpress{}
	key := client.ObjectKey{
		Name:      name,
		Namespace: project,
	}
	Expect(c.Get(context.TODO(), key, wp)).To(Succeed())
	Expect(wp.Name).To(Equal(name))
	Expect(wp.Labels).To(HaveKeyWithValue("presslabs.com/kind", "site"))
	Expect(wp.Labels).To(HaveKeyWithValue("presslabs.com/site", name))
	Expect(wp.Labels).To(HaveKeyWithValue("presslabs.com/project", project))
	Expect(wp.Annotations).To(HaveKeyWithValue("presslabs.com/created-by", userID))
	Expect(wp.Spec.Image).To(Equal(image))
	Expect(wp.Spec.Domains).To(Equal(domains))
}

var _ = Describe("API server", func() {
	var (
		// stop channel for apiserver
		stop chan struct{}
		// controller k8s client
		c client.Client
		// client connection to an RPC server
		conn *grpc.ClientConn
		// siteClient
		siteClient sites.SitesServiceClient
		// context for requests
		ctx context.Context
	)

	var (
		id, name             string
		userID               string
		project, parent      string
		organization         string
		image, primaryDomain string
		domains              []wordpressv1alpha1.Domain
		ns                   *corev1.Namespace
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

		siteClient = sites.NewSitesServiceClient(conn)

		project = fmt.Sprintf("%d", rand.Int31())
		name = fmt.Sprintf("%d", rand.Int31())
		id = fmt.Sprintf("project/%s/site/%s", project, name)
		userID = fmt.Sprintf("user#%s", name)
		parent = fmt.Sprintf("project/%s", project)
		image = fmt.Sprintf("%d", rand.Int31())
		metadata.FakeSubject = userID
		primaryDomain = fmt.Sprintf("%d", rand.Int31())
		domains = []wordpressv1alpha1.Domain{
			wordpressv1alpha1.Domain(primaryDomain),
		}
		organization = fmt.Sprintf("%d", rand.Int31())

		ctx = metadata.AddOrgInContext(context.Background(), organization)

		// create ProjectNamespace
		ns = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      project,
				Namespace: organization,
				Labels: map[string]string{
					"presslabs.com/kind":         "project",
					"presslabs.com/project":      project,
					"presslabs.com/organization": organization,
				},
				Annotations: map[string]string{
					"presslabs.com/created-by": userID,
				},
			},
		}
		Expect(c.Create(context.TODO(), ns)).To(Succeed())
	})

	AfterEach(func() {
		// close the gRPC client connection
		conn.Close()
		// stop the manager and API server
		close(stop)

		// delete ProjectNamespace
		Expect(c.Delete(context.TODO(), ns)).To(Succeed())
		Eventually(func() corev1.Namespace {
			ns := corev1.Namespace{}
			c.Get(context.TODO(), client.ObjectKey{Name: project, Namespace: organization}, &ns)
			return ns
		}).Should(BeInPhase(corev1.NamespaceTerminating))
	})

	Describe("at Create request", func() {
		It("returns AlreadyExists error when site already exists", func() {
			wp := createSite(name, userID, project, image, domains)
			Expect(c.Create(context.TODO(), wp.Unwrap())).To(Succeed())
			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}

			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.AlreadyExists))
		})

		It("returns error when parent and project are not matching", func() {
			req := sites.CreateSiteRequest{
				Parent: "another-project",
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}

			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when no parent is given", func() {
			req := sites.CreateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}

			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when parent is not fully-qualified", func() {
			req := sites.CreateSiteRequest{
				Parent: "not-fully-qualified parent",
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}

			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("creates site and generate a name when no one is given", func() {
			primaryDomain := "www.presslabs.com"
			domains[0] = wordpressv1alpha1.Domain(primaryDomain)

			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			wp, err := siteClient.CreateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(wp.PrimaryDomain).To(Equal(primaryDomain))
			Expect(wp.WordpressImage).To(Equal(image))

			// check generated name for the site
			Expect(wp.Name).Should(Not(BeEmpty()))
			generatedName, _, err := site.Resolve(wp.Name)
			Expect(err).To(Succeed())
			Expect(strings.HasPrefix(generatedName, "presslabs")).To(BeTrue())

			expectProperWordpress(c, generatedName, userID, project, image, domains)
		})

		It("returns error when name is not fully qualified", func() {
			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					Name:           "not-fully-qualified-name",
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when project name is empty", func() {
			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					Name:           "project//site/site-name",
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when site name is empty", func() {
			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					Name:           "project/project-name/site/",
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when domain is empty", func() {
			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  "",
					WordpressImage: image,
				},
			}
			_, err := siteClient.CreateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when no organization is set in metadata", func() {
			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			_, err := siteClient.CreateSite(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})

		It("creates site", func() {
			req := sites.CreateSiteRequest{
				Parent: parent,
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			wp, err := siteClient.CreateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(wp.Name).To(Equal(id))
			Expect(wp.PrimaryDomain).To(Equal(primaryDomain))
			Expect(wp.WordpressImage).To(Equal(image))
			expectProperWordpress(c, name, userID, project, image, domains)
		})
	})

	Describe("at Get request", func() {
		It("returns NotFound when site does not exist", func() {
			req := sites.GetSiteRequest{
				Name: id,
			}
			_, err := siteClient.GetSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("returns error when no organization is set in metadata", func() {
			wp := createSite(name, userID, project, image, domains)
			Expect(c.Create(context.TODO(), wp.Unwrap())).To(Succeed())
			req := sites.GetSiteRequest{
				Name: id,
			}

			_, err := siteClient.GetSite(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})

		It("returns the site", func() {
			wp := createSite(name, userID, project, image, domains)
			Expect(c.Create(context.TODO(), wp.Unwrap())).To(Succeed())
			req := sites.GetSiteRequest{
				Name: id,
			}

			resp, err := siteClient.GetSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(primaryDomain))
			Expect(resp.WordpressImage).To(Equal(image))
		})
	})

	Describe("at Delete request", func() {
		It("returns NotFound when site does not exists", func() {
			req := sites.DeleteSiteRequest{
				Name: id,
			}

			_, err := siteClient.DeleteSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("returns error when no organization is set in metadata", func() {
			wp := createSite(name, userID, project, image, domains)
			Expect(c.Create(context.TODO(), wp.Unwrap())).To(Succeed())
			req := sites.DeleteSiteRequest{
				Name: id,
			}

			_, err := siteClient.DeleteSite(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})

		It("deletes existing site", func() {
			wp := createSite(name, userID, project, image, domains)
			Expect(c.Create(context.TODO(), wp.Unwrap())).To(Succeed())
			req := sites.DeleteSiteRequest{
				Name: id,
			}

			_, err := siteClient.DeleteSite(ctx, &req)
			Expect(err).To(Succeed())

			var deletedWp wordpressv1alpha1.Wordpress
			key := client.ObjectKey{
				Name:      name,
				Namespace: project,
			}
			err = c.Get(context.TODO(), key, &deletedWp)
			Expect(errors.IsNotFound(err)).To(Equal(true))
		})
	})

	Describe("at Update request", func() {
		BeforeEach(func() {
			wp := createSite(name, userID, project, image, domains)
			Expect(c.Create(context.TODO(), wp.Unwrap())).To(Succeed())
			time.Sleep(time.Millisecond * 500)
		})

		It("returns NotFound when site does not exist", func() {
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           fmt.Sprintf("project/%s/site/inexistent", project),
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			_, err := siteClient.UpdateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("updates the wordpress image of existing site when 'site.wordpress_image' is in r.UpdateMask.GetPaths()", func() {
			newImage := "new-wordpress-image"
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: newImage,
				},
				FieldMask: types.FieldMask{
					Paths: []string{"site.wordpress_image"},
				},
			}

			resp, err := siteClient.UpdateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(primaryDomain))
			Expect(resp.WordpressImage).To(Equal(newImage))
			expectProperWordpress(c, name, userID, project, newImage, domains)
		})

		It("keeps the old value of the wordpress image when 'site.wordpress_image' is not in r.UpdateMask.GetPaths() and r.UpdateMask.GetPaths() is not empty", func() {
			newImage := "new-wordpress-image"
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: newImage,
				},
				FieldMask: types.FieldMask{
					Paths: []string{"site.primary_domain"},
				},
			}

			resp, err := siteClient.UpdateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(primaryDomain))
			Expect(resp.WordpressImage).To(Equal(image))
			expectProperWordpress(c, name, userID, project, image, domains)
		})

		It("updates the primary domain of existing site when 'site.primary_domain' is in r.UpdateMask.GetPaths()", func() {
			newPD := "new-primary-domain"
			expectedDomains := domains
			expectedDomains[0] = wordpressv1alpha1.Domain(newPD)
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  newPD,
					WordpressImage: image,
				},
				FieldMask: types.FieldMask{
					Paths: []string{"site.primary_domain"},
				},
			}

			resp, err := siteClient.UpdateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(newPD))
			Expect(resp.WordpressImage).To(Equal(image))
			expectProperWordpress(c, name, userID, project, image, expectedDomains)
		})

		It("updates only primary domain when the site have more domains", func() {
			name := fmt.Sprintf("%d", rand.Int31())
			id = fmt.Sprintf("project/%s/site/%s", project, name)
			for i := 0; i < 3; i++ {
				domains = append(domains, wordpressv1alpha1.Domain(fmt.Sprintf("%s-%02d", primaryDomain, i)))
			}
			wp := createSite(name, userID, project, image, domains)
			Expect(c.Create(context.TODO(), wp.Unwrap())).To(Succeed())

			newPD := "new-primary-domain"
			expectedDomains := domains
			expectedDomains[0] = wordpressv1alpha1.Domain(newPD)

			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  newPD,
					WordpressImage: image,
				},
				FieldMask: types.FieldMask{
					Paths: []string{"site.primary_domain"},
				},
			}

			resp, err := siteClient.UpdateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(newPD))
			Expect(resp.WordpressImage).To(Equal(image))
			expectProperWordpress(c, name, userID, project, image, expectedDomains)
		})

		It("keeps the old value of the primary domain when 'site.primary_domain' is not in r.UpdateMask.GetPaths() and r.UpdateMask.GetPaths() is not empty", func() {
			newDomain := "new-primary-domain"
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  newDomain,
					WordpressImage: image,
				},
				FieldMask: types.FieldMask{
					Paths: []string{"site.wordpress_image"},
				},
			}

			resp, err := siteClient.UpdateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(primaryDomain))
			Expect(resp.WordpressImage).To(Equal(image))
			expectProperWordpress(c, name, userID, project, image, domains)
		})

		It("returns error when 'site.primary_domain' is in r.UpdateMask.GetPaths() and primaryDomain field is empty", func() {
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					WordpressImage: image,
				},
				FieldMask: types.FieldMask{
					Paths: []string{"site.primary_domain"},
				},
			}
			_, err := siteClient.UpdateSite(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("updates primary domain and wordpress image of existing site when 'site.primary_domain' and `site.wordpress_image` are in r.UpdateMask.GetPaths()", func() {
			newPD := "new-primary-domain"
			newWI := "new-wordpress-image"
			expectedDomains := domains
			expectedDomains[0] = wordpressv1alpha1.Domain(newPD)
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  newPD,
					WordpressImage: newWI,
				},
				FieldMask: types.FieldMask{
					Paths: []string{"site.primary_domain", "site.wordpress_image"},
				},
			}

			resp, err := siteClient.UpdateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(newPD))
			Expect(resp.WordpressImage).To(Equal(newWI))
			expectProperWordpress(c, name, userID, project, newWI, expectedDomains)
		})

		It("updates primary domain and wordpress image of existing site when r.UpdateMask.GetPaths() is empty", func() {
			newPD := "new-primary-domain"
			newWI := "new-wordpress-image"
			expectedDomains := domains
			expectedDomains[0] = wordpressv1alpha1.Domain(newPD)
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  newPD,
					WordpressImage: newWI,
				},
			}

			resp, err := siteClient.UpdateSite(ctx, &req)
			Expect(err).To(Succeed())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.PrimaryDomain).To(Equal(newPD))
			Expect(resp.WordpressImage).To(Equal(newWI))
			expectProperWordpress(c, name, userID, project, newWI, expectedDomains)
		})

		It("returns error when no organization is set in metadata", func() {
			req := sites.UpdateSiteRequest{
				Site: sites.Site{
					Name:           id,
					PrimaryDomain:  primaryDomain,
					WordpressImage: image,
				},
			}
			_, err := siteClient.UpdateSite(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("at list request", func() {
		var sitesCount = 3
		BeforeEach(func() {
			for i := 1; i <= sitesCount; i++ {
				_name := fmt.Sprintf("%s-%02d", name, i)
				_image := fmt.Sprintf("%s-%02d", image, i)
				_primaryDomain := fmt.Sprintf("%s-%02d", primaryDomain, i)
				_domains := []wordpressv1alpha1.Domain{
					wordpressv1alpha1.Domain(_primaryDomain),
				}
				site := createSite(_name, userID, project, _image, _domains)
				Expect(c.Create(context.TODO(), site.Unwrap())).To(Succeed())
			}
			site := createSite(name, "user#another", project, image, domains)
			Expect(c.Create(context.TODO(), site.Unwrap())).To(Succeed())

			name := fmt.Sprintf("%s", name)
			site = createSite(name, userID, "another-project", image, domains)
			Expect(c.Create(context.TODO(), site.Unwrap())).To(Succeed())
		})

		It("returns only my sites", func() {
			req := sites.ListSitesRequest{
				Parent: parent,
			}
			Eventually(func() ([]sites.Site, error) {
				resp, err := siteClient.ListSites(ctx, &req)
				return resp.Sites, err
			}).Should(HaveLen(sitesCount + 1))
		})

		It("returns error when no parent is given", func() {
			req := sites.ListSitesRequest{}
			_, err := siteClient.ListSites(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when parent is not fully-qualified", func() {
			req := sites.ListSitesRequest{
				Parent: "not-fully-qualified",
			}
			_, err := siteClient.ListSites(ctx, &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when no organization is set in metadata", func() {
			req := sites.ListSitesRequest{
				Parent: parent,
			}
			_, err := siteClient.ListSites(context.TODO(), &req)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})
})
