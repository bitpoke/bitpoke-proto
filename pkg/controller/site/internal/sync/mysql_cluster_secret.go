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

package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/rand"
	"github.com/presslabs/controller-util/syncer"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

// NewMysqlClusterSecretSyncer returns a new syncer.Interface for reconciling MysqlCluster Secret
func NewMysqlClusterSecretSyncer(wp *wordpressv1alpha1.Wordpress, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlClusterName(wp),
			Namespace: wp.Namespace,
		},
	}

	return syncer.NewObjectSyncer("MysqlCluster Secret", wp, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Secret)

		out.StringData = map[string]string{
			"USER":     "wordpress",
			"DATABASE": "wordpress",
		}

		_, ok := out.Data["ROOT_PASSWORD"]
		if !ok {
			password, err := rand.AlphaNumericString(20)
			if err != nil {
				return err
			}
			out.StringData["ROOT_PASSWORD"] = password
		}

		_, ok = out.Data["PASSWORD"]
		if !ok {
			password, err := rand.AlphaNumericString(20)
			if err != nil {
				return err
			}
			out.StringData["PASSWORD"] = password
		}

		return nil
	})
}
