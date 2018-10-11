/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package site

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

// Site embeds wordpressv1alpha1.Wordpress and adds utility functions
type Site struct {
	*wordpressv1alpha1.Wordpress
}

type component struct {
	name       string // eg. web, database, cache
	app        string // eg. mysql, memcached
	objNameFmt string
	objName    string
}

var (
	// WordpressServiceMonitor component
	WordpressServiceMonitor = component{name: "web", app: "wordpress", objNameFmt: "%s-wordpress"}
	// MysqlCluster component
	MysqlCluster = component{name: "database", app: "mysql", objNameFmt: "%s"}
	// MysqlClusterSecret component
	MysqlClusterSecret = component{name: "database", app: "mysql", objNameFmt: "%s-mysql"}
	// MysqlServiceMonitor component
	MysqlServiceMonitor = component{name: "database", app: "mysql", objNameFmt: "%s-mysql"}
	// MemcachedService component
	MemcachedService = component{name: "cache", app: "memcached", objNameFmt: "%s-memcached"}
	// MemcachedStatefulSet component
	MemcachedStatefulSet = component{name: "cache", app: "memcached", objNameFmt: "%s-memcached"}
	// MemcachedServiceMonitor component
	MemcachedServiceMonitor = component{name: "cache", app: "memcached", objNameFmt: "%s-memcached"}
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
	return labels.Set{
		"app.kubernetes.io/part-of":  "wordpress",
		"app.kubernetes.io/instance": o.ObjectMeta.Name,
	}
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
