/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Project CRUD", func() {
	var created *Project
	var key types.NamespacedName

	BeforeEach(func() {
		key = types.NamespacedName{Name: "foo", Namespace: "default"}
		created = &Project{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		}
	})

	AfterEach(func() {
		c.Delete(context.TODO(), created)
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {
				// c.List()
				fetched := &Project{}
				Expect(c.Create(context.TODO(), created)).NotTo(HaveOccurred())

				Expect(c.Get(context.TODO(), key, fetched)).NotTo(HaveOccurred())
				Expect(fetched).To(Equal(created))

				// Test Updating the Labels
				updated := fetched.DeepCopy()
				updated.Labels = map[string]string{"hello": "world"}
				Expect(c.Update(context.TODO(), updated)).NotTo(HaveOccurred())

				Expect(c.Get(context.TODO(), key, fetched)).NotTo(HaveOccurred())
				Expect(fetched).To(Equal(updated))

				// Test Delete
				Expect(c.Delete(context.TODO(), fetched)).NotTo(HaveOccurred())
				Expect(c.Get(context.TODO(), key, fetched)).To(HaveOccurred())
			})
		})
	})
})
