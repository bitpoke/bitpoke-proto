/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

// Package v1alpha1 contains API Schema definitions for the dashboard v1alpha1
// API group
//
//go:generate go run ../../../../vendor/k8s.io/code-generator/cmd/defaulter-gen/main.go -O zz_generated.defaults -i ./... -h ../../../../hack/boilerplate.go.txt
//
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/presslabs/dashboard/pkg/apis/dashboard
// +k8s:defaulter-gen=TypeMeta
// +groupName=dashboard.presslabs.com
//
package v1alpha1
