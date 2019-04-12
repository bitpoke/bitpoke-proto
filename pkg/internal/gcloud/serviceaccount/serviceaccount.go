/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package serviceaccount

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

// CreateServiceAccount creates a service account.
func CreateServiceAccount(projectID, name, displayName string) (*iam.ServiceAccount, error) {
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

// CreateServiceAccountKey creates a service account key.
func CreateServiceAccountKey(serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
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

// DeleteServiceAccount deletes a service account.
func DeleteServiceAccount(email string) error {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return fmt.Errorf("iam.New: %v", err)
	}

	_, err = service.Projects.ServiceAccounts.Delete("projects/-/serviceAccounts/" + email).Do()
	if err != nil {
		return fmt.Errorf("Projects.ServiceAccounts.Delete: %v", err)
	}
	return nil
}
