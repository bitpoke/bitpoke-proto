/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"strconv"

	"github.com/go-ini/ini"
	"k8s.io/apimachinery/pkg/labels"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	giteaName            = "gitea"
	giteaReleaseVersion  = "1.5.0"
	giteaImage           = "docker.io/gitea/gitea:" + giteaReleaseVersion
	giteaHTTPInternalIP  = "0.0.0.0"
	giteaHTTPPort        = 80
	giteaSSHPort         = 22
	giteaRequestsMemory  = "512Mi"
	giteaRequestsCPU     = "100m"
	giteaRequestsStorage = "10Gi"
)

// GetGiteaLabels returns a set of labels that can be used to identify Gitea related resources
func GetGiteaLabels(project *dashboardv1alpha1.Project) labels.Set {
	giteaSelector := labels.Set{
		"app.kubernetes.io/name": giteaName,
	}
	return labels.Merge(project.GetDefaultLabels(), giteaSelector)
}

// GetGiteaPodLabels returns a set of labels that should be applied on Gitea related objects that are managed by the project controller
func GetGiteaPodLabels(project *dashboardv1alpha1.Project) labels.Set {
	giteaLabels := labels.Set{
		"app.kubernetes.io/version": giteaReleaseVersion,
	}
	return labels.Merge(GetGiteaLabels(project), giteaLabels)
}

func createGiteaConfig(project *dashboardv1alpha1.Project, values map[string]string) (*ini.File, error) {
	config := map[string]map[string]string{
		"DEFAULT": {
			"RUN_MODE": "prod",
		},
		"server": {
			"PROTOCOL":         "http",
			"DOMAIN":           project.GetGiteaDomain(),
			"HTTP_ADDR":        giteaHTTPInternalIP,
			"HTTP_PORT":        strconv.Itoa(giteaHTTPPort),
			"DISABLE_SSH":      "true",
			"START_SSH_SERVER": "false",
			"SSH_DOMAIN":       project.GetGiteaDomain(),
			"SSH_PORT":         strconv.Itoa(giteaSSHPort),
			"SSH_LISTEN_HOST":  "0.0.0.0",
			"SSH_LISTEN_PORT":  strconv.Itoa(giteaSSHPort),
			"OFFLINE_MODE":     "false",
			"LANDING_PAGE":     "home",
		},
		"database": {
			"DB_TYPE":        "sqlite3",
			"PATH":           "data/gitea.db",
			"SQLITE_TIMEOUT": "500",
			"LOG_SQL":        "false",
		},
		"security": {
			"INSTALL_LOCK":                      "true",
			"SECRET_KEY":                        values["SECRET_KEY"],
			"INTERNAL_TOKEN":                    values["INTERNAL_TOKEN"],
			"LOGIN_REMEMBER_DAYS":               "7",
			"COOKIE_USERNAME":                   "gitea_user",
			"COOKIE_REMEMBER_NAME":              "gitea_rememberme",
			"REVERSE_PROXY_AUTHENTICATION_USER": "X-WEBAUTH-USER",
			"MIN_PASSWORD_LENGTH":               "8",
			"IMPORT_LOCAL_PATHS":                "false",
			"DISABLE_GIT_HOOKS":                 "false",
		},
		"admin": {
			"DISABLE_REGULAR_ORG_CREATION": "true",
		},
	}

	iniConfig := ini.Empty()

	for sectionName, keyValue := range config {
		section := iniConfig.Section(sectionName)
		for key, value := range keyValue {
			_, err := section.NewKey(key, value)
			if err != nil {
				return nil, err
			}
		}
	}
	return iniConfig, nil
}
