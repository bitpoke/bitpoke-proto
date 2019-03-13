/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/presslabs/dashboard/pkg/apis"
)

func TestProjectGiteaIngress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Project Sync Suite", []Reporter{envtest.NewlineReporter{}})
}

var rts = scheme.Scheme

var _ = BeforeSuite(func() {
	apis.AddToScheme(rts)
})
