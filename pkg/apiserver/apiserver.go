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
	"time"
	"fmt"

	"net/http"
	"errors"
	"os"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	project "github.com/presslabs/dashboard/pkg/apiserver/projects/v1"
	"github.com/presslabs/dashboard/pkg/cmd/apiserver/options"
	jose "github.com/square/go-jose"
	"github.com/square/go-jose/jwt"

  "github.com/presslabs/dashboard/pkg/apiserver/jwks"
)

type grpcRunner struct {
	client client.Client
}

var log = logf.Log.WithName("apiserver")

func (s *grpcRunner) Start(stop <-chan struct{}) error {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(handleAuthentication)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(handleAuthentication)),
	)
	project.RegisterProjectsServer(grpcServer, project.NewProjectServer(s.client))

	wrappedServer := grpcweb.WrapServer(grpcServer)

	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}

	port := options.GRPCPort
	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}

	go func() {
		<-stop
		httpServer.Shutdown(context.TODO())
	}()

	log.Info("Server listening", "port", port)

	return httpServer.ListenAndServe()
}

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	m.Add(&grpcRunner{client: m.GetClient()})
	return nil
}

func handleAuthentication(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Invalid Auth Token: %v", err)
	}

	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Invalid Auth Token: %v", err)
	}

	validatedToken, err := validateToken(parsedToken)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Invalid Auth Token: %v", err)
	}

	newCtx := context.WithValue(ctx, "token", validatedToken)
	return newCtx, nil
}

func validateToken(token *jwt.JSONWebToken) (*jwt.JSONWebToken, error) {
	audience := os.Getenv("AUTH0_CLIENT_ID")
	issuer := fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))

	alg := jose.RS256

	expectedClaims := jwt.Expected{Issuer: issuer, Audience: []string{audience}}

	if len(token.Headers) < 1 {
		return nil, errors.New("No headers in the token")
	}

	header := token.Headers[0]

	if header.Algorithm != string(alg) {
		return nil, errors.New("Invalid algorithm")
	}

  jwksClient, _ := jwks.NewClient(fmt.Sprintf("%s.well-known/jwks.json", issuer))
  key, err := jwksClient.GetKey(header.KeyID)
  if err != nil {
    log.Error(err, "Cannot get key")
  }

	claims := jwt.Claims{}
	if err := token.Claims(key, &claims); err != nil {
    log.Error(err, "cannot get claims from token")
		return nil, err
	}

	expected := expectedClaims.WithTime(time.Now())
	err = claims.Validate(expected)
	return token, err
}