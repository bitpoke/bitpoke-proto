/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package organization

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	// logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/presslabs/dashboard/pkg/apis"
)

var cfg *rest.Config
var t *envtest.Environment

func TestAPIServer(t *testing.T) {
	// logf.SetLogger(logf.ZapLogger(true))
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "API Server", []Reporter{envtest.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	var err error
	t = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crds"),
			filepath.Join("..", "..", "vendor/github.com/coreos/prometheus-operator/example/prometheus-operator-crd"),
			filepath.Join("..", "..", "vendor/github.com/presslabs/wordpress-operator/config/crds"),
		},
	}
	apis.AddToScheme(scheme.Scheme)

	cfg, err = t.Start()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	t.Stop()
})

func StartTestManager(mgr manager.Manager) chan struct{} {
	stop := make(chan struct{})
	go func() {
		defer GinkgoRecover()
		Expect(mgr.Start(stop)).To(Succeed())
	}()
	return stop
}
