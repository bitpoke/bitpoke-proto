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
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
	"github.com/presslabs/dashboard/pkg/controller"
	"github.com/presslabs/dashboard/pkg/webhook"
)

// controllerManagerCmd represents the controllerManager command
var controllerManagerCmd = &cobra.Command{
	Use:     "controller-manager",
	Short:   "Start the Presslabs Dashboard Kubernetes controllers",
	Run:     runControllerManager,
	PreRunE: func(cmd *cobra.Command, args []string) error { return options.Validate() },
}

var runControllerManager = func(cmd *cobra.Command, args []string) {
	log = logf.Log.WithName("controller-manager")
	log.Info("Starting Presslabs Dashboard controller-manager...")

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		log.Error(err, "unable to create a new manager")
		os.Exit(1)
	}

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable to register types to scheme")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "unable to setup controllers")
		os.Exit(1)
	}

	// Setup webhook server
	if err := webhook.AddToManager(mgr); err != nil {
		log.Error(err, "unable to setup webhook server")
		os.Exit(1)
	}

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to start the manager")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(controllerManagerCmd)
	options.AddToFlagSet(controllerManagerCmd.Flags())
}
