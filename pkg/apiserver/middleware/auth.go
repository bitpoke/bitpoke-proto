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

	oidc "github.com/coreos/go-oidc"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/presslabs/dashboard/pkg/cmd/apiserver/options"
)

var (
	// AuthTokenContextKey is the context key for auth token
	AuthTokenContextKey = contextKey("auth-token")
)

// Claims are the custom claims
type Claims struct {
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}

// Auth verifies the authentication token present in the gRPC request
// context
func Auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	provider, err := oidc.NewProvider(ctx, options.OIDCIssuer)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "oidc provider error: %v", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: options.ClientID})

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	// Extract custom claims
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	newCtx := context.WithValue(ctx, AuthTokenContextKey, claims)
	return newCtx, nil
}
