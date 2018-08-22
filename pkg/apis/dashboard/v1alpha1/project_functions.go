package v1alpha1

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

// GetNamespaceName returns the name of the project's namespace
func (p *Project) GetNamespaceName() string {
	return p.Name
}

// GetNamespaceKey returns the project's key through which the project may be identified
func (p *Project) GetNamespaceKey() types.NamespacedName {
	return types.NamespacedName{
		Name: p.GetNamespaceName(),
	}
}

// GetProjectLabel returns a label that should be applied on objects belonging to a
// project
func (p *Project) GetProjectLabel() labels.Set {
	return labels.Set{
		"project.dashboard.presslabs.com/project": p.Name,
	}
}

// GetDeployManagerLabel returns a label that should be applied on objects managed
// by the project controller
func (p *Project) GetDeployManagerLabel() labels.Set {
	return labels.Set{
		"app.kubernetes.io/deploy-manager": "project-controller.dashboard.presslabs.com",
	}
}

// GetDefaultLabels returns a set of labels that should be applied on objects
// managed by the project controller
func (p *Project) GetDefaultLabels() labels.Set {
	return labels.Merge(p.GetProjectLabel(), p.GetDeployManagerLabel())
}

// GetPrometheusName returns the name of the Prometheus resource
func (p *Project) GetPrometheusName() string {
	return "prometheus"
}

// GetPrometheusKey returns the project's key through which the Prometheus may be identified
func (p *Project) GetPrometheusKey() types.NamespacedName {
	return types.NamespacedName{
		Namespace: p.GetNamespaceName(),
		Name:      p.GetPrometheusName(),
	}
}

// GetResourceQuotaName returns the name of the Prometheus resource
func (p *Project) GetResourceQuotaName() string {
	return p.Name
}

// GetResourceQuotaKey returns the project's key through which the Prometheus may be identified
func (p *Project) GetResourceQuotaKey() types.NamespacedName {
	return types.NamespacedName{
		Namespace: p.GetNamespaceName(),
		Name:      p.GetResourceQuotaName(),
	}
}
