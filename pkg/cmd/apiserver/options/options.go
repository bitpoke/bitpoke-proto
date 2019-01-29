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
	"os"

	"github.com/spf13/pflag"
)

// ClientID is the OpenID connect client id
var ClientID = ""

// ClientSecret is the OpenID connect client id
var ClientSecret = ""

// OIDCIssuer is the OpenID Issuer
var OIDCIssuer = ""

// BaseURL is the base url for the webapp
var BaseURL = "http://localhost:8080"

// GRPCProxyURL is the url for the gRPC proxy
var GRPCProxyURL = "http://localhost:8080"

// GRPCAddr is the address to bind the gRPC server on
var GRPCAddr = ":9090"

// HTTPAddr is the address to bind the HTTP server and gRPC proxy run
var HTTPAddr = ":8080"

// AddToFlagSet add options to a FlagSet
func AddToFlagSet(flag *pflag.FlagSet) {
	flag.StringVar(&ClientID, "oidc-client-id", os.Getenv("OIDC_CLIENT_ID"), "OpenID Cliet ID")
	flag.StringVar(&ClientSecret, "oidc-client-secret", os.Getenv("OIDC_CLIENT_SECRET"), "OpenID Cliet Secret")
	flag.StringVar(&OIDCIssuer, "oidc-issuer", os.Getenv("OIDC_ISSUER"), "The audience to validate JWT against")
	flag.StringVar(&BaseURL, "base-url", BaseURL, "Base URL for the webapp")
	flag.StringVar(&GRPCProxyURL, "grpc-proxy-url", GRPCProxyURL, "URL for the gRPC proxy")
	flag.StringVar(&GRPCAddr, "grpc-addr", GRPCAddr, "gRPC server address")
	flag.StringVar(&HTTPAddr, "http-addr", HTTPAddr, "web server address")
}

// LoadFromEnv fills in unset configs from environment variables
func LoadFromEnv() {
	if len(ClientID) == 0 {
		ClientID = os.Getenv("CLIENT_ID")
	}
	if len(ClientSecret) == 0 {
		ClientSecret = os.Getenv("CLIENT_SECRET")
	}
}
