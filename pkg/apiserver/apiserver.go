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
	// if r.Method == "POST" {
	// 	decoder := json.NewDecoder(r.Body)
	//
	// 	type Message struct {
	// 		Name, Text string
	// 	}
	// 	var m Message
	// 	err := decoder.Decode(&m)
	//
	// 	clientConfig, err := config.createClientConfigFromFile()
	// 	if err != nil {
	// 		glog.Fatalf("Failed to create a ClientConfig: %v. Exiting.", err)
	// 	}
	//
	// 	clientset, err := clientset.NewForConfig(clientConfig)
	// 	if err != nil {
	// 		glog.Fatalf("Failed to create a ClientSet: %v. Exiting.", err)
	// 	}
	//
	// 	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
	//
	// 	_, err := clientset.Core().Namespaces().Create(nsSpec)
	// }
	fmt.Fprintf(w, "Hello from Presslabs Dashboard API Server!\npath: %s\n", r.URL.Path)
}
