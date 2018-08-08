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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
	"github.com/presslabs/dashboard/pkg/util/rand"
)

const (
	GiteaSecretFailed  EventReason = "GiteaSecretFailed"
	GiteaSecretUpdated EventReason = "GiteaSecretUpdated"
)

type GiteaSecretSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.Secret
}

func NewGiteaSecretSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) *GiteaSecretSyncer {
	return &GiteaSecretSyncer{
		scheme:   r,
		existing: &corev1.Secret{},
		p:        p,
		key:      p.GetGiteaSecretKey(),
	}
}

func (s *GiteaSecretSyncer) GetKey() types.NamespacedName                 { return s.key }
func (s *GiteaSecretSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func (s *GiteaSecretSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.Secret)
	out.Labels = GetGiteaPodLabels(s.p)

	secret_key_bytes, ok := out.Data["SECRET_KEY"]
	if !ok {
		secret_key_bytes = rand.GenerateRandomBytes(32)
	}

	internal_token_bytes, ok := out.Data["INTERNAL_TOKEN"]
	if !ok {
		internal_token_bytes = rand.GenerateRandomBytes(64)
	}

	secret_key := base64.URLEncoding.EncodeToString(secret_key_bytes)
	internal_token := base64.URLEncoding.EncodeToString(internal_token_bytes)

	secrets := map[string]string{
		"SECRET_KEY":     secret_key,
		"INTERNAL_TOKEN": internal_token,
	}
	cfg, err := createGiteaConfig(s.p, secrets)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := cfg.WriteTo(&buf); err != nil {
		log.Error(err, "unable to load existing Gitea settings", "project", s.p.Name)
	}

	secrets["app.ini"] = buf.String()
	out.StringData = secrets

	return out, nil
}

func (s *GiteaSecretSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaSecretUpdated
	} else {
		return GiteaSecretFailed
	}
}
