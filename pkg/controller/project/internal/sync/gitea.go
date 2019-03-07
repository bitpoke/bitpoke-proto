/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"fmt"

	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

func giteaDomain(o *projectns.ProjectNamespace) string {
	return fmt.Sprintf("%s.%s", o.Domain(), options.GitBaseDomainURL)
}
