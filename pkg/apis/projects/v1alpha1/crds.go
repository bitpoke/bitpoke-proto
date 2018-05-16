/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package v1alpha1

import (
	"fmt"

	kutil "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	projectsOpenApi "github.com/presslabs/dashboard/pkg/openapi"
)

const (
	projectsApiPkg = "github.com/presslabs/dashboard/pkg/apis/projects"
)

// Project Custom Resource Definition
var (
	// ResourceProject contains the definition bits for Project CRD
	ResourceProject = kutil.Config{
		Group:   SchemeGroupVersion.Group,
		Version: SchemeGroupVersion.Version,

		Kind:       ResourceKindProject,
		Plural:     "projects",
		Singular:   "project",
		ShortNames: []string{"proj"},

		SpecDefinitionName:    fmt.Sprintf("%s/%s.%s", projectsApiPkg, SchemeGroupVersion.Version, ResourceKindProject),
		ResourceScope:         string(apiextensions.NamespaceScoped),
		GetOpenAPIDefinitions: projectsOpenApi.GetOpenAPIDefinitions,

		EnableValidation:        true,
		EnableStatusSubresource: true,
	}
	// ResourceProjectCRDName is the fully qualified Project CRD name (ie. projects.dashboard.presslabs.com)
	ResourceProjectCRDName = fmt.Sprintf("%s.%s", ResourceProject.Plural, ResourceProject.Group)
	// ResourceProjectCRD is the Custrom Resource Definition object for Project
	ResourceProjectCRD = kutil.NewCustomResourceDefinition(ResourceProject)
)

var CRDs = map[string]kutil.Config{
	ResourceProjectCRDName: ResourceProject,
}
