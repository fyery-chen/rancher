package client

const (
	ClusterRegistrySpecType                        = "clusterRegistrySpec"
	ClusterRegistrySpecFieldClusterID              = "clusterId"
	ClusterRegistrySpecFieldDisplayName            = "displayName"
	ClusterRegistrySpecFieldIncludeSystemComponent = "includeSystemComponent"
	ClusterRegistrySpecFieldPassword               = "password"
	ClusterRegistrySpecFieldRegistryType           = "registryType"
	ClusterRegistrySpecFieldServerAddress          = "serverAddress"
	ClusterRegistrySpecFieldUserName               = "userName"
)

type ClusterRegistrySpec struct {
	ClusterID              string `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	DisplayName            string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	IncludeSystemComponent *bool  `json:"includeSystemComponent,omitempty" yaml:"includeSystemComponent,omitempty"`
	Password               string `json:"password,omitempty" yaml:"password,omitempty"`
	RegistryType           string `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	ServerAddress          string `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	UserName               string `json:"userName,omitempty" yaml:"userName,omitempty"`
}
