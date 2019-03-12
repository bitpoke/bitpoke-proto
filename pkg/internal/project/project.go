/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package project

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// Project embeds corev1.Namespace and adds utility functions
type Project struct {
	*dashboardv1alpha1.Project
}

var (
	// RequiredLabels is a list of required Project labels
	RequiredLabels = []string{"presslabs.com/organization", "presslabs.com/project", "presslabs.com/kind"}
	// RequiredAnnotations is a list of required Project annotations
	RequiredAnnotations = []string{"presslabs.com/created-by"}
)

type component struct {
	name       string // eg. web, database, cache
	app        string // eg. mysql, memcached
	objNameFmt string
	objName    string
}

// UpdateDisplayName updates the display-name annotation
func (p *Project) UpdateDisplayName(displayName string) {
	if len(displayName) == 0 {
		p.ObjectMeta.Annotations["presslabs.com/display-name"] = p.ObjectMeta.Labels["presslabs.com/project"]
	} else {
		p.ObjectMeta.Annotations["presslabs.com/display-name"] = displayName
	}
}

// New wraps a dashboardv1alpha1.Project into a Project object
func New(p *dashboardv1alpha1.Project) *Project {
	return &Project{p}
}

// Unwrap returns the wrapped dashboardv1alpha1.Project object
func (p *Project) Unwrap() *dashboardv1alpha1.Project {
	return p.Project
}

// Labels returns default label set for Project
func (p *Project) Labels() labels.Set {
	l := labels.Set{
		"presslabs.com/project": p.GetLabels()["presslabs.com/project"],
	}

	if p.ObjectMeta.Labels != nil {
		if org, ok := p.ObjectMeta.Labels["presslabs.com/organization"]; ok {
			l["presslabs.com/organization"] = org
		}
	}

	return l
}

// ComponentLabels returns labels for a label set for a Project component
func (p *Project) ComponentLabels(component component) labels.Set {
	labels := p.Labels()
	if len(component.app) > 0 {
		labels["app.kubernetes.io/name"] = component.app
	}
	if len(component.name) > 0 {
		labels["app.kubernetes.io/component"] = component.name
	}
	return labels
}

// ComponentName returns the object name for a component
func (p *Project) ComponentName(component component) string {
	if len(component.objNameFmt) == 0 {
		return component.objName
	}
	return fmt.Sprintf(component.objNameFmt, p.GetLabels()["presslabs.com/project"])
}

// Domain returns the project's subdomain label
func (p *Project) Domain() string {
	return p.Name
}

// ValidateMetadata validates the metadata of a Project
func (p *Project) ValidateMetadata() error {
	errorList := []error{}
	// Check for some required Project Labels and Annotations
	for _, label := range RequiredLabels {
		if value, exists := p.Project.Labels[label]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", label))
		}
	}

	// This case should not be reachable in normal circumstances
	if p.Project.Labels["presslabs.com/kind"] != "project" {
		errorList = append(errorList, errors.New("label \"presslabs.com/kind\" should be \"project\""))
	}

	for _, annotation := range RequiredAnnotations {
		if value, exists := p.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}
