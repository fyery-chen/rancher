package client

const (
	GlobalRegistryTestInputType               = "globalRegistryTestInput"
	GlobalRegistryTestInputFieldDisplayName   = "displayName"
	GlobalRegistryTestInputFieldPassword      = "password"
	GlobalRegistryTestInputFieldRegistryType  = "registryType"
	GlobalRegistryTestInputFieldServerAddress = "serverAddress"
	GlobalRegistryTestInputFieldUserName      = "userName"
)

type GlobalRegistryTestInput struct {
	DisplayName   string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty"`
	RegistryType  string `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	UserName      string `json:"userName,omitempty" yaml:"userName,omitempty"`
}
