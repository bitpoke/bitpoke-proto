/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package organization

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	// logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	orgv1 "github.com/presslabs/dashboard/pkg/api/organizations/v1"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	"github.com/presslabs/dashboard/pkg/internal/organization"
	. "github.com/presslabs/dashboard/pkg/internal/testutil/gomega"
)

const (
	ctxTimeout    = time.Second * 3
	deleteTimeout = time.Second
	updateTimeout = time.Second
)

// createOrganization is a helper func that creates an organization
func createOrganization(name, displayName, createdBy string) *organization.Organization {
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
	org.UpdateDisplayName(displayName)

	return org
}

// getNamespaceFn is a helper func that returns an organization
func getNamespaceFn(ctx context.Context, c client.Client, key client.ObjectKey) func() corev1.Namespace {
	return func() corev1.Namespace {
		var orgNs corev1.Namespace
		Expect(c.Get(ctx, key, &orgNs)).To(Succeed())
		return orgNs
	}
}

// var log = logf.Log.WithName("apiserver")

var _ = Describe("API server", func() {
	var (
		// stop channel for apiserver
		stop chan struct{}
		// controller k8s client
		c client.Client
		// client connection to an RPC server
		conn *grpc.ClientConn
		// orgClient
		orgClient orgv1.OrganizationsServiceClient
	)

	BeforeEach(func() {
		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).To(Succeed())

		server := SetupAPIServer(mgr)
		// add ourselves to the server
		Add(server)

		c = mgr.GetClient()

		stop = StartTestManager(mgr)

		conn, err = grpc.Dial(server.GetGRPCAddr(), grpc.WithInsecure(), grpc.WithBlock(),
			grpc.WithTimeout(ctxTimeout))
		Expect(err).To(Succeed())

		orgClient = orgv1.NewOrganizationsServiceClient(conn)
	})

	AfterEach(func() {
		// close the gRPC client connection
		conn.Close()
		// stop the manager and API server
		close(stop)
	})

	var (
		id          string
		name        string
		displayName string
		createdBy   string
	)

	Describe("at Create request", func() {
		When("organization already exists", func() {
			BeforeEach(func() {
				id = fmt.Sprintf("%d", rand.Int31())
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())

				org := createOrganization(name, displayName, createdBy)
				Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			})

			It("should returns error", func() {
				req := orgv1.CreateOrganizationRequest{
					OrganizationId: id,
					Organization: &orgv1.Organization{
						Name:        name,
						DisplayName: displayName,
					},
				}

				_, err := orgClient.CreateOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Code()).To(Equal(codes.AlreadyExists))
			})
		})

		When("organization does not exists", func() {
			BeforeEach(func() {
				id = fmt.Sprintf("%d", rand.Int31())
				name = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())
				middleware.FakeSubject = createdBy
			})

			entries := []TableEntry{
				Entry("should creates the organization and set the given display-name", "display-name"),
				Entry("sdould creates the organization and set the default display-name", ""),
			}

			DescribeTable("", func(dispName string) {
				expectedDispName := dispName
				if len(dispName) == 0 {
					expectedDispName = name
				}
				req := orgv1.CreateOrganizationRequest{
					OrganizationId: id,
					Organization: &orgv1.Organization{
						Name:        name,
						DisplayName: dispName,
					},
				}

				resp, err := orgClient.CreateOrganization(context.TODO(), &req)
				Expect(err).To(Succeed())
				Expect(resp.Name).To(Equal(name))

				var orgNs corev1.Namespace
				key := client.ObjectKey{
					Name: organization.NamespaceName(name),
				}
				err = c.Get(context.TODO(), key, &orgNs)
				Expect(err).To(Succeed())
				Expect(orgNs.ObjectMeta.Annotations).To(HaveKeyWithValue("presslabs.com/display-name", expectedDispName))
				Expect(orgNs.ObjectMeta.Annotations).To(HaveKeyWithValue("presslabs.com/created-by", createdBy))
			}, entries...)
		})
	})

	Describe("at Get request", func() {
		When("organization exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())

				org := createOrganization(name, displayName, createdBy)
				Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			})

			It("should returns the organization", func() {
				req := orgv1.GetOrganizationRequest{
					Name: name,
				}

				resp, err := orgClient.GetOrganization(context.TODO(), &req)
				Expect(err).To(Succeed())
				Expect(resp.Name).To(Equal(name))
				Expect(resp.DisplayName).To(Equal(displayName))
			})
		})

		When("organization does not exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
			})

			It("should return error", func() {
				req := orgv1.GetOrganizationRequest{
					Name: name,
				}
				_, err := orgClient.GetOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
			})
		})
	})

	Describe("at Delete request", func() {
		When("organization exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())

				org := createOrganization(name, displayName, createdBy)
				err := c.Create(context.TODO(), org.Unwrap())
				Expect(err).To(Succeed())
			})

			It("should delete the organization", func() {
				req := orgv1.DeleteOrganizationRequest{
					Name: name,
				}
				_, err := orgClient.DeleteOrganization(context.TODO(), &req)
				Expect(err).To(Succeed())

				key := client.ObjectKey{
					Name: organization.NamespaceName(name),
				}

				Eventually(getNamespaceFn(context.TODO(), c, key), deleteTimeout).Should(
					BeInPhase(corev1.NamespaceTerminating))
			})
		})

		When("organization does not exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
			})

			It("should returns error", func() {
				req := orgv1.DeleteOrganizationRequest{
					Name: name,
				}
				_, err := orgClient.DeleteOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
			})
		})
	})

	Describe("at Update request", func() {
		When("organization exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
				// newDisplayName = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())

				org := createOrganization(name, displayName, createdBy)
				Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			})

			entries := []TableEntry{
				Entry("should update the display name with given value", "new-display-name"),
				Entry("should update the display name with default value", ""),
			}

			DescribeTable("", func(newDisplayName string) {
				expectedDisplayName := newDisplayName
				if newDisplayName == "" {
					expectedDisplayName = name
				}

				req := orgv1.UpdateOrganizationRequest{
					Organization: &orgv1.Organization{
						Name:        name,
						DisplayName: newDisplayName,
					},
				}
				_, err := orgClient.UpdateOrganization(context.TODO(), &req)
				Expect(err).To(Succeed())

				key := client.ObjectKey{
					Name: organization.NamespaceName(name),
				}

				Eventually(getNamespaceFn(context.TODO(), c, key), updateTimeout).Should(
					HaveAnnotation("presslabs.com/display-name", expectedDisplayName))
			}, entries...)
		})

		When("organization does not exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
			})

			It("should return error", func() {
				req := orgv1.UpdateOrganizationRequest{
					Organization: &orgv1.Organization{
						Name:        name,
						DisplayName: displayName,
					},
				}
				_, err := orgClient.UpdateOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
			})
		})
	})

	Describe("at List request", func() {
		When("should organizations exist", func() {
			var (
				orgList   []*organization.Organization
				pageToken string
				pageSize  int32
				noItems   int
			)

			BeforeEach(func() {
				pageToken = ""
				pageSize = int32(3)
				noItems = 3
				createdBy = fmt.Sprintf("%d", rand.Int31())

				middleware.FakeSubject = createdBy

				orgList = make([]*organization.Organization, noItems)

				for i := 0; i < noItems; i++ {
					name = fmt.Sprintf("%d", rand.Int31())
					displayName := fmt.Sprintf("%d", rand.Int31())

					org := createOrganization(name, displayName, createdBy)
					Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())

					orgList[i] = org
				}
			})

			It("list organizations", func() {
				req := orgv1.ListOrganizationsRequest{
					PageToken: pageToken,
					PageSize:  pageSize,
				}

				resp, err := orgClient.ListOrganizations(context.TODO(), &req)
				Expect(err).To(Succeed())
				Expect(resp.Organizations).To(HaveLen(noItems))
				// TODO: check pagination when it's implemented
			})
		})
	})
})
