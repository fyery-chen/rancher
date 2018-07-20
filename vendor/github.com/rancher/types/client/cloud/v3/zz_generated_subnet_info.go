package client

const (
	SubnetInfoType            = "subnetInfo"
	SubnetInfoFieldSubnetId   = "subnetId"
	SubnetInfoFieldSubnetName = "subnetName"
)

type SubnetInfo struct {
	SubnetId   string `json:"subnetId,omitempty" yaml:"subnetId,omitempty"`
	SubnetName string `json:"subnetName,omitempty" yaml:"subnetName,omitempty"`
}
