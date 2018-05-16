/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package projects

import (
	"github.com/golang/glog"

	projectsApi "github.com/presslabs/dashboard/pkg/apis/projects/v1alpha1"

	coreapi "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) syncNamespaces(project *projectsApi.Project) error {
	glog.Infof("Syncing namespace for project %s", project.ObjectMeta.Name)

	nsSpec := coreapi.Namespace{ObjectMeta: metav1.ObjectMeta{Name: project.ObjectMeta.Name}}
	_, err := c.KubeClient.CoreV1().Namespaces().Create(&nsSpec)

	glog.Infof("Created namespace for %s", project.ObjectMeta.Name)

	return err
}
