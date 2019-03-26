/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package site

import (
	"fmt"
	"path"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/dashboard/pkg/internal/project"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
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
	// MysqlCluster component
	MysqlCluster = component{name: "database", app: "mysql", objNameFmt: "%s"}
	// MysqlClusterSecret component
	MysqlClusterSecret = component{name: "database", app: "mysql", objNameFmt: "%s-mysql"}
	// MemcachedService component
	MemcachedService = component{name: "cache", app: "memcached", objNameFmt: "%s-memcached"}
	// MemcachedStatefulSet component
	MemcachedStatefulSet = component{name: "cache", app: "memcached", objNameFmt: "%s-memcached"}
)

type component struct {
	name       string // eg. web, database, cache
	app        string // eg. mysql, memcached
	objNameFmt string
	objName    string
}

const (
	// Prefix for site fully-qualified project name
	prefix = "site/"
)

// New wraps a wordpressv1alpha1.Wordpress into a Site object
func New(obj *wordpressv1alpha1.Wordpress) *Site {
	return &Site{obj}
}

// Unwrap returns the wrapped wordpressv1alpha1.Wordpress object
func (s *Site) Unwrap() *wordpressv1alpha1.Wordpress {
	return s.Wordpress
}

// Labels returns default label set for wordpressv1alpha1.Wordpress
func (s *Site) Labels() labels.Set {
	l := labels.Set{
		"app.kubernetes.io/part-of":  "wordpress",
		"app.kubernetes.io/instance": s.ObjectMeta.Name,
	}

	if s.ObjectMeta.Labels != nil {
		if org, ok := s.ObjectMeta.Labels["presslabs.com/organization"]; ok {
			l["presslabs.com/organization"] = org
		}

		if proj, ok := s.ObjectMeta.Labels["presslabs.com/project"]; ok {
			l["presslabs.com/project"] = proj
		}
	}

	return l
}

// ComponentLabels returns labels for a label set for a wordpressv1alpha1.Wordpress component
func (s *Site) ComponentLabels(component component) labels.Set {
	l := s.Labels()
	if len(component.app) > 0 {
		l["app.kubernetes.io/name"] = component.app
	}
	if len(component.name) > 0 {
		l["app.kubernetes.io/component"] = component.name
	}
	return l
}

// ComponentName returns the object name for a component
func (s *Site) ComponentName(component component) string {
	if len(component.objNameFmt) == 0 {
		return component.objName
	}
	return fmt.Sprintf(component.objNameFmt, s.ObjectMeta.Name)
}

// ValidateMetadata validates the metadata of a Site
func (s *Site) ValidateMetadata() error {
	errorList := []error{}
	// Check for some required Project Labels and Annotations
	for _, l := range RequiredLabels {
		if value, exists := s.Wordpress.Labels[l]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", l))
		}
	}

	for _, annotation := range RequiredAnnotations {
		if value, exists := s.Wordpress.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}

// FQName returns the fully-qualified site name
func FQName(projName, siteName string) string {
	return path.Join(project.FQName(projName), prefix, siteName)
}

// Resolve resolves an fully-qualified site name to a k8s object name.
// The function returns site name, project name and error
func Resolve(name string) (string, string, error) {
	if path.Clean(name) != name {
		return "", "", fmt.Errorf("site resources fully-qualified name must be in form project/PROJECT-NAME/site/SITE-NAME")
	}

	matched, err := path.Match("project/*/site/*", name)
	if err != nil || !matched {
		return "", "", fmt.Errorf("site resources fully-qualified name must be in form project/PROJECT-NAME/site/SITE-NAME")
	}

	names := strings.Split(name, "/")
	return names[3], names[1], nil
}

// ResolveToObjectKey resolves an fully-qualified site name to a k8s object name.
// The function returns the object key from FQName and an error
func ResolveToObjectKey(c client.Client, fqSiteName, orgName string) (*client.ObjectKey, error) {
	siteName, projName, err := Resolve(fqSiteName)
	if err != nil {
		return nil, err
	}

	ns, err := projectns.Lookup(c, projName, orgName)
	if err != nil {
		return nil, err
	}

	key := client.ObjectKey{
		Name:      siteName,
		Namespace: ns.Name,
	}
	return &key, nil
}
