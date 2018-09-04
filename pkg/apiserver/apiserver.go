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
	"net"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	project "github.com/presslabs/dashboard/pkg/apiserver/projects/v1"
	"github.com/presslabs/dashboard/pkg/cmd/apiserver/options"
	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type grpcRunner struct {
	client client.Client
}

var log = logf.Log.WithName("apiserver")

func (s *grpcRunner) Start(stop <-chan struct{}) error {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(middleware.Auth)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(middleware.Auth)),
	)
	project.RegisterProjectsServer(grpcServer, project.NewProjectServer(s.client))

	wrappedServer := grpcweb.WrapServer(grpcServer)

	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    options.HTTPAddr,
		Handler: http.HandlerFunc(handler),
	}

	errChan := make(chan error)

	lis, err := net.Listen("tcp", options.GRPCAddr)
	if err != nil {
		return err
	}

	go func() {
		log.Info("gRPC Server listening", "address", options.GRPCAddr)
		err := grpcServer.Serve(lis)
		errChan <- err
	}()

	go func() {
		log.Info("gRPC Web Server listening", "address", options.HTTPAddr)
		err := httpServer.ListenAndServe()
		errChan <- err
	}()

	go func() {
		<-stop
		httpServer.Shutdown(context.TODO())
		lis.Close()
	}()

	return <-errChan
}

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	m.Add(&grpcRunner{client: m.GetClient()})
	return nil
}
