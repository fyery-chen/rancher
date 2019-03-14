package client

const (
	CloudProviderImageListInputType              = "cloudProviderImageListInput"
	CloudProviderImageListInputFieldProviderType = "providerType"
)

type CloudProviderImageListInput struct {
	ProviderType string `json:"providerType,omitempty" yaml:"providerType,omitempty"`
}
