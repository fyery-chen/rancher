package client

const (
	GetRepositoryTagsInputType                = "getRepositoryTagsInput"
	GetRepositoryTagsInputFieldPage           = "page"
	GetRepositoryTagsInputFieldPageSize       = "pageSize"
	GetRepositoryTagsInputFieldRegistryName   = "registryName"
	GetRepositoryTagsInputFieldRepositoryName = "repositoryName"
)

type GetRepositoryTagsInput struct {
	Page           int64  `json:"page,omitempty" yaml:"page,omitempty"`
	PageSize       int64  `json:"pageSize,omitempty" yaml:"pageSize,omitempty"`
	RegistryName   string `json:"registryName,omitempty" yaml:"registryName,omitempty"`
	RepositoryName string `json:"repositoryName,omitempty" yaml:"repositoryName,omitempty"`
}
