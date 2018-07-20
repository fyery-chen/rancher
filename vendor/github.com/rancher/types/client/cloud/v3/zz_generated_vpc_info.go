package client

const (
	VpcInfoType            = "vpcInfo"
	VpcInfoFieldSubnetInfo = "subnetInfo"
	VpcInfoFieldVpcId      = "vpcId"
	VpcInfoFieldVpcName    = "vpcName"
)

type VpcInfo struct {
	SubnetInfo []SubnetInfo `json:"subnetInfo,omitempty" yaml:"subnetInfo,omitempty"`
	VpcId      string       `json:"vpcId,omitempty" yaml:"vpcId,omitempty"`
	VpcName    string       `json:"vpcName,omitempty" yaml:"vpcName,omitempty"`
}
