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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

const (
	GiteaFailed  EventReason = "GiteaFailed"
	GiteaUpdated EventReason = "GiteaUpdated"
)

type GiteaDeploymentSyncer struct {
	scheme   *runtime.Scheme
	p        *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *appsv1.Deployment
}

func NewGiteaDeploymentSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) *GiteaDeploymentSyncer {
	return &GiteaDeploymentSyncer{
		scheme:   r,
		existing: &appsv1.Deployment{},
		p:        p,
		key:      p.GetGiteaDeploymentKey(),
	}
}

func (s *GiteaDeploymentSyncer) GetKey() types.NamespacedName                 { return s.key }
func (s *GiteaDeploymentSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

func (s *GiteaDeploymentSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*appsv1.Deployment)

	out.Labels = GetGiteaPodLabels(s.p)

	replicas := int32(giteaReplicas)
	maxUnavailable := intstr.FromInt(giteaMaxUnavailable)
	maxSurge := intstr.FromInt(giteaMaxSurge)

	out.Spec.Replicas = &replicas
	out.Spec.Strategy = appsv1.DeploymentStrategy{
		Type: "RollingUpdate",
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &maxUnavailable,
			MaxSurge:       &maxSurge,
		},
	}
	out.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: GetGiteaLabels(s.p),
	}

	out.Spec.Template = corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: GetGiteaPodLabels(s.p),
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: s.p.GetGiteaSecretName(),
						},
					},
				},
				{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: s.p.GetGiteaPVCName(),
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
							ContainerPort: int32(giteaHTTPInternalPort),
						},
						{
							Name:          "ssh",
							ContainerPort: 22,
						},
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/",
								Port: intstr.FromInt(giteaHTTPInternalPort),
							},
						},
						InitialDelaySeconds: 30,
					},
					ReadinessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/",
								Port: intstr.FromInt(giteaHTTPInternalPort),
							},
						},
						InitialDelaySeconds: 30,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"memory": resource.MustParse("512Mi"),
							"cpu":    resource.MustParse("100m"),
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

	return out, nil
}

func (s *GiteaDeploymentSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaUpdated
	} else {
		return GiteaFailed
	}
}
