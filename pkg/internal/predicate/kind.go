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

package predicate

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// KindPredicate allows filtering only the organization namespaces
type KindPredicate struct {
	kind string
}

var _ predicate.Predicate = &KindPredicate{}

// NewKindPredicate return a new KindPredicate
func NewKindPredicate(kind string) *KindPredicate {
	return &KindPredicate{kind}
}

func (p *KindPredicate) isOfKind(resource metav1.Object) bool {
	return resource.GetLabels()["presslabs.com/kind"] == p.kind
}

// Create returns true if the Create event should be processed
func (p *KindPredicate) Create(e event.CreateEvent) bool {
	return p.isOfKind(e.Meta)
}

// Delete returns true if the Delete event should be processed
func (p *KindPredicate) Delete(e event.DeleteEvent) bool {
	return p.isOfKind(e.Meta)
}

// Update returns true if the Update event should be processed
func (p *KindPredicate) Update(e event.UpdateEvent) bool {
	return p.isOfKind(e.MetaNew)
}

// Generic returns true if the Generic event should be processed
func (p *KindPredicate) Generic(e event.GenericEvent) bool {
	return p.isOfKind(e.Meta)
}
