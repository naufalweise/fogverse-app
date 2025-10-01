package service

import (
	"context"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/deployer/constant"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	Name        string `json:"name" binding:"required"`
	ImageURL    string `json:"imageURL" binding:"required"`
	Replicas    int32  `json:"replicas"`
	InputTopic  string `json:"inputTopic" binding:"required"`
	OutputTopic string `json:"outputTopic" binding:"required"`
	State
}

type State struct {
	ReadyReplicas int32 `json:"readyReplicas"`
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
	return Service{Name: deployment.Name, ImageURL: container.Image, Replicas: deployment.Status.Replicas, InputTopic: inputTopic, OutputTopic: outputTopic, State: State{ReadyReplicas: deployment.Status.ReadyReplicas}}
}
func List(k kubernetes.Interface) ([]Service, error) {
	deploymentsClient := k.AppsV1().Deployments("kafka")
	deploymentList, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/managed-by=" + constant.KubeLabel})
	if err != nil {
		return nil, err
	}
	services := make([]Service, 0)
	for _, deployment := range deploymentList.Items {
		services = append(services, deploymentToService(&deployment))
	}
	return services, nil
}
