/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
)

// giteaLabels returns a set of labels that can be used to identify Gitea related resources
func giteaLabels(project *dashboardv1alpha1.Project) labels.Set {
	giteaSelector := labels.Set{
		"app.kubernetes.io/name": giteaName,
	}
	return labels.Merge(getDefaultLabels(project), giteaSelector)
}

// giteaPodLabels returns a set of labels that should be applied on Gitea related objects that are managed by the project controller
func giteaPodLabels(project *dashboardv1alpha1.Project) labels.Set {
	giteaPodLabels := labels.Set{
		"app.kubernetes.io/version": giteaReleaseVersion,
	}
	return labels.Merge(giteaLabels(project), giteaPodLabels)
}

func giteaDomain(project *dashboardv1alpha1.Project) string {
	return fmt.Sprintf("%s.%s", getProjectDomainName(project), options.GitBaseDomainURL)
}

// giteaSecretName returns the name of the Gitea Secret
func giteaSecretName(project *dashboardv1alpha1.Project) string {
	return "gitea-conf"
}

// giteaPVCName returns the name of the Gitea PVC
func giteaPVCName(project *dashboardv1alpha1.Project) string {
	return "gitea" // nolint
}

// giteaDeploymentName returns the name of the Gitea Deployment
func giteaDeploymentName(project *dashboardv1alpha1.Project) string {
	return "gitea" // nolint
}

// giteaServiceName returns the name of the Gitea Service
func giteaServiceName(project *dashboardv1alpha1.Project) string {
	return "gitea" // nolint
}

// giteaIngressName returns the name of the Gitea Service
func giteaIngressName(project *dashboardv1alpha1.Project) string {
	return "gitea" // nolint
}
