// kubernetes/client.go
// This package contains all the logic for interacting with the Kubernetes API.

package kubernetes

import (
	"context"
	"log"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func GetKubeClient() (*kubernetes.Clientset, error) {
	// Try to create an in-cluster config.
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

		// Use the current context in kubeconfig.
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Create the clientset from the config.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func CreateDeployment(clientset *kubernetes.Clientset, deploymentName, imageURL string, replicas int32) error {
	// Get the deployments client for the 'default' namespace.
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	// Define the Deployment object.
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  deploymentName,
							Image: imageURL,
							Ports: []apiv1.ContainerPort{
								{
									// Assuming the deployed container exposes port 80.
									// This can be customized if needed.
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create the Deployment.
	log.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// CreateJob creates a new job in the "default" namespace.
func CreateJob(clientset *kubernetes.Clientset, name, imageURL string) error {
	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: imageURL,
						},
					},
					// For a job, we want it to stop after completion, not restart indefinitely.
					// 'Never' ensures it doesn't restart. 'OnFailure' is another option.
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}

	log.Println("Creating job...")
	result, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("Created job %q.\n", result.GetObjectMeta().GetName())
	return nil
}
