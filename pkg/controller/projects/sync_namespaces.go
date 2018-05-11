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
