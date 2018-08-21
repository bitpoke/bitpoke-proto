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

package v1alpha1

import (
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("dashboard")

var cfg *rest.Config
var c client.Client

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "..", "config", "crds")},
	}

	err := SchemeBuilder.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Error(err, "unable to add scheme to scheme builder")
		os.Exit(1)
	}

	if cfg, err = t.Start(); err != nil {
		log.Error(err, "unable to start test environment")
		os.Exit(1)
	}

	if c, err = client.New(cfg, client.Options{Scheme: scheme.Scheme}); err != nil {
		log.Error(err, "unable to create a new client")
		os.Exit(1)
	}

	code := m.Run()
	t.Stop()
	os.Exit(code)
}
