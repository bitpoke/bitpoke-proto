/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package util

import (
	"context"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	gomegatypes "github.com/onsi/gomega/types"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// BeInPhase is a helper func that returns a mather to check for namespace phase
func BeInPhase(phase corev1.NamespacePhase) gomegatypes.GomegaMatcher {
	return MatchFields(IgnoreExtras, Fields{
		"Status": MatchFields(IgnoreExtras, Fields{
			"Phase": Equal(phase),
		}),
	})
}

// HaveAnnotation is a helper func that returns a matcher to check for namespace
// annotations
func HaveAnnotation(annKey, annValue string) gomegatypes.GomegaMatcher {
	return MatchFields(IgnoreExtras, Fields{
		"ObjectMeta": MatchFields(IgnoreExtras, Fields{
			"Annotations": HaveKeyWithValue(annKey, annValue),
		}),
	})
}

// GetNamespace is a helper func that returns an organization
func GetNamespace(ctx context.Context, c client.Client, key client.ObjectKey) func() corev1.Namespace {
	return func() corev1.Namespace {
		var orgNs corev1.Namespace
		Expect(c.Get(ctx, key, &orgNs)).To(Succeed())
		return orgNs
	}
}
