/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package controller

import (
	"github.com/presslabs/dashboard/pkg/apiserver"
)

// AddToServerFuncs is a list of functions to add all Controllers to the
// Manager and to the API server
var AddToServerFuncs []func(server *apiserver.APIServer) error

// AddToServer adds all Controllers to the Manager and to the API server
func AddToServer(server *apiserver.APIServer) error {
	for _, f := range AddToServerFuncs {
		if err := f(server); err != nil {
			return err
		}
	}
	return nil
}
