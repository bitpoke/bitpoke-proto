/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/handlers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	grpcstatus "google.golang.org/grpc/status"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	logf "github.com/presslabs/controller-util/log"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/cmd/apiserver/options"
)

// APIServerOptions contains manager, GRPC address, HTTP address and AuthFunc
// nolint: golint
type APIServerOptions struct {
	Manager  manager.Manager
	GRPCAddr string
	HTTPAddr string
	AuthFunc grpc_auth.AuthFunc
}

// APIServer is the API Server that contains GRPC Server, HTTP Server and client
type APIServer struct {
	Manager    manager.Manager
	GRPCServer *grpc.Server
	HTTPServer *http.Server
	serverMux  *http.ServeMux
	grpcWeb    *grpcweb.WrappedGrpcServer
	grpcAddr   string
	packrBox   *packr.Box
	tmpl       *template.Template
	env        *webappEnv
}

type webappEnv struct {
	ClientID   string `json:"REACT_APP_OIDC_CLIENT_ID"`
	OIDCIssuer string `json:"REACT_APP_OIDC_ISSUER"`
	BaseURL    string `json:"REACT_APP_API_URL"`
}

var log = logf.Log.WithName("apiserver")

func defaultOpts(opts *APIServerOptions) *APIServerOptions {
	if opts.AuthFunc == nil {
		opts.AuthFunc = metadata.Auth
	}
	return opts
}

// NewAPIServer creates a new API Server
func NewAPIServer(opts *APIServerOptions) (*APIServer, error) {
	opts = defaultOpts(opts)

	// create recovery function which keeps the eventual grpc error code
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			type grpcStatus interface {
				GRPCStatus() *grpcstatus.Status
			}
			switch v := p.(type) {
			case grpcStatus:
				return v.GRPCStatus().Err()
			default:
				return grpcstatus.Errorf(codes.Unknown, "panic triggered: %v", v)
			}
		}),
	}

	// Create the gRPC server
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zap.L()),
			grpc_auth.UnaryServerInterceptor(opts.AuthFunc),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(zap.L()),
			grpc_auth.StreamServerInterceptor(opts.AuthFunc),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	// register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Create the HTTP handler and server
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    opts.HTTPAddr,
		Handler: http.HandlerFunc(mux.ServeHTTP),
	}
	box := packr.NewBox("../../app/build")
	if !box.Has("index.html") {
		return nil, fmt.Errorf("Cannot find 'index.html' web server entry point. You need to build the webapp first.")
	}

	t := template.New("index")
	t.Funcs(template.FuncMap{
		"toJSON": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	})
	if _, err := t.Parse(box.String("index.html")); err != nil {
		return nil, err
	}

	env := webappEnv{
		ClientID:   options.ClientID,
		OIDCIssuer: options.OIDCIssuer,
		BaseURL:    options.BaseURL,
	}

	s := &APIServer{
		Manager:    opts.Manager,
		GRPCServer: grpcServer,
		HTTPServer: httpServer,
		serverMux:  mux,
		grpcAddr:   opts.GRPCAddr,
		packrBox:   &box,
		tmpl:       t,
		env:        &env,
	}

	s.grpcWeb = grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(s.allowedOrigin))
	handler := handlers.CustomLoggingHandler(os.Stderr, s, s.logFormatter)
	mux.Handle("/", handler)
	return s, nil
}

func (s *APIServer) logFormatter(_ io.Writer, params handlers.LogFormatterParams) {
	msg := fmt.Sprintf("%s %s", params.Request.Method, params.Request.URL)
	log.V(1).Info(msg, "code", params.StatusCode, "length", params.Size, "duration", time.Since(params.TimeStamp))
}

// GetGRPCAddr returns the GRPC address
func (s *APIServer) GetGRPCAddr() string {
	return s.grpcAddr
}

// GetHTTPAddr returns the HTTP address
func (s *APIServer) GetHTTPAddr() string {
	return s.HTTPServer.Addr
}

func (s *APIServer) startGRPCServer() error {
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return err
	}

	log.Info("gRPC Server listening", "address", options.GRPCAddr)
	err = s.GRPCServer.Serve(lis)
	return err
}

func (s *APIServer) allowedOrigin(origin string) bool {
	return true
}

func (s *APIServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		s.tmpl.Execute(resp, map[string]interface{}{"Env": s.env})
	} else if s.grpcWeb.IsGrpcWebRequest(req) || s.grpcWeb.IsAcceptableGrpcCorsRequest(req) {
		s.grpcWeb.ServeHTTP(resp, req)
	} else {
		http.FileServer(s.packrBox).ServeHTTP(resp, req)
	}
}

func (s *APIServer) startHTTPServer() error {
	log.Info("Web Server listening", "address", s.HTTPServer.Addr)
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
