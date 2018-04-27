/*
Copyright 2014 The Kubernetes Authors.

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

// Package app implements a controller manager that runs a set of active
// controllers, like projects controller, site controller or edge controller
package app

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/appscode/kutil/tools/clientcmd"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"

	dashclientset "github.com/presslabs/dashboard/pkg/client/clientset/versioned"
	dashintscheme "github.com/presslabs/dashboard/pkg/client/clientset/versioned/scheme"
	dashinformers "github.com/presslabs/dashboard/pkg/client/informers/externalversions"

	"github.com/presslabs/dashboard/cmd/controller/app/options"
	"github.com/presslabs/dashboard/pkg/controller"
	"github.com/presslabs/dashboard/pkg/version"
)

const (
	controllerAgentName = "presslabs-dashboard-controller"
)

// NewControllerManagerCommand creates a *cobra.Command object with default parameters
func NewControllerManagerCommand(stopCh <-chan struct{}) *cobra.Command {
	o := options.NewControllerManagerOptions()
	cmd := &cobra.Command{
		Use:   "presslabs-controller",
		Short: fmt.Sprintf("Presslabs Dashboard Controller (%s)", version.Get()),
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
	glog.Infof("Starting Presslabs Dashboard Controller (%s)...", version.Get())
	ctx, err := buildControllerContext(c)

	if err != nil {
		return err
	}

	run := func(_ <-chan struct{}) {
		var wg sync.WaitGroup
		glog.V(4).Infof("Starting shared informer factories")
		ctx.KubeSharedInformerFactory.Start(stopCh)
		ctx.DashboardSharedInformerFactory.Start(stopCh)
		wg.Wait()
		glog.Fatalf("Control loops exited")
	}

	run(stopCh)
	return nil

	panic("unreachable")
}

func buildControllerContext(c *options.ControllerManagerOptions) (*controller.Context, error) {
	// Create a Kubernetes api client
	kubeCfg, err := clientcmd.BuildConfigFromContext(c.Kubeconfig, "")

	// Create a Kubernetes api client
	cl, err := kubernetes.NewForConfig(kubeCfg)

	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes client: %s", err.Error())
	}

	// Create a Navigator api client
	intcl, err := dashclientset.NewForConfig(kubeCfg)

	if err != nil {
		return nil, fmt.Errorf("error creating internal group client: %s", err.Error())
	}

	// Create event broadcaster
	// Add oxygen types to the default Kubernetes Scheme so Events can be
	// logged properly
	dashintscheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.V(4).Infof)
	eventBroadcaster.StartRecordingToSink(&core.EventSinkImpl{Interface: cl.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: controllerAgentName})

	kubeSharedInformerFactory := informers.NewFilteredSharedInformerFactory(cl, time.Second*30, "", nil)
	dashboardInformerFactory := dashinformers.NewFilteredSharedInformerFactory(intcl, time.Second*30, "", nil)
	return &controller.Context{
		KubeClient:                     cl,
		KubeSharedInformerFactory:      kubeSharedInformerFactory,
		Recorder:                       recorder,
		DashboardClient:                intcl,
		DashboardSharedInformerFactory: dashboardInformerFactory,
	}, nil
}
