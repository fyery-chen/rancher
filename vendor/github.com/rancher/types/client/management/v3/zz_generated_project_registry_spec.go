package client

const (
	ProjectRegistrySpecType               = "projectRegistrySpec"
	ProjectRegistrySpecFieldDisplayName   = "displayName"
	ProjectRegistrySpecFieldPassword      = "password"
	ProjectRegistrySpecFieldProjectID     = "projectId"
	ProjectRegistrySpecFieldRegistryType  = "registryType"
	ProjectRegistrySpecFieldServerAddress = "serverAddress"
	ProjectRegistrySpecFieldUserName      = "userName"
)

type ProjectRegistrySpec struct {
	DisplayName   string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty"`
	ProjectID     string `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RegistryType  string `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	UserName      string `json:"userName,omitempty" yaml:"userName,omitempty"`
}
