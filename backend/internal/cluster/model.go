package cluster

import (
	"encoding/json"
	"fmt"

	"log"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/constant"
	"gitlab.cs.ui.ac.id/fogverse/app/pkg/brominmax"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Topic struct {
	Name             string            `json:"name" binding:"required"`
	Partitions       int               `json:"partitions"`
	Replicas         int               `json:"replicas"`
	UsePartitionAlgo bool              `json:"usePartitionAlgo" default:"false"`
	PartitionConfig  *brominmax.Config `json:"partitionConfig,omitempty"`
}

func (t *Topic) validate() error {
	if t.UsePartitionAlgo {
		if t.PartitionConfig == nil {
			return fmt.Errorf("partitionConfig is required")
		}
	} else {
		if t.Partitions <= 0 {
			return fmt.Errorf("invalid partitions: %d", t.Partitions)
		}
		if t.Replicas <= 0 {
			return fmt.Errorf("invalid replicas: %d", t.Replicas)
		}
	}
	return nil
}

func (t *Topic) Init() error {
	if err := t.validate(); err != nil {
		return err
	}
	if t.UsePartitionAlgo {
		err := t.calculatePartitions()
		if err != nil {
			return err
		}
	}
	return nil
}

// convert topic to kubernetes resource configuration
func (t *Topic) getKubeConfig() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kafka.strimzi.io/v1beta2",
			"kind":       "KafkaTopic",
			"metadata": map[string]interface{}{
				"name":      t.Name,
				"namespace": constant.Namespace,
				"labels": map[string]interface{}{
					"strimzi.io/cluster": "my-cluster",
				},
			},
			"spec": map[string]interface{}{
				"partitions": t.Partitions,
				"replicas":   t.Replicas,
			},
		},
	}
	if t.UsePartitionAlgo {
		partitionConfigJson, err := json.Marshal(t.PartitionConfig)
		if err != nil {
			log.Printf("failed converting kubeconfig to topic '%s': %v", t.Name, err)
			return obj
		}
		obj.SetAnnotations(map[string]string{
			"fogverse-app/partition-config": string(partitionConfigJson),
		})
	}
	return obj
}

// NewTopicFromKubeConfig converts kafka topic kube config to Topic
func NewTopicFromKubeConfig(obj *unstructured.Unstructured) (*Topic, error) {
	partitions, found, err := unstructured.NestedInt64(obj.Object, "spec", "partitions")
	if err != nil || !found {
		return nil, fmt.Errorf("failed to extract 'spec.partitions' in KafkaTopic: %v", err)
	}

	// The path to replicas is "spec", "replicas"
	replicas, found, err := unstructured.NestedInt64(obj.Object, "spec", "replicas")
	if err != nil || !found {
		replicas = 1 // default replica is 1 if not specified
	}

	// get partition config if exist
	var partitionConfig *brominmax.Config
	partitionConfigStr := obj.GetAnnotations()["fogverse-app/partition-config"]
	if partitionConfigStr != "" {
		if err := json.Unmarshal([]byte(partitionConfigStr), &partitionConfig); err != nil {
			log.Printf("failed to unmarshal of json partition config in topic '%s': %v", obj.GetName(), err)
		}
	}

	return &Topic{
		Name:             obj.GetName(),
		Partitions:       int(partitions),
		Replicas:         int(replicas),
		UsePartitionAlgo: partitionConfig != nil,
		PartitionConfig:  partitionConfig}, nil
}

// calculatePartitions fills partitions and replicas field using partition algo and partition config field
func (t *Topic) calculatePartitions() error {
	if t.PartitionConfig == nil {
		return fmt.Errorf("partitionConfig is required")
	}
	partitions, replicas, err := brominmax.GetPartitions(*t.PartitionConfig)
	if err != nil {
		return err
	}
	t.Partitions = partitions
	t.Replicas = replicas
	return nil
}
