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

	"github.com/presslabs/controller-util/syncer"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

// NewGiteaDeploymentSyncer returns a new syncer.Interface for reconciling Gitea Deployment
func NewGiteaDeploymentSyncer(proj *dashboardv1alpha1.Project) syncer.Interface {
	obj := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      proj.GetGiteaDeploymentName(),
			Namespace: proj.GetNamespaceName(),
		},
	}

	return syncer.New("GiteaDeployment", proj, obj, func(existing runtime.Object) error {
		out := existing.(*appsv1.Deployment)

		out.Labels = GetGiteaPodLabels(proj)

		out.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: GetGiteaLabels(proj),
		}

		out.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: GetGiteaPodLabels(proj),
			},
			Spec: corev1.PodSpec{
				Volumes: []corev1.Volume{
					{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: proj.GetGiteaSecretName(),
							},
						},
					},
					{
						Name: "data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: proj.GetGiteaPVCName(),
							},
						},
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
								"memory": resource.MustParse(giteaRequestsMemory),
								"cpu":    resource.MustParse(giteaRequestsCPU),
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "data",
								MountPath: "/path",
							},
							{
								Name:      "config",
								MountPath: "/data/gitea/conf/app.ini",
								SubPath:   "app.ini",
							},
						},
					},
				},
			},
		}
		return nil
	})
}
