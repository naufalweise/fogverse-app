package job

import (
	"context"
	"log"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/constant"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Job struct {
	Name        string `json:"name" binding:"required"`
	ImageURL    string `json:"imageURL" binding:"required"`
	InputTopic  string `json:"inputTopic" binding:"required"`
	OutputTopic string `json:"outputTopic" binding:"required"`
	Completion  int32  `json:"completion" binding:"required"`  // completion parameter in kubernetes job
	Parallelism int32  `json:"parallelism" binding:"required"` // parallelism parameter in kubernetes job
}

type Svc interface {
	List(ctx context.Context) ([]Job, error)
	Create(ctx context.Context, job Job) error
}

type SvcImpl struct {
	kubeClient kubernetes.Interface
}

func NewSvc(kubeClient kubernetes.Interface) *SvcImpl {
	return &SvcImpl{kubeClient: kubeClient}
}

func (s *SvcImpl) Create(ctx context.Context, job Job) error {

	jobsClient := s.kubeClient.BatchV1().Jobs(constant.Namespace)

	kubeJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: job.Name,
		},
		Spec: batchv1.JobSpec{
			Completions: &job.Completion,
			Parallelism: &job.Parallelism,
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  job.Name,
							Image: job.ImageURL,
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
	result, err := jobsClient.Create(ctx, kubeJob, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("Created job %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func (s *SvcImpl) List(ctx context.Context) ([]Job, error) {
	return nil, nil
}
