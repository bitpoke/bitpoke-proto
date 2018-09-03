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
	// GiteaSecretFailed is the event reason for a failed Gitea Secret reconcile
	GiteaSecretFailed EventReason = "GiteaSecretFailed"
	// GiteaSecretUpdated is the event reason for a successful Gitea Secret reconcile
	GiteaSecretUpdated EventReason = "GiteaSecretUpdated"
)

type giteaSecretSyncer struct {
	scheme   *runtime.Scheme
	proj     *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *corev1.Secret
}

// NewGiteaSecretSyncer returns a new sync.Interface for reconciling Gitea PVC
func NewGiteaSecretSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &giteaSecretSyncer{
		scheme:   r,
		existing: &corev1.Secret{},
		proj:     p,
		key:      p.GetGiteaSecretKey(),
	}
}

// GetKey returns the giteaSecretSyncer key through which an existing object may be identified
func (s *giteaSecretSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *giteaSecretSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Gitea Secret
func (s *giteaSecretSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*corev1.Secret)
	out.Labels = GetGiteaPodLabels(s.proj)

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
	cfg, err := createGiteaConfig(s.proj, secrets)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := cfg.WriteTo(&buf); err != nil {
		log.Error(err, "unable to load existing Gitea settings", "project", s.proj.Name)
	}

	secrets["app.ini"] = buf.String()
	out.StringData = secrets

	return out, nil
}

func (s *giteaSecretSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaSecretUpdated
	}
	return GiteaSecretFailed
}
