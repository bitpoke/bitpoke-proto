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

package podplacement_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/presslabs/dashboard/pkg/apis"
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
)

var cfg *rest.Config
var t *envtest.Environment

func TestProjectListing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Webhook Server", []Reporter{envtest.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	var err error

	options.WebhookSecretName = ""
	options.WebhookService = ""

	t = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crds"),
			//filepath.Join("..", "..", "vendor/github.com/coreos/prometheus-operator/example/prometheus-operator-crd"),
			//filepath.Join("..", "..", "vendor/github.com/presslabs/wordpress-operator/config/crds"),
		},
	}
	apis.AddToScheme(scheme.Scheme)

	cfg, err = t.Start()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	t.Stop()
})

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager) chan struct{} {
	stop := make(chan struct{})
	go func() {
		Expect(mgr.Start(stop)).To(Succeed())
	}()
	return stop
}
