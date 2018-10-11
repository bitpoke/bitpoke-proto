package project

import (
	"fmt"


	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// GetNamespaceName returns the name of the Project's Namespace
func GetNamespaceName(project *dashboardv1alpha1.Project) string {
	return fmt.Sprintf("proj-%s", project.Name)
}
