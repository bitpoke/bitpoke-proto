/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package webhook

import (
	server "github.com/presslabs/dashboard/pkg/webhook/default_server"
)

func init() {
	// AddToManagerFuncs is a list of functions to create webhook servers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, server.Add)
}
