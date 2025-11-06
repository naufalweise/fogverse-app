package service

import (
	"context"
	"fmt"
	"log"

	kedav2 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	kedaclientset "github.com/kedacore/keda/v2/pkg/generated/clientset/versioned"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/constant"
	"gitlab.cs.ui.ac.id/fogverse/app/pkg/kafkautils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Svc is a service layer of 'service' module
type Svc interface {
	Create(ctx context.Context, service Service) error
	List(ctx context.Context) ([]Service, error)
	CreateAutoscale(ctx context.Context, service Service) error
}

type SvcImpl struct {
	kubeClient kubernetes.Interface
	kedaClient kedaclientset.Interface
}

func NewSvc(kubeClient kubernetes.Interface, kedaClient kedaclientset.Interface) Svc {
	return &SvcImpl{kubeClient: kubeClient, kedaClient: kedaClient}
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

	if service.UseKeda {
		log.Println("Creating Keda service...")
		err := s.CreateAutoscale(ctx, service)
		if err != nil {
			return err
		}
	}

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

// CreateAutoscale create KEDA scaled object to autoscale service
func (s *SvcImpl) CreateAutoscale(ctx context.Context, service Service) error {

	// init default keda config
	config := KedaConfig{}
	if service.KedaConfig != nil {
		config = *service.KedaConfig
	}
	if config.LagThreshold == nil {
		defaultLagThreshold := int32(10)
		config.LagThreshold = &defaultLagThreshold
	}

	// create scaled object
	scaledObject := &kedav2.ScaledObject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("fogverse-run-%s", service.Name),
			Namespace: constant.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": constant.KubeLabel,
			},
		},
		Spec: kedav2.ScaledObjectSpec{
			ScaleTargetRef: &kedav2.ScaleTarget{
				Name: service.Name,
			},
			MinReplicaCount: config.MinReplicaCount,
			MaxReplicaCount: config.MaxReplicaCount,
			CooldownPeriod:  config.CooldownPeriod,
			PollingInterval: config.PollingInterval,
			Triggers: []kedav2.ScaleTriggers{
				{
					Type: "apache-kafka",
					Metadata: map[string]string{
						"bootstrapServers": kafkautils.GetBootstrapServersInternal(),
						"consumerGroup":    service.Name,
						"topic":            service.InputTopic,
						"lagThreshold":     fmt.Sprint(*config.LagThreshold),
					},
				},
			},
		},
	}
	result, err := s.kedaClient.KedaV1alpha1().ScaledObjects(constant.Namespace).Create(ctx, scaledObject, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating ScaledObject: %v", err)
		return err
	}
	log.Printf("Successfully created ScaledObject: %s\n", result.GetName())
	return nil
}
