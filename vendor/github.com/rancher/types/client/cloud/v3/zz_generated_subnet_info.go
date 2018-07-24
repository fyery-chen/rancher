package client

const (
	SubnetInfoType            = "subnetInfo"
	SubnetInfoFieldSubnetName = "name"
	SubnetInfoFieldVpcId      = "vpc_id"
)

type SubnetInfo struct {
	SubnetName string `json:"name,omitempty" yaml:"name,omitempty"`
	VpcId      string `json:"vpc_id,omitempty" yaml:"vpc_id,omitempty"`
}
