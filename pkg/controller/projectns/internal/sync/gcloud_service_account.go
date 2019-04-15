/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
	"github.com/presslabs/dashboard/pkg/internal/gcloud/serviceaccount"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

// NewGCloudServiceAccountSyncer returns a new syncer.Interface for reconciling gcloud service account
func NewGCloudServiceAccountSyncer(proj *projectns.ProjectNamespace, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(projectns.GcloudServiceAccountSecret),
			Namespace: proj.ComponentName(projectns.Namespace),
		},
	}

	return syncer.NewObjectSyncer("GcloudServiceAccount", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Secret)

		if out.CreationTimestamp.IsZero() {
			sa, err := serviceaccount.CreateServiceAccount(options.GCloudProjectID,
				proj.ObjectMeta.Labels["presslabs.com/project"],
				proj.ObjectMeta.Annotations["presslabs.com/display-name"])
			if err != nil {
				return err
			}

			saKey, err := serviceaccount.CreateServiceAccountKey(sa.Email)
			if err != nil {
				return err
			}

			secretData, err := saKey.MarshalJSON()
			if err != nil {
				return err
			}

			out.Labels = labels.Merge(out.Labels, controllerLabels)
			out.Data = map[string][]byte{
				"SERVICE_ACCOUNT_KEY":  secretData,
				"SERVICE_ACCOUNT_MAIL": []byte(sa.Email),
			}
		}

		return nil
	})
}
