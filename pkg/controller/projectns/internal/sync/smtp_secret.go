/*
Copyright 2019 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/controller-util/syncer"
	"github.com/presslabs/dashboard/pkg/cmd/manager/options"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

func namespace() string {
	if ns := os.Getenv("MY_NAMESPACE"); ns != "" {
		return ns
	}
	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}
	return corev1.NamespaceDefault
}

// NewSMTPSecretSyncer returns a new syncer.Interface for reconciling smtp credentials
func NewSMTPSecretSyncer(proj *projectns.ProjectNamespace, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	objLabels := labels.Set{
		"presslabs.com/kind":                "smtp",
		"dashboard.presslabs.com/reconcile": "true",
	}

	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.ComponentName(projectns.SMTPSecret),
			Namespace: proj.Name,
		},
	}

	return syncer.NewObjectSyncer("SMTPSecret", proj.Unwrap(), obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*corev1.Secret)

		// get default smtp secret
		defaultSecret := &corev1.Secret{}
		key := client.ObjectKey{
			Name:      options.SMTPSecret,
			Namespace: namespace(),
		}
		if err := cl.Get(context.TODO(), key, defaultSecret); err != nil {
			return err
		}

		out.Labels = labels.Merge(labels.Merge(out.Labels, objLabels), controllerLabels)
		out.Data = defaultSecret.Data

		return nil
	})
}
