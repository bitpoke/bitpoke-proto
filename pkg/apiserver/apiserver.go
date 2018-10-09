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

	"github.com/gobuffalo/packr"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	projectv1 "github.com/presslabs/dashboard/pkg/api/core/v1"
	"github.com/presslabs/dashboard/pkg/apiserver/middleware"
	"github.com/presslabs/dashboard/pkg/cmd/apiserver/options"
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
	projectv1.RegisterProjectsServer(grpcServer, projectv1.NewProjectServer(s.client))

	box := packr.NewBox("../../app/build")
	if !box.Has("index.html") {
		panic("Cannot find 'index.html' web server entry point. You need to build the webapp first.")
	}

	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	handler := func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		}
		// Fall back to other servers.
		http.DefaultServeMux.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    options.HTTPAddr,
		Handler: http.HandlerFunc(handler),
	}

	errChan := make(chan error, 2)

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
		log.Info("Web Server listening", "address", options.HTTPAddr)
		http.Handle("/", http.FileServer(box))
		err := httpServer.ListenAndServe()
		errChan <- err
	}()

	go func() {
		<-stop
		err := httpServer.Shutdown(context.TODO())
		if err != nil {
			log.Error(err, "unable to shutdown HTTP server properly")
		}

		err = lis.Close()
		if err != nil {
			log.Error(err, "unable to close gRPC server properly")
		}
	}()

	return <-errChan
}

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	return m.Add(&grpcRunner{client: m.GetClient()})
}
