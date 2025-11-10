package service

// Service is a type of fogverse run containers
// workload that keeps running.
type Service struct {
	Name        string `json:"name" binding:"required"`
	ImageURL    string `json:"imageURL" binding:"required"`
	Replicas    *int32 `json:"replicas"`
	InputTopic  string `json:"inputTopic" binding:"required"`
	OutputTopic string `json:"outputTopic" binding:"required"`
	State
	KedaConfig *KedaConfig `json:"kedaConfig"`
	UseKeda    bool        `json:"useKeda" default:"false"`
}

type State struct {
	ReadyReplicas int32 `json:"readyReplicas"`
}

type KedaConfig struct {
	MinReplicaCount *int32 `json:"minReplicaCount"`
	MaxReplicaCount *int32 `json:"maxReplicaCount"`
	CooldownPeriod  *int32 `json:"cooldownPeriod"`
	PollingInterval *int32 `json:"pollingInterval"`
	LagThreshold    *int32 `json:"lagThreshold"`
}
