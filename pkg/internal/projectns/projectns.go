/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package projectns

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProjectNamespace embeds corev1.Namespace and adds utility functions
type ProjectNamespace struct {
	*corev1.Namespace
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

var (
	// Namespace component
	Namespace = component{objNameFmt: "proj-%s"}
	// LimitRange component
	LimitRange = component{objName: "presslabs-dashboard"}
	// ResourceQuota component
	ResourceQuota = component{objName: "presslabs-dashboard"}
	// PrometheusServiceAccount component
	PrometheusServiceAccount = component{app: "prometheus", objName: "prometheus"}
	// PrometheusRoleBinding for ServiceAccount component
	PrometheusRoleBinding = component{app: "prometheus", objName: "dashboard.presslabs.com:project:prometheus"}
	// Prometheus component
	Prometheus = component{app: "prometheus", objName: "prometheus"}
	// OwnerRoleBinding component
	OwnerRoleBinding = component{objName: "dashboard.presslabs.com:project:owners"}
	// MemberRoleBinding component
	MemberRoleBinding = component{objName: "dashboard.presslabs.com:project:members"}
	// WordpressServiceMonitor component
	WordpressServiceMonitor = component{app: "prometheus", objName: "wordpress"}
	// MysqlServiceMonitor component
	MysqlServiceMonitor = component{app: "prometheus", objName: "mysql"}
	// MemcachedServiceMonitor component
	MemcachedServiceMonitor = component{app: "prometheus", objName: "memcached"}
	// SMTPSecret component
	SMTPSecret = component{objName: "default-smtp-credentials"}
)

// UpdateDisplayName updates the display-name annotation
func (p *ProjectNamespace) UpdateDisplayName(displayName string) {
	if len(displayName) == 0 {
		p.ObjectMeta.Annotations["presslabs.com/display-name"] = p.ObjectMeta.Labels["presslabs.com/project"]
	} else {
		p.ObjectMeta.Annotations["presslabs.com/display-name"] = displayName
	}
}

// New wraps a Namespace into a ProjectNamespace object
func New(p *corev1.Namespace) *ProjectNamespace {
	return &ProjectNamespace{p}
}

// Unwrap returns the wrapped Namespace object
func (p *ProjectNamespace) Unwrap() *corev1.Namespace {
	return p.Namespace
}

// Labels returns default label set for ProjectNamespace
func (p *ProjectNamespace) Labels() labels.Set {
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

// ComponentLabels returns labels for a label set for a ProjectNamespace component
func (p *ProjectNamespace) ComponentLabels(component component) labels.Set {
	l := p.Labels()
	if len(component.app) > 0 {
		l["app.kubernetes.io/name"] = component.app
	}
	if len(component.name) > 0 {
		l["app.kubernetes.io/component"] = component.name
	}
	return l
}

// ComponentName returns the object name for a component
func (p *ProjectNamespace) ComponentName(component component) string {
	if len(component.objNameFmt) == 0 {
		return component.objName
	}
	return fmt.Sprintf(component.objNameFmt, p.GetLabels()["presslabs.com/project"])
}

// Domain returns the project's subdomain label
func (p *ProjectNamespace) Domain() string {
	return p.Name
}

// ValidateMetadata validates the metadata of a ProjectNamespace
func (p *ProjectNamespace) ValidateMetadata() error {
	errorList := []error{}
	// Check for some required Project Labels and Annotations
	for _, label := range RequiredLabels {
		if value, exists := p.Namespace.Labels[label]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", label))
		}
	}

	// This case should not be reachable in normal circumstances
	if p.Namespace.Labels["presslabs.com/kind"] != "project" {
		errorList = append(errorList, errors.New("label \"presslabs.com/kind\" should be \"project\""))
	}

	for _, annotation := range RequiredAnnotations {
		if value, exists := p.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}

// Lookup gets the namespace for a project name and organization name
func Lookup(c client.Client, projName, orgName string) (*ProjectNamespace, error) {
	nsList := &corev1.NamespaceList{}
	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set{
				"presslabs.com/kind":         "project",
				"presslabs.com/project":      projName,
				"presslabs.com/organization": orgName,
			},
		),
	}

	if err := c.List(context.TODO(), listOptions, nsList); err != nil {
		return nil, err
	}

	return New(&nsList.Items[0]), nil
}
