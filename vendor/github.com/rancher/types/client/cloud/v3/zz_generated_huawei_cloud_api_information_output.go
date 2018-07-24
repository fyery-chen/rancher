package client

const (
	HuaweiCloudApiInformationOutputType               = "huaweiCloudApiInformationOutput"
	HuaweiCloudApiInformationOutputFieldAvailableZone = "availabilityZoneInfo"
	HuaweiCloudApiInformationOutputFieldHighwaySubnet = "networks"
	HuaweiCloudApiInformationOutputFieldMessage       = "message"
	HuaweiCloudApiInformationOutputFieldNodeFlavor    = "flavors"
	HuaweiCloudApiInformationOutputFieldSshKeyName    = "keypairs"
	HuaweiCloudApiInformationOutputFieldSubnetInfo    = "subnets"
	HuaweiCloudApiInformationOutputFieldVpcInfo       = "vpcs"
)

type HuaweiCloudApiInformationOutput struct {
	AvailableZone []AvailabilityZone `json:"availabilityZoneInfo,omitempty" yaml:"availabilityZoneInfo,omitempty"`
	HighwaySubnet []HighwaySubnet    `json:"networks,omitempty" yaml:"networks,omitempty"`
	Message       string             `json:"message,omitempty" yaml:"message,omitempty"`
	NodeFlavor    []NodeFlavor       `json:"flavors,omitempty" yaml:"flavors,omitempty"`
	SshKeyName    []SshKey           `json:"keypairs,omitempty" yaml:"keypairs,omitempty"`
	SubnetInfo    []SubnetInfo       `json:"subnets,omitempty" yaml:"subnets,omitempty"`
	VpcInfo       []VpcInfo          `json:"vpcs,omitempty" yaml:"vpcs,omitempty"`
}
