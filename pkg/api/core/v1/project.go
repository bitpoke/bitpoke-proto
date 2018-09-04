package v1

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	dashboardv1alpha1 "github.com/presslabs/dashboard/pkg/apis/dashboard/v1alpha1"
)

type projectsServer struct {
	client client.Client
}

func (s *projectsServer) List(r *ListRequest, stream Projects_ListServer) error {
	projects := &dashboardv1alpha1.ProjectList{}

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

func newFromK8s(p *dashboardv1alpha1.Project) *Project {
	return &Project{Id: p.Name, Name: p.Name}
}
