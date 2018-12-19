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
	. "github.com/onsi/gomega"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	// logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	orgv1 "github.com/presslabs/dashboard/pkg/api/organizations/v1"
	"github.com/presslabs/dashboard/pkg/apiserver/errors"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	apiserverutil "github.com/presslabs/dashboard/pkg/apiserver/util"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

const (
	ctxTimeout    = time.Second * 3
	deleteTimeout = time.Second
	updateTimeout = time.Second
)

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

		conn, err = grpc.Dial(server.GetGRPCAddr(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(ctxTimeout))
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

				org := organization.New(name, displayName, createdBy)
				Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			})

			It("returns error", func() {
				req := orgv1.CreateOrganizationRequest{
					OrganizationId: id,
					Organization: &orgv1.Organization{
						Name:        name,
						DisplayName: displayName,
					},
				}

				_, err := orgClient.CreateOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Message()).To(Equal(errors.AlreadyExists.Msg))
			})
		})

		When("organization not exists", func() {
			BeforeEach(func() {
				id = fmt.Sprintf("%d", rand.Int31())
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())

				createdBy = fmt.Sprintf("%d", rand.Int31())
				middleware.FakeSubject = createdBy
			})

			It("successfully creates an organization", func() {
				req := orgv1.CreateOrganizationRequest{
					OrganizationId: id,
					Organization: &orgv1.Organization{
						Name:        name,
						DisplayName: displayName,
					},
				}

				resp, err := orgClient.CreateOrganization(context.TODO(), &req)
				Expect(err).To(Succeed())
				Expect(resp.Name).To(Equal(name))

				// check org
				var orgNs corev1.Namespace
				key := client.ObjectKey{
					Name: organization.NamespaceName(name),
				}
				err = c.Get(context.TODO(), key, &orgNs)
				Expect(err).To(Succeed())
				Expect(orgNs.ObjectMeta.Annotations).To(HaveKeyWithValue("presslabs.com/display-name", displayName))
				Expect(orgNs.ObjectMeta.Annotations).To(HaveKeyWithValue("presslabs.com/created-by", createdBy))
			})
		})
	})

	Describe("at Get request", func() {
		When("organization exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())

				org := organization.New(name, displayName, createdBy)
				Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			})

			It("returns the organization", func() {
				req := orgv1.GetOrganizationRequest{
					Name: name,
				}

				resp, err := orgClient.GetOrganization(context.TODO(), &req)
				Expect(err).To(Succeed())
				Expect(resp.Name).To(Equal(name))
				Expect(resp.DisplayName).To(Equal(displayName))
			})
		})

		When("organization not exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
			})

			It("return error", func() {
				req := orgv1.GetOrganizationRequest{
					Name: name,
				}
				_, err := orgClient.GetOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Message()).To(Equal(errors.NotFound.Msg))
			})
		})
	})

	Describe("at Delete request", func() {
		When("organization exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())

				org := organization.New(name, displayName, createdBy)
				err := c.Create(context.TODO(), org.Unwrap())
				Expect(err).To(Succeed())
			})

			It("delete the organization", func() {
				req := orgv1.DeleteOrganizationRequest{
					Name: name,
				}
				_, err := orgClient.DeleteOrganization(context.TODO(), &req)
				Expect(err).To(Succeed())

				key := client.ObjectKey{
					Name: organization.NamespaceName(name),
				}

				Eventually(apiserverutil.GetNamespace(context.TODO(), c, key), deleteTimeout).Should(apiserverutil.BeInPhase(corev1.NamespaceTerminating))
			})
		})

		When("organization not exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
			})

			It("returns error", func() {
				req := orgv1.DeleteOrganizationRequest{
					Name: name,
				}
				_, err := orgClient.DeleteOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Message()).To(Equal(errors.NotFound.Msg))
			})
		})
	})

	Describe("at Update request", func() {
		When("organization exists", func() {
			var (
				newDisplayName string
			)

			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
				newDisplayName = fmt.Sprintf("%d", rand.Int31())
				createdBy = fmt.Sprintf("%d", rand.Int31())

				org := organization.New(name, displayName, createdBy)
				Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			})

			It("update the organization", func() {
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

				Eventually(apiserverutil.GetNamespace(context.TODO(), c, key), updateTimeout).Should(apiserverutil.HaveAnnotation("presslabs.com/display-name", newDisplayName))
			})
		})

		When("organization not exists", func() {
			BeforeEach(func() {
				name = fmt.Sprintf("%d", rand.Int31())
				displayName = fmt.Sprintf("%d", rand.Int31())
			})

			It("return error", func() {
				req := orgv1.UpdateOrganizationRequest{
					Organization: &orgv1.Organization{
						Name:        name,
						DisplayName: displayName,
					},
				}
				_, err := orgClient.UpdateOrganization(context.TODO(), &req)
				Expect(status.Convert(err).Message()).To(Equal(errors.NotFound.Msg))
			})
		})
	})

	Describe("at List request", func() {
		When("organizations exist", func() {
			var (
				orgList   []*organization.Organization
				pageToken string
				pageSize  int32
				noItems   int
				checked   []bool
				noChecked int
			)

			BeforeEach(func() {
				pageToken = ""
				pageSize = int32(3)
				noItems = 3
				checked = make([]bool, noItems)
				createdBy = fmt.Sprintf("%d", rand.Int31())

				middleware.FakeSubject = createdBy

				orgList = make([]*organization.Organization, noItems)

				for i := 0; i < noItems; i++ {
					name = fmt.Sprintf("%d", rand.Int31())
					displayName := fmt.Sprintf("%d", rand.Int31())

					org := organization.New(name, displayName, createdBy)
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
				Expect(len(resp.Organizations)).To(Equal(noItems))

				noChecked = 0
				for _, org := range resp.Organizations {
					for i, myOrg := range orgList {
						if organization.NamespaceName(org.Name) == myOrg.Unwrap().ObjectMeta.Name && !checked[i] {
							checked[i] = true
							noChecked++
							break
						}
					}
				}
				Expect(noChecked).To(Equal(noItems))
			})
		})
	})
})
