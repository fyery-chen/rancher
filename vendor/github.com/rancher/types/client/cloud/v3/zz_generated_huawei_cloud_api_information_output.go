package client

const (
	HuaweiCloudApiInformationOutputType               = "huaweiCloudApiInformationOutput"
	HuaweiCloudApiInformationOutputFieldAvailableZone = "availableZone"
	HuaweiCloudApiInformationOutputFieldMessage       = "message"
	HuaweiCloudApiInformationOutputFieldNodeFlavor    = "nodeFlavor"
	HuaweiCloudApiInformationOutputFieldSshKeyName    = "sshKeyName"
	HuaweiCloudApiInformationOutputFieldVpcInfo       = "vpcInfo"
)

type HuaweiCloudApiInformationOutput struct {
	AvailableZone []AvailableZone `json:"availableZone,omitempty" yaml:"availableZone,omitempty"`
	Message       string          `json:"message,omitempty" yaml:"message,omitempty"`
	NodeFlavor    []NodeFlavor    `json:"nodeFlavor,omitempty" yaml:"nodeFlavor,omitempty"`
	SshKeyName    []string        `json:"sshKeyName,omitempty" yaml:"sshKeyName,omitempty"`
	VpcInfo       []VpcInfo       `json:"vpcInfo,omitempty" yaml:"vpcInfo,omitempty"`
}
