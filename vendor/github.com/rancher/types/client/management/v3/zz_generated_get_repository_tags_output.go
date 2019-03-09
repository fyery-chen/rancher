package client

const (
	GetRepositoryTagsOutputType                    = "getRepositoryTagsOutput"
	GetRepositoryTagsOutputFieldMessage            = "message"
	GetRepositoryTagsOutputFieldRepositoryTagsInfo = "repositoryTagsInfo"
)

type GetRepositoryTagsOutput struct {
	Message            string          `json:"message,omitempty" yaml:"message,omitempty"`
	RepositoryTagsInfo []RepositoryTag `json:"repositoryTagsInfo,omitempty" yaml:"repositoryTagsInfo,omitempty"`
}
