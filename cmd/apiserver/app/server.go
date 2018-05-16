// Package app implements a server that runs the Presslabs Dashboard API Server
//
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/presslabs/dashboard/cmd/apiserver/app/options"
	"github.com/presslabs/dashboard/pkg/apiserver"
	"github.com/presslabs/dashboard/pkg/version"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand(stopCh <-chan struct{}) *cobra.Command {
	o := options.NewControllerManagerOptions()
	cmd := &cobra.Command{
		Use:   "presslabs-apiserver",
		Short: fmt.Sprintf("Presslabs Dashboard API Server (%s)", version.Get()),
		Run: func(cmd *cobra.Command, args []string) {

			if err := o.Validate(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := Run(o, stopCh); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	flags := cmd.Flags()
	o.AddFlags(flags)

	return cmd
}

// Run Presslabs Controller Manager.  This should never exit.
func Run(c *options.ControllerManagerOptions, stopCh <-chan struct{}) error {
	glog.Infof("Starting Presslabs Dashboard API Server (%s)...", version.Get())

	http.HandleFunc("/", apiserver.RootHandler)
	http.HandleFunc("/projects", apiserver.ProjectsHandler)

	run := func(_ <-chan struct{}) {

		srv := &http.Server{Addr: c.ListenAddr}
		go func() {
			glog.Infof("Start listening on %s", c.ListenAddr)
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				glog.Fatalf(err.Error())
			}
		}()
		<-stopCh
		glog.Infof("Shutting down Presslabs Dashboard API Server...")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		srv.Shutdown(ctx)
		glog.Infof("Presslabs Dashboard API Server gracefully stopped.")
		os.Exit(0)
	}

	run(stopCh)

	panic("unreachable")
}
