package client

const (
	NodeFlavorType       = "nodeFlavor"
	NodeFlavorFieldName  = "name"
	NodeFlavorFieldRam   = "ram"
	NodeFlavorFieldVcpus = "vcpus"
)

type NodeFlavor struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Ram   int64  `json:"ram,omitempty" yaml:"ram,omitempty"`
	Vcpus string `json:"vcpus,omitempty" yaml:"vcpus,omitempty"`
}
