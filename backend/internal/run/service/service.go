package service

import (
	"context"
	"fmt"
	"log"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/constant"
	"gitlab.cs.ui.ac.id/fogverse/app/pkg/kafkautils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	Name        string `json:"name" binding:"required"`
	ImageURL    string `json:"imageURL" binding:"required"`
	Replicas    *int32 `json:"replicas"`
	InputTopic  string `json:"inputTopic" binding:"required"`
	OutputTopic string `json:"outputTopic" binding:"required"`
	State
}

type State struct {
	ReadyReplicas int32 `json:"readyReplicas"`
}

// Svc is a service layer of 'service' module
type Svc interface {
	Create(ctx context.Context, service Service) error
	List(ctx context.Context) ([]Service, error)
}

type SvcImpl struct {
	kubeClient kubernetes.Interface
}

func NewSvc(kubeClient kubernetes.Interface) Svc {
	return &SvcImpl{kubeClient: kubeClient}
}

func deploymentToService(deployment *v1.Deployment) Service {
	container := deployment.Spec.Template.Spec.Containers[0]
	inputTopic := ""
	outputTopic := ""
	for _, env := range container.Env {
		if env.Name == "INPUT_TOPIC" {
			inputTopic = env.Value
		} else if env.Name == "OUTPUT_TOPIC" {
			outputTopic = env.Value
		}
	}
	return Service{Name: deployment.Name, ImageURL: container.Image, Replicas: deployment.Spec.Replicas, InputTopic: inputTopic, OutputTopic: outputTopic, State: State{ReadyReplicas: deployment.Status.ReadyReplicas}}
}

func (s *SvcImpl) Create(ctx context.Context, service Service) error {
	deploymentsClient := s.kubeClient.AppsV1().Deployments(constant.Namespace)

	bootstrapServers := kafkautils.GetBootstrapServersInternal()

	// Define the Deployment object.
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: service.Name,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": constant.KubeLabel,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: service.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": service.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": service.Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  service.Name,
							Image: service.ImageURL,
							Ports: []apiv1.ContainerPort{
								{
									// Assuming the deployed container exposes port 80.
									// This can be customized if needed.
									ContainerPort: 80,
								},
							},
							Env: []apiv1.EnvVar{
								{Name: "INPUT_TOPIC", Value: service.InputTopic},
								{Name: "OUTPUT_TOPIC", Value: service.OutputTopic},
								{Name: "REPLICAS", Value: fmt.Sprint(service.Replicas)},
								{Name: "BOOTSTRAP_SERVERS", Value: bootstrapServers},
							},
						},
					},
				},
			},
		},
	}

	// Create the Deployment.
	log.Println("Creating deployment...")
	result, err := deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func (s *SvcImpl) List(ctx context.Context) ([]Service, error) {
	deploymentsClient := s.kubeClient.AppsV1().Deployments("kafka")
	deploymentList, err := deploymentsClient.List(ctx, metav1.ListOptions{LabelSelector: "app.kubernetes.io/managed-by=" + constant.KubeLabel})
	if err != nil {
		return nil, err
	}
	services := make([]Service, 0)
	for _, deployment := range deploymentList.Items {
		services = append(services, deploymentToService(&deployment))
	}
	return services, nil
}
