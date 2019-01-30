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

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/presslabs/dashboard/pkg/apis"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/apiserver/controller"
	"github.com/presslabs/dashboard/pkg/cmd/apiserver/options"
)

// apiserverCmd represents the controllerManager command
var apiserverCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "Start the Presslabs Dashboard API server",
	Run:   runAPIServer,
}

var runAPIServer = func(cmd *cobra.Command, args []string) {
	options.LoadFromEnv()
	log = logf.Log.WithName("apiserver")
	log.Info("Starting Presslabs Dashboard apiserver...")

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		log.Error(err, "unable to create a new manager")
		os.Exit(1)
	}

	// Setup Scheme for all resources
	err = apis.AddToScheme(mgr.GetScheme())
	if err != nil {
		log.Error(err, "unable to register types to scheme")
		os.Exit(1)
	}

	opts := &apiserver.APIServerOptions{
		Manager:  mgr,
		HTTPAddr: options.HTTPAddr,
		GRPCAddr: options.GRPCAddr,
	}

	server, err := apiserver.NewAPIServer(opts)
	if err != nil {
		log.Error(err, "unable to setup apiserver")
		os.Exit(1)
	}

	err = controller.AddToServer(server)
	if err != nil {
		log.Error(err, "unable to setup apiserver controllers")
		os.Exit(1)
	}

	err = mgr.Add(server)
	if err != nil {
		log.Error(err, "unable to add to managaer")
		os.Exit(1)
	}

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to start the manager")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(apiserverCmd)
	options.AddToFlagSet(apiserverCmd.Flags())
}
