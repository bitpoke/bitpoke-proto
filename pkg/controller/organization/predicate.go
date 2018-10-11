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

package organization

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// Predicate allows filtering only the organization namespaces
type Predicate struct{}

// IsObjectOrganization checks that the given meta Object is an Organization
func IsObjectOrganization(resource metav1.Object) bool {
	return resource.GetLabels()["presslabs.com/kind"] == "organization"
}

// Create returns true if the Create event should be processed
func (p *Predicate) Create(e event.CreateEvent) bool {
	return IsObjectOrganization(e.Meta)
}

// Delete returns true if the Delete event should be processed
func (p *Predicate) Delete(e event.DeleteEvent) bool {
	return IsObjectOrganization(e.Meta)
}

// Update returns true if the Update event should be processed
func (p *Predicate) Update(e event.UpdateEvent) bool {
	return IsObjectOrganization(e.MetaNew)
}

// Generic returns true if the Generic event should be processed
func (p *Predicate) Generic(e event.GenericEvent) bool {
	return IsObjectOrganization(e.Meta)
}
