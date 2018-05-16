/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package apiserver

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Presslabs Dashboard API Server!\npath: %s\n", r.URL.Path)
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Presslabs Dashboard API Server!\npath: %s\n", r.URL.Path)
}
