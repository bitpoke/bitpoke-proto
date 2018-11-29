/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package controller

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToServerFuncs is a list of functions to add all Controllers to the
// Manager and to the API server
var AddToServerFuncs []func(m manager.Manager, auth grpc_auth.AuthFunc, grpcAddr, httpAddr string) error

// AddToServer adds all Controllers to the Manager and to the API server
func AddToServer(m manager.Manager, auth grpc_auth.AuthFunc, grpcAddr, httpAddr string) error {
	for _, f := range AddToServerFuncs {
		if err := f(m, auth, grpcAddr, httpAddr); err != nil {
			return err
		}
	}
	return nil
}
