/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package defaultserver

import (
	"fmt"

	"github.com/presslabs/dashboard/pkg/webhook/default_server/project/validating"
)

func init() {
	for k, v := range validating.Builders {
		_, found := builderMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in builder map: %v", k))
		}
		builderMap[k] = v
	}
	for k, v := range validating.HandlerMap {
		_, found := HandlerMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in handler map: %v", k))
		}
		_, found = builderMap[k]
		if !found {
			log.V(1).Info(fmt.Sprintf(
				"can't find webhook builder name %q in builder map", k))
			continue
		}
		HandlerMap[k] = v
	}
}
