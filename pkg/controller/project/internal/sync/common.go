/*
Copyright 2018 Pressinfra SRL.

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

package sync

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

var log = logf.Log.WithName("project-controller")

// getNamespaceName returns the name of the project's namespace
func getNamespaceName(project *dashboardv1alpha1.Project) string {
	return fmt.Sprintf("proj-%s-%s", project.Namespace, project.Name)
}

// getProjectDomainName returns the DNS domain name for a project
func getProjectDomainName(project *dashboardv1alpha1.Project) string {
	return fmt.Sprintf("%s-%s", project.Name, project.Namespace)
}

// getProjectLabel returns a label that should be applied on objects belonging to a
// project
func getProjectLabel(project *dashboardv1alpha1.Project) labels.Set {
	return labels.Set{
		"project.dashboard.presslabs.com/project": project.Name,
	}
}

// getDeployManagerLabel returns a label that should be applied on objects managed
// by the project controller
func getDeployManagerLabel(project *dashboardv1alpha1.Project) labels.Set { // nolint: unparam
	return labels.Set{
		"app.kubernetes.io/deploy-manager": "project-controller.dashboard.presslabs.com",
	}
}

// getDefaultLabels returns a set of labels that should be applied on objects
// managed by the project controller
func getDefaultLabels(project *dashboardv1alpha1.Project) labels.Set {
	return labels.Merge(getProjectLabel(project), getDeployManagerLabel(project))
}
