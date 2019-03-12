/*
Copyright 2019 Pressinfra SRL
This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package header

import (
	"context"

	"github.com/presslabs/dashboard/pkg/apiserver/internal/status"

	"google.golang.org/grpc/metadata"
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

// OrgFromContext returns organzation id from context
func OrgFromContext(ctx context.Context) string {
	md, hasMD := metadata.FromIncomingContext(ctx)
	if !hasMD {
		panic(status.Unauthenticatedf("no organization id value in context"))
	}
	if val, ok := md[organizationTokenContextKey]; ok {
		return val[0]
	}
	panic(status.Unauthenticatedf("no organization id value in context"))
}
