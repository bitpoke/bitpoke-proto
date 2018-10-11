/*
Copyright 2018 Pressinfra SRL.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sync_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/controller/site/internal/sync"
	"github.com/presslabs/dashboard/pkg/internal/site"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

func getEnvVar(env []corev1.EnvVar, name string) corev1.EnvVar {
	for _, e := range env {
		if e.Name == name {
			return e
		}
	}
	return corev1.EnvVar{}
}

var _ = Describe("WordpressSyncer", func() {
	var wp *site.Site

	BeforeEach(func() {
		wp = site.New(&wordpressv1alpha1.Wordpress{ObjectMeta: metav1.ObjectMeta{Name: "wp", Namespace: "default"}})
		wpSyncer := sync.NewWordpressSyncer(wp, fake.NewFakeClient(), scheme.Scheme).(*syncer.ObjectSyncer)

		Expect(wpSyncer.SyncFn(wp.Unwrap())).To(Succeed())
	})

	DescribeTable("when syncing", func(name string, value interface{}) {
		switch v := value.(type) {
		case string:
			Expect(getEnvVar(wp.Spec.Env, name)).To(MatchFields(IgnoreExtras, Fields{
				"Value": Equal(v),
			}))
		case *corev1.EnvVarSource:
			Expect(getEnvVar(wp.Spec.Env, name)).To(MatchFields(IgnoreExtras, Fields{
				"ValueFrom": Equal(v),
			}))
		default:
			panic(fmt.Sprintf("%T is not a string nor a *v1.EnvVarSource", value))
		}
	},
		Entry("successfully sets MEMCACHED_DISCOVERY_SERVICE", "MEMCACHED_DISCOVERY_SERVICE", "wp-memcached.default"),
		Entry("successfully sets WORDPRESS_DB_HOST", "WORDPRESS_DB_HOST", "wp-mysql-master.default"),
		Entry("successfully sets WORDPRESS_DB_USER", "WORDPRESS_DB_USER", &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wp-mysql",
				},
				Key: "USER",
			},
		}),
		Entry("successfully sets WORDPRESS_DB_PASSWORD", "WORDPRESS_DB_PASSWORD", &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wp-mysql",
				},
				Key: "PASSWORD",
			},
		}),
		Entry("successfully sets WORDPRESS_DB_NAME", "WORDPRESS_DB_NAME", &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "wp-mysql",
				},
				Key: "DATABASE",
			},
		}),
	)
})
