/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package projects

import (
	"github.com/golang/glog"

	projectsv1 "github.com/presslabs/dashboard/pkg/apis/projects/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (c *Controller) syncNamespaces(project *projectsv1.Project) error {
	glog.Infof("Syncing namespace for project %s", project.ObjectMeta.Name)

	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.ObjectMeta.Name,
			Labels: labels.Set{
				"dashboard.presslabs.com/project": project.ObjectMeta.Name,
			},
		},
	}
	_, err := c.KubeClient.CoreV1().Namespaces().Create(&ns)

	if err != nil {
		glog.Errorf("Error while creating namespace for %s: %v", project.ObjectMeta.Name, err)
	} else {
		glog.Infof("Created namespace for %s", project.ObjectMeta.Name)
	}
	return err
}
