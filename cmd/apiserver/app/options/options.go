package options

import (
	"github.com/spf13/pflag"
)

const (
	defaultListenAddr = ":8080"
)

type ControllerManagerOptions struct {
	ListenAddr string
}

func NewControllerManagerOptions() *ControllerManagerOptions {
	return &ControllerManagerOptions{
		ListenAddr: defaultListenAddr,
	}
}

func (o *ControllerManagerOptions) Validate() error {
	return nil
}

func (o *ControllerManagerOptions) AddFlags(*pflag.FlagSet) {
	pflag.StringVar(&o.ListenAddr, "address", defaultListenAddr, "Address for the API Server to listen on")
}
