/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package gomega

import (
	// nolint: golint,stylecheck
	. "github.com/onsi/gomega"
	// nolint: golint,stylecheck
	. "github.com/onsi/gomega/gstruct"

	gomegatypes "github.com/onsi/gomega/types"
	corev1 "k8s.io/api/core/v1"
	// logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// BeInPhase is a helper func that returns a matcher to check for namespace phase
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

// HaveLabel is a helper func that returns a matcher to check for namespace
// label
func HaveLabel(labelKey, labelValue string) gomegatypes.GomegaMatcher {
	return MatchFields(IgnoreExtras, Fields{
		"ObjectMeta": MatchFields(IgnoreExtras, Fields{
			"Labels": HaveKeyWithValue(labelKey, labelValue),
		}),
	})
}
