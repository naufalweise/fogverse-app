package cluster

import (
	"context"
	"fmt"
	"log"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/constant"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type SvcImpl struct {
	dynamicClient dynamic.Interface
	topicClient   dynamic.ResourceInterface
}

var topicGvr = schema.GroupVersionResource{
	Group:    "kafka.strimzi.io",
	Version:  "v1beta2",
	Resource: "kafkatopics",
}

func NewSvcImpl(dynamicClient dynamic.Interface) *SvcImpl {
	topicClient := dynamicClient.Resource(topicGvr).Namespace(constant.Namespace)
	return &SvcImpl{dynamicClient: dynamicClient, topicClient: topicClient}

}

func (s *SvcImpl) CreateTopic(ctx context.Context, t Topic) error {
	err := t.Init()
	if err != nil {
		return err
	}
	_, err = s.topicClient.Create(ctx, t.getKubeConfig(), metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Printf("Created topic: %s", t.getKubeConfig().GetName())
	return nil
}

func (s *SvcImpl) DeleteTopic(ctx context.Context, topicName string) error {
	err := s.topicClient.Delete(ctx, topicName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("Deleted topic: %s", topicName)
	return nil
}

func (s *SvcImpl) UpdateTopic(ctx context.Context, t Topic) error {
	err := t.Init()
	if err != nil {
		return err
	}
	_, err = s.topicClient.Apply(ctx, t.Name, t.getKubeConfig(), metav1.ApplyOptions{})
	if err != nil {
		return err
	}
	log.Printf("Updated topic: %s", t.getKubeConfig().GetName())
	return nil
}

func (s *SvcImpl) GetTopic(ctx context.Context, topicName string) (*Topic, error) {
	res, err := s.topicClient.Get(ctx, topicName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	topic, err := NewTopicFromKubeConfig(res)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (s *SvcImpl) ListTopics(ctx context.Context) ([]Topic, error) {
	res, err := s.topicClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("no topics found")
	}
	topicList := make([]Topic, 0, len(res.Items))
	for _, item := range res.Items {
		topic, err := NewTopicFromKubeConfig(&item)
		if err != nil {
			log.Println(err)
			continue
		}
		topicList = append(topicList, *topic)
	}

	return topicList, nil
}
