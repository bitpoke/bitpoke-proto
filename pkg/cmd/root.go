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
	goflag "flag"
	"fmt"
	"os"

	// homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
	"github.com/go-logr/logr"
	// enable GKE cluster login
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	logf "github.com/presslabs/controller-util/log"
)

var cfg *rest.Config
var log logr.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "presslabs-dashboard",
	Short: "Presslabs Dashboard for WordPress",
	Long: `Presslabs Dashboard for WordPress is a Kubernetes controller and Web
application for managing WordPress deployments at scale.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// setup logging
		logf.SetLogger(logf.ZapLogger(true))

		// configure Kubernetes rest.Client
		cfg, err = config.GetConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// add goflag flags
	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	// remove glog inserted flags
	// nolint: gosec
	// we really don't care about these errors
	_ = rootCmd.PersistentFlags().MarkHidden("alsologtostderr")
	_ = rootCmd.PersistentFlags().MarkHidden("log_backtrace_at")
	_ = rootCmd.PersistentFlags().MarkHidden("log_dir")
	_ = rootCmd.PersistentFlags().MarkHidden("logtostderr")
	_ = rootCmd.PersistentFlags().MarkHidden("stderrthreshold")
	_ = rootCmd.PersistentFlags().MarkHidden("v")
	_ = rootCmd.PersistentFlags().MarkHidden("vmodule")
}
