package options

import (
	"github.com/spf13/pflag"
)

type ControllerManagerOptions struct {
	Kubeconfig string
}

func NewControllerManagerOptions() *ControllerManagerOptions {
	return &ControllerManagerOptions{
		Kubeconfig: "",
	}
}

func (o *ControllerManagerOptions) Validate() error {
	return nil
}

func (o *ControllerManagerOptions) AddFlags(*pflag.FlagSet) {
}
