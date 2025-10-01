package deployer

import (
	"context"
	"fmt"
	"log"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/deployer/constant"
	"gitlab.cs.ui.ac.id/fogverse/app/pkg/kafkautils"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Deployer struct {
	KubeInterface kubernetes.Interface
	DeployRequest DeployRequest
}

//func (d *Deployer) New(request DeployRequest) *Deployer {
//	clientset, err := kubeutils.GetKubeClient()
//	if err != nil {
//		panic(err)
//	}
//	return &Deployer{
//		KubeInterface: clientset,
//		DeployRequest: request,
//	}
//}

func (d *Deployer) Deploy() error {
	var err error
	switch d.DeployRequest.Type {

	case "job":
		err = d.createJob()
	default:
		err = d.createDeployment()
	}
	return err
}

func (d *Deployer) createDeployment() error {
	deploymentsClient := d.KubeInterface.AppsV1().Deployments(constant.Namespace)

	bootstrapServers := kafkautils.GetBootstrapServersInternal()

	var replicas int32 = 1
	if d.DeployRequest.Replicas != nil {
		replicas = *d.DeployRequest.Replicas
	}
	// Define the Deployment object.
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.DeployRequest.Name,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": constant.KubeLabel,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: d.DeployRequest.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": d.DeployRequest.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": d.DeployRequest.Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  d.DeployRequest.Name,
							Image: d.DeployRequest.ImageURL,
							Ports: []apiv1.ContainerPort{
								{
									// Assuming the deployed container exposes port 80.
									// This can be customized if needed.
									ContainerPort: 80,
								},
							},
							Env: []apiv1.EnvVar{
								{Name: "INPUT_TOPIC", Value: d.DeployRequest.InputTopic},
								{Name: "OUTPUT_TOPIC", Value: d.DeployRequest.OutputTopic},
								{Name: "REPLICAS", Value: fmt.Sprint(replicas)},
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
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// CreateJob creates a new job in the "default" namespace.
func (d *Deployer) createJob() error {
	jobsClient := d.KubeInterface.BatchV1().Jobs(apiv1.NamespaceDefault)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.DeployRequest.Name,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  d.DeployRequest.Name,
							Image: d.DeployRequest.ImageURL,
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

// CreateTopic creates kafka topic via kubernetes api via strimzi
//func CreateTopic(clientset *kubernetes.Clientset, name, imageURL string) error {
//	clientset.
//}
