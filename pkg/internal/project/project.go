/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package project

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// Project embeds dashboardv1alpha1.Project and adds utility functions
type Project struct {
	*corev1.Namespace
}

var (
	// RequiredLabels is a list of required Project labels
	RequiredLabels = []string{"presslabs.com/organization", "presslabs.com/project", "presslabs.com/kind"}
	// RequiredAnnotations is a list of required Project annotations
	RequiredAnnotations = []string{"presslabs.com/created-by"}
)

type component struct {
	name            string // eg. web, database, cache
	app             string // eg. mysql, memcached
	objNameFmt      string
	objName         string
	objNamespaceFmt string
	objNamespace    string
	kind            runtime.Object
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
	PrometheusRoleBinding = component{app: "prometheus", objName: "prometheus"}
	// Prometheus component
	Prometheus = component{app: "prometheus", objName: "prometheus"}
	// GiteaDeployment component
	GiteaDeployment = component{name: "web", app: "gitea", objName: "gitea"}
	// GiteaService component
	GiteaService = component{name: "web", app: "gitea", objName: "gitea"}
	// GiteaIngress component
	GiteaIngress = component{name: "web", app: "gitea", objName: "gitea"}
	// GiteaPVC component
	GiteaPVC = component{name: "web", app: "gitea", objName: "gitea"}
	// GiteaSecret component
	GiteaSecret = component{name: "web", app: "gitea", objName: "gitea-conf"}
	// OwnerRoleBinding component
	OwnerRoleBinding = component{
		kind:            &rbacv1.RoleBinding{},
		objName:         "owner",
		objNamespaceFmt: "%s",
	}
	// MemberRoleBinding component
	MemberRoleBinding = component{
		kind:            &rbacv1.RoleBinding{},
		objName:         "member",
		objNamespaceFmt: "%s",
	}
	// WordpressServiceMonitor component
	WordpressServiceMonitor = component{app: "prometheus", objName: "wordpress"}
	// MysqlServiceMonitor component
	MysqlServiceMonitor = component{app: "prometheus", objName: "mysql"}
	// MemcachedServiceMonitor component
	MemcachedServiceMonitor = component{app: "prometheus", objName: "memcached"}
)

// NamespaceName returns the name of the project's namespace
func NamespaceName(name string) string {
	return fmt.Sprintf("proj-%s", name)
}

// UpdateDisplayName updates the display-name annotation
func (p *Project) UpdateDisplayName(displayName string) {
	if len(displayName) == 0 {
		p.Namespace.ObjectMeta.Annotations["presslabs.com/display-name"] = p.Namespace.ObjectMeta.Labels["presslabs.com/project"]
	} else {
		p.Namespace.ObjectMeta.Annotations["presslabs.com/display-name"] = displayName
	}
}

// New wraps a dashboardv1alpha1.Project into a Project object
func New(obj *corev1.Namespace) *Project {
	return &Project{obj}
}

// Unwrap returns the wrapped dashboardv1alpha1.Project object
func (p *Project) Unwrap() *corev1.Namespace {
	return p.Namespace
}

// Labels returns default label set for dashboardv1alpha1.Project
func (p *Project) Labels() labels.Set {
	labels := labels.Set{
		"presslabs.com/project": p.GetLabels()["presslabs.com/project"],
	}

	if p.ObjectMeta.Labels != nil {
		if org, ok := p.ObjectMeta.Labels["presslabs.com/organization"]; ok {
			labels["presslabs.com/organization"] = org
		}
	}

	return labels
}

// ComponentLabels returns labels for a label set for a dashboardv1alpha1.Project component
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

// ComponentNamespace returns the object namespace for a component
func (p *Project) ComponentNamespace(component component) string {
	if len(component.objNamespaceFmt) == 0 {
		return component.objNamespace
	}
	return fmt.Sprintf(component.objNamespaceFmt, p.ObjectMeta.Name)
}

// Domain returns the project's subdomain label
func (p *Project) Domain() string {
	return p.Name
}

// ComponentObject returns a default object for the component
func (p *Project) ComponentObject(component component) runtime.Object {
	obj := component.kind.DeepCopyObject().(metav1.Object)
	obj.SetName(p.ComponentName(component))
	obj.SetNamespace(p.ComponentNamespace(component))
	return obj.(runtime.Object)
}

// ValidateMetadata validates the metadata of a Project
func (p *Project) ValidateMetadata() error {
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
		if value, exists := p.Namespace.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}
