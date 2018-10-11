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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/rand"
	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/site"
)

// NewMysqlClusterSecretSyncer returns a new syncer.Interface for reconciling MysqlCluster Secret
func NewMysqlClusterSecretSyncer(wp *site.Site, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := wp.ComponentLabels(site.MysqlClusterSecret)

	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wp.ComponentName(site.MysqlClusterSecret),
			Namespace: wp.Namespace,
		},
	}

	return syncer.NewObjectSyncer("MysqlCluster Secret", wp.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Secret)

		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		stringData := make(map[string]string)

		_, ok := out.Data["USER"]
		if !ok {
			stringData["USER"] = "wordpress"
		}

		_, ok = out.Data["DATABASE"]
		if !ok {
			stringData["DATABASE"] = "wordpress"
		}

		_, ok = out.Data["ROOT_PASSWORD"]
		if !ok {
			password, err := rand.AlphaNumericString(20)
			if err != nil {
				return err
			}
			stringData["ROOT_PASSWORD"] = password
		}

		_, ok = out.Data["PASSWORD"]
		if !ok {
			password, err := rand.AlphaNumericString(20)
			if err != nil {
				return err
			}
			stringData["PASSWORD"] = password
		}

		if len(stringData) > 0 {
			out.StringData = stringData
		}

		return nil
	})
}
