/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package v1alpha1

import (
	"fmt"

	kutilv1 "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/presslabs/dashboard/pkg/openapi"
)

const (
	projectsPkg = "github.com/presslabs/dashboard/pkg/apis/projects"
)

// Project Custom Resource Definition
var (
	// ResourceProject contains the definition bits for Project CRD
	ResourceProject = kutilv1.Config{
		Group:   SchemeGroupVersion.Group,
		Version: SchemeGroupVersion.Version,

		Kind:       ResourceKindProject,
		Plural:     "projects",
		Singular:   "project",
		ShortNames: []string{"proj"},

		SpecDefinitionName:    fmt.Sprintf("%s/%s.%s", projectsPkg, SchemeGroupVersion.Version, ResourceKindProject),
		ResourceScope:         string(apiextensionsv1.NamespaceScoped),
		GetOpenAPIDefinitions: openapi.GetOpenAPIDefinitions,

		EnableValidation:        true,
		EnableStatusSubresource: true,
	}
	// ResourceProjectCRDName is the fully qualified Project CRD name (ie. projects.dashboard.presslabs.com)
	ResourceProjectCRDName = fmt.Sprintf("%s.%s", ResourceProject.Plural, ResourceProject.Group)
	// ResourceProjectCRD is the Custrom Resource Definition object for Project
	ResourceProjectCRD = kutilv1.NewCustomResourceDefinition(ResourceProject)
)

var CRDs = map[string]kutilv1.Config{
	ResourceProjectCRDName: ResourceProject,
}
