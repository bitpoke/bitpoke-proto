/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"bytes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/internal/project"

	"github.com/presslabs/controller-util/rand"
)

// NewGiteaSecretSyncer returns a new syncer.Interface for reconciling Gitea Secret
func NewGiteaSecretSyncer(proj *project.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := proj.ComponentLabels(project.GiteaSecret)

	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(project.GiteaSecret),
			Namespace: proj.ComponentName(project.Namespace),
		},
	}

	return syncer.NewObjectSyncer("GiteaSecret", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Secret)
		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)

		if len(out.Data) == 0 {
			out.Data = make(map[string][]byte)
		}

		if len(out.Data["SECRET_KEY"]) == 0 {
			r, err := rand.AlphaNumericString(20)
			if err != nil {
				return err
			}
			out.Data["SECRET_KEY"] = []byte(r)
		}

		if len(out.Data["INTERNAL_TOKEN"]) == 0 {
			r, err := rand.AlphaNumericString(20)
			if err != nil {
				return err
			}
			out.Data["INTERNAL_TOKEN"] = []byte(r)
		}

		cfg, err := createGiteaConfig(proj, out.Data)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if _, err := cfg.WriteTo(&buf); err != nil {
			log.Error(err, "unable to load existing Gitea settings", "project", proj.Name)
		}
		out.Data["app.ini"] = buf.Bytes()

		return nil
	})
}
