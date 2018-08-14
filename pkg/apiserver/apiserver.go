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

package apiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	project "github.com/presslabs/dashboard/pkg/apiserver/projects/v1"
	"google.golang.org/grpc"
)

type grpcRunner struct {
	client client.Client
}

func (s *grpcRunner) Start(stop <-chan struct{}) error {
	var (
		httpServer http.Server
		opts       []grpc.ServerOption
		port       = 9090
	)

	grpcServer := grpc.NewServer(opts...)
	project.RegisterProjectsServer(grpcServer, project.NewProjectServer(s.client))

	wrappedServer := grpcweb.WrapServer(grpcServer)

	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}

	resources := grpcweb.ListGRPCResources(grpcServer)

	httpServer = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
		// Handler: c.Handler(http.HandlerFunc(handler)),
	}

	go func() {
		<-stop
		httpServer.Shutdown(context.TODO())
	}()

	log.Printf("Server started on http://0.0.0.0:%d", port)
	log.Printf("Available resources: %v", resources)

	return httpServer.ListenAndServe()
}

func AddToManager(m manager.Manager) error {
	m.Add(&grpcRunner{client: m.GetClient()})
	return nil
}
