/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"sort"
	"strconv"

	"github.com/go-ini/ini"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

const (
	giteaName            = "gitea"
	giteaVersion         = "1.5.2"
	giteaImage           = "docker.io/gitea/gitea:" + giteaVersion
	giteaHTTPInternalIP  = "0.0.0.0"
	giteaHTTPPort        = 8080
	giteaSSHPort         = 22
	giteaRequestsMemory  = "256Mi"
	giteaRequestsCPU     = "100m"
	giteaRequestsStorage = "10Gi"
)

func createGiteaConfig(project *projectns.ProjectNamespace, data map[string][]byte) (*ini.File, error) {
	cfg := ini.Empty()

	config := map[string]map[string]string{
		"DEFAULT": {
			"RUN_MODE": "prod",
		},
		"server": {
			"PROTOCOL":         "http",
			"DOMAIN":           giteaDomain(project),
			"HTTP_ADDR":        giteaHTTPInternalIP,
			"HTTP_PORT":        strconv.Itoa(giteaHTTPPort),
			"DISABLE_SSH":      "true",
			"START_SSH_SERVER": "false",
			"SSH_DOMAIN":       giteaDomain(project),
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
			"SECRET_KEY":                        string(data["SECRET_KEY"]),
			"INTERNAL_TOKEN":                    string(data["INTERNAL_TOKEN"]),
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

	sections := []string{}
	for section := range config {
		sections = append(sections, section)
	}
	sort.Strings(sections)

	for _, sectionName := range sections {
		section := cfg.Section(sectionName)
		keys := []string{}
		for key := range config[sectionName] {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			if _, err := section.NewKey(key, config[sectionName][key]); err != nil {
				return nil, err
			}
		}
	}

	return cfg, nil
}
