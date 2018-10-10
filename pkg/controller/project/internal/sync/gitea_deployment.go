/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package sync

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/imdario/mergo"

	"github.com/presslabs/controller-util/mergo/transformers"
	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

var (
	resGiteaRequestsMemory = resource.MustParse(giteaRequestsMemory)
	resGiteaRequestsCPU    = resource.MustParse(giteaRequestsCPU)
)

// NewGiteaDeploymentSyncer returns a new syncer.Interface for reconciling Gitea Deployment
func NewGiteaDeploymentSyncer(proj *dashboardv1alpha1.Project, cl client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      giteaDeploymentName(proj),
			Namespace: getNamespaceName(proj),
		},
	}

	return syncer.NewObjectSyncer("GiteaDeployment", proj, obj, cl, scheme, func(existing runtime.Object) error {
		out := existing.(*appsv1.Deployment)

		out.Labels = giteaPodLabels(proj)

		out.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: giteaLabels(proj),
		}

		out.Spec.Template.ObjectMeta = metav1.ObjectMeta{
			Labels: giteaPodLabels(proj),
		}

		spec := corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "config-secret",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: giteaSecretName(proj),
						},
					},
				},
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: giteaPVCName(proj),
						},
					},
				},
			},
			InitContainers: []corev1.Container{
				{
					Name:  "gitea-config-init",
					Image: "busybox",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config-secret",
							MountPath: "/secret",
						},
						{
							Name:      "config",
							MountPath: "/conf/",
						},
					},
					Args: []string{"cp", "/secret/app.ini", "/conf/app.ini"},
				},
			},
			Containers: []corev1.Container{
				{
					Name:            giteaName,
					Image:           giteaImage,
					ImagePullPolicy: "IfNotPresent",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: int32(giteaHTTPPort),
						},
						{
							Name:          "ssh",
							ContainerPort: giteaSSHPort,
						},
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/",
								Port: intstr.FromInt(giteaHTTPPort),
							},
						},
						InitialDelaySeconds: 30,
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/",
								Port: intstr.FromInt(giteaHTTPPort),
							},
						},
						InitialDelaySeconds: 30,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"memory": resGiteaRequestsMemory,
							"cpu":    resGiteaRequestsCPU,
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "data",
							MountPath: "/path",
						},
						{
							Name:      "config",
							MountPath: "/data/gitea/conf/",
						},
					},
				},
			},
		}

		err := mergo.Merge(&out.Spec.Template.Spec, spec, mergo.WithTransformers(transformers.PodSpec))
		if err != nil {
			return err
		}

		return nil
	})
}
