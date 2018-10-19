package v1

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type projectsServer struct {
	client client.Client
}

func (s *projectsServer) List(r *ListRequest, stream Projects_ListServer) error {
	projects := &corev1.NamespaceList{}

	err := s.client.List(context.TODO(), &client.ListOptions{}, projects)
	if err != nil {
		return err
	}

	for _, project := range projects.Items {
		err = stream.Send(newFromK8s(&project))
		if err != nil {
			return err
		}
	}

	return nil
}

// NewProjectServer creates a new gRPC server for projects
func NewProjectServer(client client.Client) ProjectsServer {
	return &projectsServer{
		client: client,
	}
}

func newFromK8s(p *corev1.Namespace) *Project {
	return &Project{Id: p.Name, Name: p.Name}
}
