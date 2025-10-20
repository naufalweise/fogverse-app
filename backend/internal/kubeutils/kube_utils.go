// kubernetes/client.go
// This package contains all the logic for interacting with the Kubernetes API.

package kubeutils

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// GetKubeClient provides a Kubernetes clientset.
// It tries to use in-cluster configuration first, which is suitable for when the
// application is running inside a pod in the Kubernetes cluster.
// If that fails, it falls back to using the local kubeconfig file, which is
// useful for local development and testing.
func GetKubeClient(config *rest.Config) (*kubernetes.Clientset, error) {
	// Create the clientset from the config.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func GetKubeConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Could not create in-cluster config: %v. Falling back to kubeconfig.", err)
		// Fallback to using the user's kubeconfig file.
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			// Handle cases where home directory is not available.
			kubeconfig = os.Getenv("KUBECONFIG")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	return config, nil
}
