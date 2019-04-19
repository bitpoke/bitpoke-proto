/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package impersonate

import (
	"context"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/apiserver/status"
)

// Client returns impersonated client
func Client(userName string, cfg *rest.Config) client.Client {
	if userName == "" {
		panic(status.InternalErrorf("empty impersonation user"))
	}
	mcfg := rest.CopyConfig(cfg)
	mcfg.Impersonate = rest.ImpersonationConfig{
		UserName: userName,
	}
	c, err := client.New(mcfg, client.Options{})
	if err != nil {
		panic(err)
	}
	return c
}

// ClientFromContext returns user ID from context and impersonated client
func ClientFromContext(ctx context.Context, cfg *rest.Config) (client.Client, string) {
	userID := metadata.RequireUserID(ctx)
	return Client(userID, cfg), userID
}
