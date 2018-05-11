/*
Copyright 2018 Pressinfra SRL

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package projects

import (
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/cache"

	projectsApi "github.com/presslabs/dashboard/pkg/apis/projects/v1alpha1"
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
		project := obj.(*projectsApi.Project).DeepCopy()

		err = c.syncNamespaces(project)
		return err
	}
	return nil
}
