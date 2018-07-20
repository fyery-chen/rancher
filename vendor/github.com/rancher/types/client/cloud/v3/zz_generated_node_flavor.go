package client

const (
	NodeFlavorType            = "nodeFlavor"
	NodeFlavorFieldCpu        = "cpu"
	NodeFlavorFieldFlavorName = "flavorName"
	NodeFlavorFieldMemory     = "memory"
)

type NodeFlavor struct {
	Cpu        int64  `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	FlavorName string `json:"flavorName,omitempty" yaml:"flavorName,omitempty"`
	Memory     string `json:"memory,omitempty" yaml:"memory,omitempty"`
}
