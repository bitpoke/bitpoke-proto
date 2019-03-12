/*
Copyright 2019 Pressinfra SRL
This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package metadata

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"
)

var (
	// organizationTokenContextKey is the context key for organization_id token
	organizationTokenContextKey = "organization"
)

// AddOrgInContext adds organization id in context
func AddOrgInContext(ctx context.Context, org string) context.Context {
	md := metadata.New(map[string]string{organizationTokenContextKey: org})
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return newCtx
}

// RequireOrganizationNamespace returns organzation id from context
func RequireOrganizationNamespace(ctx context.Context) string {
	md, hasMD := metadata.FromIncomingContext(ctx)
	if !hasMD {
		panic(status.InvalidArgumentf("no organization id value in context"))
	}

	val, ok := md[organizationTokenContextKey]
	if !ok || val[0] == "" {
		panic(status.InvalidArgumentf("no organization id value in context"))
	}

	return val[0]
}
