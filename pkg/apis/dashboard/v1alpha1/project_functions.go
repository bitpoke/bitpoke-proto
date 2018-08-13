package v1alpha1

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

func (p *Project) GetNamespaceName() string {
	return p.Name
}

func (p *Project) GetNamespaceKey() types.NamespacedName {
	return types.NamespacedName{
		Name: p.GetNamespaceName(),
	}
}

func (p *Project) GetProjectLabel() labels.Set {
	return labels.Set{
		"project.dashboard.presslabs.com/project": p.Name,
	}
}

func (p *Project) GetDeployManagerLabel() labels.Set {
	return labels.Set{
		"app.kubernetes.io/deploy-manager": "project-controller.dashboard.presslabs.com",
	}
}

func (p *Project) GetDefaultLabels() labels.Set {
	return labels.Merge(p.GetProjectLabel(), p.GetDeployManagerLabel())
}

func (p *Project) GetPrometheusLabels() labels.Set {
	labels := p.GetDefaultLabels()
	labels["app.kubernetes.io/name"] = "prometheus"

	return labels
}

func (p *Project) GetPrometheusName() string {
	return "prometheus"
}

func (p *Project) GetPrometheusKey() types.NamespacedName {
	return types.NamespacedName{
		Namespace: p.GetNamespaceName(),
		Name:      p.GetPrometheusName(),
	}
}
