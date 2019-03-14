package client

const (
	CloudProviderImageListOutputType         = "cloudProviderImageListOutput"
	CloudProviderImageListOutputFieldImages  = "images"
	CloudProviderImageListOutputFieldMessage = "message"
)

type CloudProviderImageListOutput struct {
	Images  []ImageInfo `json:"images,omitempty" yaml:"images,omitempty"`
	Message string      `json:"message,omitempty" yaml:"message,omitempty"`
}
