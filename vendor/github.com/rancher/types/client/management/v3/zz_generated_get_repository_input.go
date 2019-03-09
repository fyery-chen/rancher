package client

const (
	GetRepositoryInputType              = "getRepositoryInput"
	GetRepositoryInputFieldIsAll        = "isAll"
	GetRepositoryInputFieldPage         = "page"
	GetRepositoryInputFieldPageSize     = "pageSize"
	GetRepositoryInputFieldProjectID    = "projectId"
	GetRepositoryInputFieldRegistryName = "registryName"
)

type GetRepositoryInput struct {
	IsAll        bool   `json:"isAll,omitempty" yaml:"isAll,omitempty"`
	Page         int64  `json:"page,omitempty" yaml:"page,omitempty"`
	PageSize     int64  `json:"pageSize,omitempty" yaml:"pageSize,omitempty"`
	ProjectID    int64  `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RegistryName string `json:"registryName,omitempty" yaml:"registryName,omitempty"`
}
