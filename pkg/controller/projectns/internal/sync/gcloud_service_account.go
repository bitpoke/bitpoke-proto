/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

// createServiceAccount creates a service account.
func createServiceAccount(projectID, name, displayName string) (*iam.ServiceAccount, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	request := &iam.CreateServiceAccountRequest{
		AccountId: name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: displayName,
		},
	}
	account, err := service.Projects.ServiceAccounts.Create("projects/"+projectID, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Create: %v", err)
	}
	return account, nil
}

// createKey creates a service account key.
func createKey(serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	request := &iam.CreateServiceAccountKeyRequest{}
	key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
	}
	return key, nil
}

// NewGCloudServiceAccountSyncer returns a new syncer.Interface for reconciling gcloud service account
func NewGCloudServiceAccountSyncer(proj *projectns.ProjectNamespace, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := labels.Set{
		"presslabs.com/kind": "gcloud-service-account-secret",
	}

	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(projectns.GcloudServiceAccountSecret),
			Namespace: proj.ComponentName(projectns.Namespace),
		},
	}

	return syncer.NewObjectSyncer("GcloudServiceAccount", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Secret)

		if out.CreationTimestamp.IsZero() {
			sa, err := createServiceAccount(options.GCloudProjectID,
				proj.ObjectMeta.Labels["presslabs.com/project"],
				proj.ObjectMeta.Annotations["presslabs.com/display-name"])
			if err != nil {
				return err
			}

			saKey, err := createKey(sa.Email)
			if err != nil {
				return err
			}

			secretData, err := saKey.MarshalJSON()
			if err != nil {
				return err
			}

			out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)
			out.Data = map[string][]byte{
				"SERVICE_ACCOUNT_KEY":  secretData,
				"SERVICE_ACCOUNT_MAIL": []byte(sa.Email),
			}
		}

		return nil
	})
}
