package client

const (
	ProjectRegistryTestInputType               = "projectRegistryTestInput"
	ProjectRegistryTestInputFieldDisplayName   = "displayName"
	ProjectRegistryTestInputFieldPassword      = "password"
	ProjectRegistryTestInputFieldProjectName   = "projectId"
	ProjectRegistryTestInputFieldRegistryType  = "registryType"
	ProjectRegistryTestInputFieldServerAddress = "serverAddress"
	ProjectRegistryTestInputFieldUserName      = "userName"
)

type ProjectRegistryTestInput struct {
	DisplayName   string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty"`
	ProjectName   string `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RegistryType  string `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	UserName      string `json:"userName,omitempty" yaml:"userName,omitempty"`
}
