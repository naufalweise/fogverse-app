package deployer

// DeployRequest defines the structure of the JSON body for the /deploy endpoint.
// It includes validation tags to ensure the fields are not empty.
type DeployRequest struct {
	Name        string `json:"name" binding:"required"`
	ImageURL    string `json:"imageURL" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=service job"`
	Replicas    *int32 `json:"replicas"` // Replicas are optional, only used for service
	InputTopic  string `json:"inputTopic" binding:"required"`
	OutputTopic string `json:"outputTopic" binding:"required"`
}
