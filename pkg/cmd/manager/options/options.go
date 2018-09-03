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
	"net/url"

	"github.com/spf13/pflag"
)

// GitBaseDomainURL is the the base domain used to obtain the git repo domain for projects
var GitBaseDomainURL = "git.presslabs.net"

// AddToFlagSet add options to a FlagSet
func AddToFlagSet(flag *pflag.FlagSet) {
	flag.StringVar(&GitBaseDomainURL, "git-base-domain", GitBaseDomainURL, "The base git domain")
}

// Validate validates the arguments
func Validate() error {
	_, err := url.Parse(GitBaseDomainURL)
	if err != nil {
		return err
	}
	return nil
}
