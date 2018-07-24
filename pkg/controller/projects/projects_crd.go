/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package projects

import (
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/cache"

	projects "github.com/presslabs/dashboard/pkg/apis/projects/v1alpha1"
	projectslister "github.com/presslabs/dashboard/pkg/client/listers/projects/v1alpha1"
)

type ProjectsContext struct {
	// Projects CRD
	projectsQueue    *queue.Worker
	projectsInformer cache.SharedIndexInformer
	projectsLister   projectslister.ProjectLister
}

func (c *Controller) initProjectsWorker() {
	c.ProjectsContext = &ProjectsContext{
		projectsInformer: c.DashboardSharedInformerFactory.Dashboard().V1alpha1().Projects().Informer(),
		projectsLister:   c.DashboardSharedInformerFactory.Dashboard().V1alpha1().Projects().Lister(),
		projectsQueue:    queue.New("project", maxRetries, threadiness, c.reconcileProjects),
	}

	c.projectsInformer.AddEventHandler(queue.NewEventHandler(c.projectsQueue.GetQueue(), func(old interface{}, new interface{}) bool {
		return true
	}))
}

func (c *Controller) reconcileProjects(key string) error {
	obj, exists, err := c.projectsInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if exists {
		glog.Infof("Sync/Add/Update for Projects %s", key)

		project := obj.(*projects.Project).DeepCopy()
		namespace, err := c.syncNamespace(project)
		if err != nil {
			return err
		}

		project = obj.(*projects.Project).DeepCopy()
		_, err = c.syncPrometheus(project, namespace)
		return err
	}
	return nil
}
