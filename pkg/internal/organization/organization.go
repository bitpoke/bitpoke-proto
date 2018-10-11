/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package organization

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// Organization embeds dashboardv1alpha1.Organization and adds utility functions
type Organization struct {
	*corev1.Namespace
}

var (
	// RequiredLabels is a list of required Organization labels
	RequiredLabels = []string{"presslabs.com/organization", "presslabs.com/kind"}
	// RequiredAnnotations is a list of required Organization annotations
	RequiredAnnotations = []string{"presslabs.com/created-by"}
)

// Component is a component type of Organization
type Component struct {
	name            string // eg. web, database, cache
	app             string // eg. mysql, memcached
	objNameFmt      string
	objName         string
	objNamespaceFmt string
	objNamespace    string
	kind            runtime.Object
}

var (
	// OwnerClusterRole component
	OwnerClusterRole = Component{
		kind:       &rbacv1.ClusterRole{},
		objNameFmt: "dashboard.presslabs.com:organization:%s:owner",
	}
	// OwnerClusterRoleBinding component
	OwnerClusterRoleBinding = Component{
		kind:       &rbacv1.ClusterRoleBinding{},
		objNameFmt: "dashboard.presslabs.com:organization:%s:owners",
	}
	// MemberRoleBinding component
	MemberRoleBinding = Component{
		kind:            &rbacv1.RoleBinding{},
		objName:         "members",
		objNamespaceFmt: "%s",
	}
)

// New wraps a dashboardv1alpha1.Organization into a Organization object
func New(obj *corev1.Namespace) *Organization {
	return &Organization{obj}
}

// Unwrap returns the wrapped dashboardv1alpha1.Organization object
func (o *Organization) Unwrap() *corev1.Namespace {
	return o.Namespace
}

// Labels returns default label set for corev1.Namespace organization
func (o *Organization) Labels() labels.Set {
	labels := labels.Set{
		"presslabs.com/organization": o.ObjectMeta.Labels["presslabs.com/organization"],
	}

	return labels
}

// ComponentLabels returns labels for a label set for a dashboardv1alpha1.Organization component
func (o *Organization) ComponentLabels(component Component) labels.Set {
	labels := o.Labels()
	if len(component.app) > 0 {
		labels["app.kubernetes.io/name"] = component.app
	}
	if len(component.name) > 0 {
		labels["app.kubernetes.io/component"] = component.name
	}
	return labels
}

// ComponentName returns the object name for a component
func (o *Organization) ComponentName(component Component) string {
	if len(component.objNameFmt) == 0 {
		return component.objName
	}
	return fmt.Sprintf(component.objNameFmt, o.ObjectMeta.Labels["presslabs.com/organization"])
}

// ComponentNamespace returns the object name for a component
func (o *Organization) ComponentNamespace(component Component) string {
	if len(component.objNamespaceFmt) == 0 {
		return component.objNamespace
	}
	return fmt.Sprintf(component.objNamespaceFmt, o.ObjectMeta.Name)
}

// ComponentObject returns a default object for the component
func (o *Organization) ComponentObject(component Component) runtime.Object {
	obj := component.kind.DeepCopyObject().(metav1.Object)
	obj.SetName(o.ComponentName(component))
	obj.SetNamespace(o.ComponentNamespace(component))
	return obj.(runtime.Object)
}

// ValidateMetadata validates the metadata of a Organization
func (o *Organization) ValidateMetadata() error {
	errorList := []error{}
	// Check for some required Organization Labels and Annotations
	for _, label := range RequiredLabels {
		if value, exists := o.Namespace.Labels[label]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", label))
		}
	}
	for _, annotation := range RequiredAnnotations {
		if value, exists := o.Namespace.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}
