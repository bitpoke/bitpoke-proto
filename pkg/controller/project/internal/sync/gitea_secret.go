/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"bytes"
	"encoding/base64"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/util/rand"
)

// NewGiteaSecretSyncer returns a new syncer.Interface for reconciling Gitea Secret
func NewGiteaSecretSyncer(proj *dashboardv1alpha1.Project) syncer.Interface {
	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      giteaSecretName(proj),
			Namespace: getNamespaceName(proj),
		},
	}

	return syncer.New("GiteaSecret", proj, obj, func(existing runtime.Object) error {
		out := existing.(*corev1.Secret)
		out.Labels = giteaPodLabels(proj)

		secretKeyBytes, ok := out.Data["SECRET_KEY"]
		if !ok {
			secretKeyBytes = rand.GenerateRandomBytes(32)
		}

		internalTokenBytes, ok := out.Data["INTERNAL_TOKEN"]
		if !ok {
			internalTokenBytes = rand.GenerateRandomBytes(64)
		}

		secretKey := base64.URLEncoding.EncodeToString(secretKeyBytes)
		internalToken := base64.URLEncoding.EncodeToString(internalTokenBytes)

		secrets := map[string]string{
			"SECRET_KEY":     secretKey,
			"INTERNAL_TOKEN": internalToken,
		}
		cfg, err := createGiteaConfig(proj, secrets)

		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if _, err := cfg.WriteTo(&buf); err != nil {
			log.Error(err, "unable to load existing Gitea settings", "project", proj.Name)
		}

		secrets["app.ini"] = buf.String()
		out.StringData = secrets

		return nil
	})
}
