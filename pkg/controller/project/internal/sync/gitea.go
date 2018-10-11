/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
	"github.com/presslabs/dashboard/pkg/internal/project"
)

var (
	giteaVersionLabels = labels.Set{
		"app.kubernetes.io/version": giteaReleaseVersion,
	}
)

func giteaDomain(o *project.Project) string {
	return fmt.Sprintf("%s.%s", o.Domain(), options.GitBaseDomainURL)
}
