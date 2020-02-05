package api

// kubeSecret
// swagger:model kubeSecret
type KubeSecret struct {
	// required: true
	// example: Secret
	Kind string `json:"kind"`
	// required: true
	// example: v1
	APIVersion string         `json:"apiVersion"`
	Metadata   SecretMetadata `json:"metadata"`
	// required: true
	// example: {"admin":"YWRtaW4=", "password":"cGFzc3dvcmQ="}
	Data map[string][]byte `json:"data"`
}

// Secret Metadata
// swagger:model secretMetadata
type SecretMetadata struct {
	// required: true
	// example: myAwesomeSecret
	Name string `json:"name"`
	// required: true
	// example: destinationNamespace
	Namespace string `json:"namespace"`
}

type SecretData struct {
	Admin    []byte `json:"admin"`
	Password []byte `json:"password"`
}

// Secret Request
// swagger:model secretRequest
type SecretRequest struct {
	// required: true
	// example: my-secret
	Name string `json:"name"`
	// required: true
	// example: app-namespace
	Namespace string `json:"namespace"`
	// required: true
	// example: {"admin": "my-serviceaccount", "password": "start123", ".dockerconfigjson": "eyJhdXRocyI6eyJodHRwczovL2luZGV4LmRvY2tlci5...", "certificate": "InBhc3N3b3JkIjoiSWJjdDk0MjAiLCJlbWFpbCI6InNlYmFz..."}
	Data map[string]string `json:"data"`
}