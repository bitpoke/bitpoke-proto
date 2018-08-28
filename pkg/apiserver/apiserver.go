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
	// "encoding/base64"
	"fmt"
	"log"
	"net/http"
	// "os"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	// "github.com/auth0-community/auth0"
	project "github.com/presslabs/dashboard/pkg/apiserver/projects/v1"
	// "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type grpcRunner struct {
	client client.Client
}

func (s *grpcRunner) Start(stop <-chan struct{}) error {
	var (
		httpServer http.Server
		port       = 9090
	)

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(handleAuthentication)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(handleAuthentication)),
	)
	project.RegisterProjectsServer(grpcServer, project.NewProjectServer(s.client))

	wrappedServer := grpcweb.WrapServer(grpcServer)

	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}

	resources := grpcweb.ListGRPCResources(grpcServer)

	httpServer = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}

	go func() {
		<-stop
		httpServer.Shutdown(context.TODO())
	}()

	log.Printf("Server started on http://0.0.0.0:%d", port)
	log.Printf("Available resources: %v", resources)

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
		return nil, err
	}

	log.Println("string token>>>> %v", token)

	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	log.Println("TOKKEEENN>>>> %v", parsedToken)

	// secret, _ := base64.URLEncoding.DecodeString(os.Getenv("AUTH0_CLIENT_SECRET"))
	// secretProvider := auth0.NewKeyProvider(secret)
	// audience := os.Getenv("AUTH0_CLIENT_ID")

	// configuration := auth0.NewConfiguration(secretProvider, []string{audience}, os.Getenv("AUTH0_DOMAIN"), jose.HS256)
	// validator := auth0.NewValidator(configuration, nil)

	// token, err := validator.ValidateRequest(r)

	// if err != nil {
	// 	fmt.Println("Token is not valid:", token)
	// }

	// tokenInfo, err := parseToken(token)
	// if err != nil {
	// 	return nil, grpc.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	// }
	// grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))
	// newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
	return ctx, nil
}

// func validateToken

// // NewValidator creates a new
// // validator with the provided configuration.
// func NewValidator(config Configuration, extractor RequestTokenExtractor) *JWTValidator {
// 	if extractor == nil {
// 		extractor = RequestTokenExtractorFunc(FromHeader)
// 	}
// 	return &JWTValidator{config, extractor}
// }

// // ValidateRequest validates the token within
// // the http request.
// func (v *JWTValidator) ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error) {
// 	token, err := v.extractor.Extract(r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(token.Headers) < 1 {
// 		return nil, ErrNoJWTHeaders
// 	}

// 	// trust secret provider when sig alg not configured and skip check
// 	if v.config.signIn != "" {
// 		header := token.Headers[0]
// 		if header.Algorithm != string(v.config.signIn) {
// 			return nil, ErrInvalidAlgorithm
// 		}
// 	}

// 	claims := jwt.Claims{}
// 	key, err := v.config.secretProvider.GetSecret(r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = token.Claims(key, &claims); err != nil {
// 		return nil, err
// 	}

// 	expected := v.config.expectedClaims.WithTime(time.Now())
// 	err = claims.Validate(expected)
// 	return token, err
// }

// // Claims unmarshall the claims of the provided token
// func (v *JWTValidator) Claims(r *http.Request, token *jwt.JSONWebToken, values ...interface{}) error {
// 	key, err := v.config.secretProvider.GetSecret(r)
// 	if err != nil {
// 		return err
// 	}
// 	return token.Claims(key, values...)
// }
