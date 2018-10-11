package project

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

// RequiredLabels defines the required labels for a Site resource
func RequiredLabels() []string {
	return []string{"presslabs.com/organization", "presslabs.com/project", "presslabs.com/site"}
}

// RequiredAnnotations defines the required annotations for a Site resource
func RequiredAnnotations() []string {
	return []string{"presslabs.com/created-by"}
}

// SetMetadata sets the required metadata for a Site resources
func SetMetadata(objMeta *metav1.ObjectMeta, site *wordpressv1alpha1.Wordpress, createdBy string) {
	objMeta.Labels["presslabs.com/site"] = site.Name
	objMeta.Labels["presslabs.com/project"] = site.Labels["presslabs.com/project"]
	objMeta.Labels["presslabs.com/organization"] = site.Labels["presslabs.com/organization"]
	objMeta.Annotations["presslabs.com/created-by"] = createdBy
}

// ValidateMetadata validates the metadata of a Site
func ValidateMetadata(site *wordpressv1alpha1.Wordpress) error {
	errorList := []error{}
	// Check for some required Project Labels and Annotations
	for _, label := range RequiredLabels() {
		if value, exists := site.Labels[label]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required label \"%s\" is missing", label))
		}
	}
	for _, annotation := range RequiredAnnotations() {
		if value, exists := site.Annotations[annotation]; !exists || value == "" {
			errorList = append(errorList, fmt.Errorf("required annotation \"%s\" is missing", annotation))
		}
	}

	return utilerrors.Flatten(utilerrors.NewAggregate(errorList))
}
