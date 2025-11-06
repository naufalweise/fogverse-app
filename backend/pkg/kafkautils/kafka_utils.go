package kafkautils

import (
	"context"
	"errors"
	"fmt"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/constant"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetBootstrapServers gets the external IP of kafka cluster bootstrap servers via nodeport
func GetBootstrapServers(ctx context.Context, k kubernetes.Interface) (string, error) {
	serviceName := "my-cluster-kafka-nodeport-bootstrap"
	//nodeName := "lima-rancher-desktop"
	namespace := constant.Namespace
	// get service port
	service, err := k.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	servicePort := service.Spec.Ports[0].NodePort

	// get node ip
	nodes, err := k.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if len(nodes.Items) < 1 {
		return "", errors.New("no kubernetes nodes found")
	}
	// use the first node
	node := nodes.Items[0]

	// get external ip
	externalIp := ""
	for _, address := range node.Status.Addresses {
		if address.Type == v1.NodeExternalIP {
			externalIp = address.Address
			break
		}
	}
	return fmt.Sprintf("%s:%d", externalIp, servicePort), nil
}

func GetBootstrapServersInternal() string {
	return fmt.Sprintf("my-cluster-kafka-bootstrap.%s:9092", constant.Namespace)
}
