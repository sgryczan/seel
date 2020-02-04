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
