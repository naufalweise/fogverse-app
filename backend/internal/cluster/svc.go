package cluster

import "context"

// Svc service layer of cluster package
type Svc interface {
	CreateTopic(ctx context.Context, topic Topic) error
	UpdateTopic(ctx context.Context, topic Topic) error // Replace topic configuration with new configuration
	DeleteTopic(ctx context.Context, topicName string) error
	GetTopic(ctx context.Context, topicName string) (*Topic, error)
	ListTopics(ctx context.Context) ([]Topic, error)
}
