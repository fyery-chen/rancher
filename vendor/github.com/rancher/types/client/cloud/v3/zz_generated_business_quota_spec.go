package client

const (
	BusinessQuotaSpecType             = "businessQuotaSpec"
	BusinessQuotaSpecFieldDescription = "description"
	BusinessQuotaSpecFieldDisplayName = "displayName"
	BusinessQuotaSpecFieldNodeCount   = "nodeCount"
)

type BusinessQuotaSpec struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	NodeCount   int64  `json:"nodeCount,omitempty" yaml:"nodeCount,omitempty"`
}
