/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package apis

import dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, dashboardv1alpha1.SchemeBuilder.AddToScheme)
}
