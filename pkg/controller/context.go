package controller

import (
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	dashclientset "github.com/presslabs/dashboard/pkg/client/clientset/versioned"
	dashinformers "github.com/presslabs/dashboard/pkg/client/informers/externalversions"
)

// Context contains various types that are used by controller implementations.
// We purposely don't have specific informers/listers here, and instead keep
// a reference to a SharedInformerFactory so that controllers can choose
// themselves which listers are required.
type Context struct {
	// KubeClient is a Kubernetes clientset
	KubeClient kubernetes.Interface
	// Recorder to record events to
	Recorder record.EventRecorder
	// KubeSharedInformerFactory can be used to obtain shared
	// SharedIndexInformer instances for Kubernetes types
	KubeSharedInformerFactory kubeinformers.SharedInformerFactory

	// DashboardClient is a Presslabs Dashboard clientset
	DashboardClient dashclientset.Interface

	// DashboardSharedInformerFactory can be used to obtain shared
	// SharedIndexInformer instances for Presslabs Dashboard types
	DashboardSharedInformerFactory dashinformers.SharedInformerFactory
}
