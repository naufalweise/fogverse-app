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
	namespace := "kafka"
	service, err := k.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	fmt.Print(service.Spec.Ports)
	serviceIP := service.Spec.ClusterIP
	servicePort := service.Spec.Ports[0].NodePort
	return fmt.Sprintf("%s:%d", serviceIP, servicePort), nil
}
