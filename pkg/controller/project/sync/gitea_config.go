/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"strconv"

	"github.com/go-ini/ini"
	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	giteaName             = "gitea"
	giteaReleaseVersion   = "1.5.0"
	giteaImage            = "docker.io/gitea/gitea:" + giteaReleaseVersion
	giteaHTTPInternalIP   = "0.0.0.0"
	giteaHTTPInternalPort = 80
	giteaSshInternalPort  = 22
	giteaReplicas         = 1
	giteaMaxUnavailable   = 1
	giteaMaxSurge         = 0
)

func GetGiteaLabels(project *dashboardv1alpha1.Project) labels.Set {
	giteaSelector := labels.Set{
		"app.kubernetes.io/name": giteaName,
	}
	return labels.Merge(project.GetDefaultLabels(), giteaSelector)
}

func GetGiteaPodLabels(project *dashboardv1alpha1.Project) labels.Set {
	giteaLabels := labels.Set{
		"app.kubernetes.io/version": giteaReleaseVersion,
	}
	return labels.Merge(GetGiteaLabels(project), giteaLabels)
}

func getExistingKeyValueOrDefault(section *ini.Section, name string, defvalue string) string {
	if section == nil || section.Key(name) == nil {
		return defvalue
	}
	return section.Key(name).Value()
}

func createGiteaConfig(project *dashboardv1alpha1.Project, values map[string]string) (*ini.File, error) {
	cfg := ini.Empty()
	default_section := cfg.Section("DEFAULT")
	default_section.NewKey("RUN_MODE", "prod")
	server := cfg.Section("server")
	server.NewKey("PROTOCOL", "http")
	// server.NewKey("ROOT_URL", "https://ingress")
	server.NewKey("HTTP_ADDR", giteaHTTPInternalIP)
	server.NewKey("HTTP_PORT", strconv.Itoa(giteaHTTPInternalPort))
	server.NewKey("DISABLE_SSH", "true")
	server.NewKey("START_SSH_SERVER", "false")
	server.NewKey("SSH_DOMAIN", "sshDomain")
	server.NewKey("SSH_PORT", strconv.Itoa(giteaSshInternalPort))
	server.NewKey("SSH_LISTEN_HOST", "0.0.0.0")
	server.NewKey("SSH_LISTEN_PORT", strconv.Itoa(giteaSshInternalPort))
	server.NewKey("OFFLINE_MODE", "false")
	server.NewKey("LANDING_PAGE", "home")
	database := cfg.Section("database")
	database.NewKey("DB_TYPE", "sqlite3")
	database.NewKey("PATH", "data/gitea.db")
	database.NewKey("SQLITE_TIMEOUT", "500")
	database.NewKey("LOG_SQL", "false")
	security := cfg.Section("security")
	security.NewKey("INSTALL_LOCK", "true")
	security.NewKey("SECRET_KEY", values["SECRET_KEY"])
	security.NewKey("INTERNAL_TOKEN", values["INTERNAL_TOKEN"])
	security.NewKey("LOGIN_REMEMBER_DAYS", "7")
	security.NewKey("COOKIE_USERNAME", "gitea_user")
	security.NewKey("COOKIE_REMEMBER_NAME", "gitea_rememberme")
	security.NewKey("REVERSE_PROXY_AUTHENTICATION_USER", "X-WEBAUTH-USER")
	security.NewKey("MIN_PASSWORD_LENGTH", "8")
	security.NewKey("IMPORT_LOCAL_PATHS", "false")
	security.NewKey("DISABLE_GIT_HOOKS", "false")
	admin := cfg.Section("admin")
	admin.NewKey("DISABLE_REGULAR_ORG_CREATION", "true")
	return cfg, nil
}
