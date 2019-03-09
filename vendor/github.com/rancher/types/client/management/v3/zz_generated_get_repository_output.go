package client

const (
	GetRepositoryOutputType                = "getRepositoryOutput"
	GetRepositoryOutputFieldMessage        = "message"
	GetRepositoryOutputFieldRepositoryInfo = "repositoryInfo"
)

type GetRepositoryOutput struct {
	Message        string             `json:"message,omitempty" yaml:"message,omitempty"`
	RepositoryInfo []HarborRepository `json:"repositoryInfo,omitempty" yaml:"repositoryInfo,omitempty"`
}
