package client

const (
	ClusterRegistryTestInputType               = "clusterRegistryTestInput"
	ClusterRegistryTestInputFieldClusterName   = "clusterId"
	ClusterRegistryTestInputFieldDisplayName   = "displayName"
	ClusterRegistryTestInputFieldPassword      = "password"
	ClusterRegistryTestInputFieldRegistryType  = "registryType"
	ClusterRegistryTestInputFieldServerAddress = "serverAddress"
	ClusterRegistryTestInputFieldUserName      = "userName"
)

type ClusterRegistryTestInput struct {
	ClusterName   string `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	DisplayName   string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty"`
	RegistryType  string `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	UserName      string `json:"userName,omitempty" yaml:"userName,omitempty"`
}
