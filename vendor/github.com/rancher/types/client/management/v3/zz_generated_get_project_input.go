package client

const (
	GetProjectInputType              = "getProjectInput"
	GetProjectInputFieldPage         = "page"
	GetProjectInputFieldPageSize     = "pageSize"
	GetProjectInputFieldRegistryName = "registryName"
)

type GetProjectInput struct {
	Page         int64  `json:"page,omitempty" yaml:"page,omitempty"`
	PageSize     int64  `json:"pageSize,omitempty" yaml:"pageSize,omitempty"`
	RegistryName string `json:"registryName,omitempty" yaml:"registryName,omitempty"`
}
