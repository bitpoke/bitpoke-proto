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

package options

import (
	"github.com/spf13/pflag"
)

// GRPCAddr is the address to bind the gRPC server on
var GRPCAddr = ":9090"

// HTTPAddr is the address to bind the gRPC web proxy server on
var HTTPAddr = ":8080"

// AddToFlagSet add options to a FlagSet
func AddToFlagSet(flag *pflag.FlagSet) {
	flag.StringVar(&GRPCAddr, "grpc-addr", GRPCAddr, "GRPC address to use")
	flag.StringVar(&HTTPAddr, "http-addr", HTTPAddr, "HTTP address to use")
}
