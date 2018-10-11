package project

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// RequiredLabels defines the required labels for a Project resource
func RequiredLabels() []string {
	return []string{"presslabs.com/organization", "presslabs.com/project"}
}

// RequiredAnnotations defines the required annotations for a Project resource
func RequiredAnnotations() []string {
	return []string{"presslabs.com/created-by"}
}

// GetNamespaceName returns the name of the Project's Namespace
func GetNamespaceName(project *dashboardv1alpha1.Project) string {
	return fmt.Sprintf("proj-%s", project.Name)
}

// SetMetadata sets the required metadata for a Project resources
func SetMetadata(objMeta *metav1.ObjectMeta, project *dashboardv1alpha1.Project, org, createdBy string) {
	objMeta.Labels["presslabs.com/organization"] = org
	objMeta.Labels["presslabs.com/project"] = project.Name
	objMeta.Annotations["presslabs.com/created-by"] = createdBy
}

// ValidateMetadata validates the metadata of a Project
func ValidateMetadata(project *dashboardv1alpha1.Project) error {
	errorList := []error{}
	// Check for some required Project Labels and Annotations
	for _, label := range RequiredLabels() {
		if value, exists := project.Labels[label]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", label))
		}
	}
	for _, annotation := range RequiredAnnotations() {
		if value, exists := project.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}
