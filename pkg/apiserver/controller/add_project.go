/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package controller

import "github.com/presslabs/dashboard/pkg/apiserver/controller/project"

func init() {
	// AddToServerFuncs is a list of functions to create controllers and add
	// them to the api server
	AddToServerFuncs = append(AddToServerFuncs, project.Add)
}
