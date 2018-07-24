package client

const (
	VpcInfoType         = "vpcInfo"
	VpcInfoFieldVpcName = "name"
)

type VpcInfo struct {
	VpcName string `json:"name,omitempty" yaml:"name,omitempty"`
}
