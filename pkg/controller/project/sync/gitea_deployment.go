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
	// GiteaDeploymentFailed is the event reason for a failed Gitea Deployment reconcile
	GiteaDeploymentFailed EventReason = "GiteaDeploymentFailed"
	// GiteaDeploymentUpdated is the event reason for a successful Gitea Deployment reconcile
	GiteaDeploymentUpdated EventReason = "GiteaDeploymentUpdated"
)

type giteaDeploymentSyncer struct {
	scheme   *runtime.Scheme
	proj     *dashboardv1alpha1.Project
	key      types.NamespacedName
	existing *appsv1.Deployment
}

// NewGiteaDeploymentSyncer returns a new sync.Interface for reconciling Gitea Deployment
func NewGiteaDeploymentSyncer(p *dashboardv1alpha1.Project, r *runtime.Scheme) Interface {
	return &giteaDeploymentSyncer{
		scheme:   r,
		existing: &appsv1.Deployment{},
		proj:     p,
		key:      p.GetGiteaDeploymentKey(),
	}
}

// GetKey returns the giteaDeploymentSyncer key through which an existing object may be identified
func (s *giteaDeploymentSyncer) GetKey() types.NamespacedName { return s.key }

// GetExistingObjectPlaceholder returns a Placeholder object if an existing one is not found
func (s *giteaDeploymentSyncer) GetExistingObjectPlaceholder() runtime.Object { return s.existing }

// T is the transform function used to reconcile the Gitea Deployment
func (s *giteaDeploymentSyncer) T(in runtime.Object) (runtime.Object, error) {
	out := in.(*appsv1.Deployment)

	out.Labels = GetGiteaPodLabels(s.proj)

	out.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: GetGiteaLabels(s.proj),
	}

	out.Spec.Template = corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: GetGiteaPodLabels(s.proj),
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: s.proj.GetGiteaSecretName(),
						},
					},
				},
				{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: s.proj.GetGiteaPVCName(),
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

	return out, nil
}

func (s *giteaDeploymentSyncer) GetErrorEventReason(err error) EventReason {
	if err == nil {
		return GiteaDeploymentUpdated
	}
	return GiteaDeploymentFailed
}
