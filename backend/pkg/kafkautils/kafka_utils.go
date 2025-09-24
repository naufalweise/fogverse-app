package kafkautils

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetBootstrapServers gets the external IP of kafka cluster bootstrap servers via nodeport
func GetBootstrapServers(k kubernetes.Interface) (string, error) {
	serviceName := "my-cluster-kafka-nodeport-bootstrap"
	//nodeName := "lima-rancher-desktop"
	namespace := "kafka"
	// get service port
	service, err := k.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	servicePort := service.Spec.Ports[0].NodePort

	// get node ip
	//node, err := k.CoreV1().Nodes().Get(context.TODO())
	serviceIP := ""
	return fmt.Sprintf("%s:%d", serviceIP, servicePort), nil
}

func GetBootstrapServersInternal() string {
	return "my-cluster-kafka-bootstrap:9092"
}
