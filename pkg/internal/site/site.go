/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package site

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

// Site embeds wordpressv1alpha1.Wordpress and adds utility functions
type Site struct {
	*wordpressv1alpha1.Wordpress
}

var (
	// RequiredLabels is a list of required Site labels
	RequiredLabels = []string{"presslabs.com/organization", "presslabs.com/project"}
	// RequiredAnnotations is a list of required Site annotations
	RequiredAnnotations = []string{"presslabs.com/created-by"}
)

type component struct {
	name       string // eg. web, database, cache
	app        string // eg. mysql, memcached
	objNameFmt string
	objName    string
}

var (
	// MysqlCluster component
	MysqlCluster = component{name: "database", app: "mysql", objNameFmt: "%s"}
	// MysqlClusterSecret component
	MysqlClusterSecret = component{name: "database", app: "mysql", objNameFmt: "%s-mysql"}
	// MemcachedService component
	MemcachedService = component{name: "cache", app: "memcached", objNameFmt: "%s-memcached"}
	// MemcachedStatefulSet component
	MemcachedStatefulSet = component{name: "cache", app: "memcached", objNameFmt: "%s-memcached"}
)

// New wraps a wordpressv1alpha1.Wordpress into a Site object
func New(obj *wordpressv1alpha1.Wordpress) *Site {
	return &Site{obj}
}

// Unwrap returns the wrapped wordpressv1alpha1.Wordpress object
func (o *Site) Unwrap() *wordpressv1alpha1.Wordpress {
	return o.Wordpress
}

// Labels returns default label set for wordpressv1alpha1.Wordpress
func (o *Site) Labels() labels.Set {
	labels := labels.Set{
		"app.kubernetes.io/part-of":  "wordpress",
		"app.kubernetes.io/instance": o.ObjectMeta.Name,
	}

	if o.ObjectMeta.Labels != nil {
		if org, ok := o.ObjectMeta.Labels["presslabs.com/organization"]; ok {
			labels["presslabs.com/organization"] = org
		}

		if proj, ok := o.ObjectMeta.Labels["presslabs.com/project"]; ok {
			labels["presslabs.com/project"] = proj
		}
	}

	return labels
}

// ComponentLabels returns labels for a label set for a wordpressv1alpha1.Wordpress component
func (o *Site) ComponentLabels(component component) labels.Set {
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
func (o *Site) ComponentName(component component) string {
	if len(component.objNameFmt) == 0 {
		return component.objName
	}
	return fmt.Sprintf(component.objNameFmt, o.ObjectMeta.Name)
}

// ValidateMetadata validates the metadata of a Site
func (o *Site) ValidateMetadata() error {
	errorList := []error{}
	// Check for some required Project Labels and Annotations
	for _, label := range RequiredLabels {
		if value, exists := o.Wordpress.Labels[label]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", label))
		}
	}

	for _, annotation := range RequiredAnnotations {
		if value, exists := o.Wordpress.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}