package client

const (
	BusinessQuotaSpecType             = "businessQuotaSpec"
	BusinessQuotaSpecFieldBusinessId  = "businessId"
	BusinessQuotaSpecFieldCpuQuota    = "cpuQuota"
	BusinessQuotaSpecFieldDescription = "description"
	BusinessQuotaSpecFieldDisplayName = "displayName"
	BusinessQuotaSpecFieldMemoryQuota = "memoryQuota"
)

type BusinessQuotaSpec struct {
	BusinessId  string `json:"businessId,omitempty" yaml:"businessId,omitempty"`
	CpuQuota    int64  `json:"cpuQuota,omitempty" yaml:"cpuQuota,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	MemoryQuota int64  `json:"memoryQuota,omitempty" yaml:"memoryQuota,omitempty"`
}
