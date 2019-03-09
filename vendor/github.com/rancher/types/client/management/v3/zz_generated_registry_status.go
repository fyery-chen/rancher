package client

const (
	RegistryStatusType               = "registryStatus"
	RegistryStatusFieldRegistryState = "RegistryState"
)

type RegistryStatus struct {
	RegistryState string `json:"RegistryState,omitempty" yaml:"RegistryState,omitempty"`
}
