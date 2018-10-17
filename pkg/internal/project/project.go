/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package project

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// Project embeds dashboardv1alpha1.Project and adds utility functions
type Project struct {
	*dashboardv1alpha1.Project
}

var (
	// RequiredLabels is a list of required Project labels
	RequiredLabels = []string{"presslabs.com/organization", "presslabs.com/project"}
	// RequiredAnnotations is a list of required Project annotations
	RequiredAnnotations = []string{"presslabs.com/created-by"}
)

type component struct {
	name       string // eg. web, database, cache
	app        string // eg. mysql, memcached
	Version    string
	objNameFmt string
	objName    string
}

var (
	// Namespace component
	Namespace = component{objNameFmt: "proj-%s"}
	// LimitRange component
	LimitRange = component{objName: "presslabs-dashboard"}
	// ResourceQuota component
	ResourceQuota = component{objName: "presslabs-dashboard"}
	// Prometheus component
	Prometheus = component{app: "prometheus", objName: "prometheus", Version: "v2.3.2"}
	// GiteaDeployment component
	GiteaDeployment = component{name: "web", app: "gitea", objName: "gitea", Version: "1.5.2"}
	// GiteaService component
	GiteaService = component{name: "web", app: "gitea", objName: "gitea"}
	// GiteaIngress component
	GiteaIngress = component{name: "web", app: "gitea", objName: "gitea"}
	// GiteaPVC component
	GiteaPVC = component{name: "web", app: "gitea", objName: "gitea"}
	// GiteaSecret component
	GiteaSecret = component{name: "web", app: "gitea", objName: "gitea-conf"}
)

// New wraps a dashboardv1alpha1.Project into a Project object
func New(obj *dashboardv1alpha1.Project) *Project {
	return &Project{obj}
}

// Unwrap returns the wrapped dashboardv1alpha1.Project object
func (o *Project) Unwrap() *dashboardv1alpha1.Project {
	return o.Project
}

// Labels returns default label set for dashboardv1alpha1.Project
func (o *Project) Labels() labels.Set {
	labels := labels.Set{
		"presslabs.com/project": o.Name,
	}

	if o.ObjectMeta.Labels != nil {
		if org, ok := o.ObjectMeta.Labels["presslabs.com/organization"]; ok {
			labels["presslabs.com/organization"] = org
		}
	}

	return labels
}

// ComponentLabels returns labels for a label set for a dashboardv1alpha1.Project component
func (o *Project) ComponentLabels(component component) labels.Set {
	labels := o.Labels()
	if len(component.app) > 0 {
		labels["app.kubernetes.io/name"] = component.app
	}
	if len(component.name) > 0 {
		labels["app.kubernetes.io/component"] = component.name
	}
	if len(component.Version) > 0 {
		labels["app.kubernetes.io/version"] = component.Version
	}
	return labels
}

// ComponentName returns the object name for a component
func (o *Project) ComponentName(component component) string {
	if len(component.objNameFmt) == 0 {
		return component.objName
	}
	return fmt.Sprintf(component.objNameFmt, o.ObjectMeta.Name)
}

// Domain returns the project's subdomain label
func (o *Project) Domain() string {
	return o.Name
}

// ValidateMetadata validates the metadata of a Project
func (o *Project) ValidateMetadata() error {
	errorList := []error{}
	// Check for some required Project Labels and Annotations
	for _, label := range RequiredLabels {
		if value, exists := o.Project.Labels[label]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", label))
		}
	}
	for _, annotation := range RequiredAnnotations {
		if value, exists := o.Project.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}
