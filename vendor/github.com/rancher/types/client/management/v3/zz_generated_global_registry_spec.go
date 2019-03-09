package client

const (
	GlobalRegistrySpecType               = "globalRegistrySpec"
	GlobalRegistrySpecFieldDisplayName   = "displayName"
	GlobalRegistrySpecFieldPassword      = "password"
	GlobalRegistrySpecFieldRegistryType  = "registryType"
	GlobalRegistrySpecFieldServerAddress = "serverAddress"
	GlobalRegistrySpecFieldUserName      = "userName"
)

type GlobalRegistrySpec struct {
	DisplayName   string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty"`
	RegistryType  string `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	UserName      string `json:"userName,omitempty" yaml:"userName,omitempty"`
}
