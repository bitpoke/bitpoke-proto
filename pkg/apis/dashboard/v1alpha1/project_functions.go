package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
)

// GetNamespaceName returns the name of the project's namespace
func (p *Project) GetNamespaceName() string {
	return fmt.Sprintf("proj-%s-%s", p.Namespace, p.Name)
}

// GetNamespaceKey returns the key through which the project may be identified
func (p *Project) GetNamespaceKey() types.NamespacedName {
	return types.NamespacedName{
		Name: p.GetNamespaceName(),
	}
}

// GetOrganizationName returns the name of the organization which the project belongs too
func (p *Project) GetOrganizationName() string {
	return p.Namespace
}

// GetProjectDomainName returns the DNS domain name for a project
func (p *Project) GetProjectDomainName() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Namespace)
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

// GetPrometheusKey returns the key through which the Prometheus may be identified
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

// GetResourceQuotaKey returns the key through which the Prometheus may be identified
func (p *Project) GetResourceQuotaKey() types.NamespacedName {
	return types.NamespacedName{
		Namespace: p.GetNamespaceName(),
		Name:      p.GetResourceQuotaName(),
	}
}

// GetGiteaSecretName returns the name of the Gitea Secret
func (p *Project) GetGiteaSecretName() string {
	return "gitea-conf"
}

// GetGiteaSecretKey returns the key through which the Gitea Secret may be identified
func (p *Project) GetGiteaSecretKey() types.NamespacedName {
	return types.NamespacedName{Name: p.GetGiteaSecretName(), Namespace: p.GetNamespaceName()}
}

// GetGiteaPVCName returns the name of the Gitea PVC
func (p *Project) GetGiteaPVCName() string {
	return "gitea" // nolint
}

// GetGiteaPVCKey returns the key through which the Gitea Secret may be identified
func (p *Project) GetGiteaPVCKey() types.NamespacedName {
	return types.NamespacedName{Name: p.GetGiteaPVCName(), Namespace: p.GetNamespaceName()}
}

// GetGiteaDeploymentName returns the name of the Gitea Deployment
func (p *Project) GetGiteaDeploymentName() string {
	return "gitea" // nolint
}

// GetGiteaDeploymentKey returns the key through which the Gitea Deployment may be identified
func (p *Project) GetGiteaDeploymentKey() types.NamespacedName {
	return types.NamespacedName{Name: p.GetGiteaDeploymentName(), Namespace: p.GetNamespaceName()}
}

// GetGiteaServiceName returns the name of the Gitea Service
func (p *Project) GetGiteaServiceName() string {
	return "gitea" // nolint
}

// GetGiteaServiceKey returns the key through which the Gitea Service may be identified
func (p *Project) GetGiteaServiceKey() types.NamespacedName {
	return types.NamespacedName{Name: p.GetGiteaServiceName(), Namespace: p.GetNamespaceName()}
}

// GetGiteaIngressName returns the name of the Gitea Service
func (p *Project) GetGiteaIngressName() string {
	return "gitea" // nolint
}

// GetGiteaIngressKey returns the key through which the Gitea Service may be identified
func (p *Project) GetGiteaIngressKey() types.NamespacedName {
	return types.NamespacedName{Name: p.GetGiteaIngressName(), Namespace: p.GetNamespaceName()}
}

// GetGiteaDomain returns the project Git Repository URL Path
func (p *Project) GetGiteaDomain() string {
	return fmt.Sprintf("%s.%s", p.GetProjectDomainName(), options.GitBaseDomainURL)
}
