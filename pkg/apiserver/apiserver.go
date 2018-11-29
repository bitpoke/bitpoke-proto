/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package apiserver

import (
	"context"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"reflect"

	"github.com/gobuffalo/packr"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/presslabs/dashboard/pkg/cmd/apiserver/options"
)

// APIServer is the API Server that contains GRPC Server, HTTP Server and client
type APIServer struct {
	Client     client.Client
	GRPCServer *grpc.Server
	GRPCAddr   string
	HTTPServer *http.Server
	serveMux   *http.ServeMux
}

type config struct {
	OIDCIssuer   string `jsenv:"REACT_APP_OIDC_ISSUER"`
	ClientID     string `jsenv:"REACT_APP_OIDC_CLIENT_ID"`
	GRPCProxyURL string `jsenv:"REACT_APP_API_URL"`
}

var log = logf.Log.WithName("apiserver")

func serveConfig(resp http.ResponseWriter, req *http.Request) {
	cfg := config{
		OIDCIssuer:   options.OIDCIssuer,
		ClientID:     options.ClientID,
		GRPCProxyURL: options.GRPCProxyURL,
	}

	resp.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(resp, "window.env = {};\n")
	_cfg := reflect.ValueOf(cfg)
	for i := 0; i < _cfg.NumField(); i++ {
		field := _cfg.Type().Field(i)
		value := _cfg.Field(i)
		varName := field.Tag.Get("jsenv")
		if len(varName) > 0 {
			switch value.Kind() {
			case reflect.String:
				fmt.Fprintf(resp, "window.env.%s = \"%s\";\n", varName, template.JSEscapeString(value.String()))
			default:
				panic("jsenv must be strings")
			}
		}
	}
}

func (s *APIServer) startGRPCServer() error {
	lis, err := net.Listen("tcp", s.GRPCAddr)
	if err != nil {
		return err
	}

	log.Info("gRPC Server listening", "address", options.GRPCAddr)
	err = s.GRPCServer.Serve(lis)
	return err
}

func (s *APIServer) startHTTPServer() error {
	box := packr.NewBox("../../app/build")
	// if !box.Has("index.html") {
	// panic("Cannot find 'index.html' web server entry point. You need to build the webapp first.")
	// }

	log.Info("Web Server listening", "address", s.HTTPServer.Addr)
	s.serveMux.Handle("/", http.FileServer(box))
	if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Start starts API server
func (s *APIServer) Start(stop <-chan struct{}) error {

	errChan := make(chan error, 2)

	go func() {
		err := s.startGRPCServer()
		errChan <- err
	}()

	go func() {
		err := s.startHTTPServer()
		errChan <- err
	}()

	go func() {
		<-stop
		s.GracefullShutdown()
	}()

	return <-errChan
}

// GracefullShutdown will gracefully shutdown the API Server
func (s *APIServer) GracefullShutdown() {
	if err := s.HTTPServer.Shutdown(context.TODO()); err != nil {
		log.Error(err, "unable to shutdown HTTP server properly")
	}

	s.GRPCServer.GracefulStop()
}

// AddToServer adds all Controllers to the Manager and to the API server
func AddToServer(m manager.Manager, auth grpc_auth.AuthFunc, grpcAddr, httpAddr string) (*APIServer, error) {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth)),
	)

	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	sMux := http.NewServeMux()
	handler := func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		} else if req.URL.Path == "/env.js" {
			serveConfig(resp, req)
		} else {
			sMux.ServeHTTP(resp, req)
		}
	}

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: http.HandlerFunc(handler),
	}

	apiServer := &APIServer{
		Client:     m.GetClient(),
		GRPCServer: grpcServer,
		GRPCAddr:   grpcAddr,
		HTTPServer: httpServer,
		serveMux:   sMux,
	}

	// register reflection service on gRPC server
	reflection.Register(grpcServer)

	return apiServer, m.Add(apiServer)
}
