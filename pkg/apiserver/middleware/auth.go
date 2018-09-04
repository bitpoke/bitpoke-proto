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

package middleware

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/presslabs/dashboard/pkg/apiserver/jwks"
	jose "github.com/square/go-jose"
	"github.com/square/go-jose/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Auth verifies the authentication token present in the gRPC request context
func Auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	validatedToken, err := validateToken(parsedToken)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	newCtx := context.WithValue(ctx, "token", validatedToken)
	return newCtx, nil
}

func validateToken(token *jwt.JSONWebToken) (*jwt.Claims, error) {
	audience := os.Getenv("AUTH0_CLIENT_ID")
	issuer := fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))
	alg := jose.RS256
	expectedClaims := jwt.Expected{Issuer: issuer, Audience: []string{audience}}

	if len(token.Headers) < 1 {
		return nil, errors.New("No headers in the token")
	}

	header := token.Headers[0]

	if header.Algorithm != string(alg) {
		return nil, errors.New("invalid algorithm")
	}

	jwksClient, _ := jwks.NewClient(fmt.Sprintf("%s.well-known/jwks.json", issuer))
	key, err := jwksClient.GetKey(header.KeyID)
	if err != nil {
		return nil, err
	}

	claims := jwt.Claims{}
	if err = token.Claims(key, &claims); err != nil {
		return nil, err
	}

	expected := expectedClaims.WithTime(time.Now())
	err = claims.Validate(expected)
	return &claims, err
}
