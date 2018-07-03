package client

const (
	BusinessQuotaCheckType              = "businessQuotaCheck"
	BusinessQuotaCheckFieldBusinessName = "businessName"
	BusinessQuotaCheckFieldCpuQuota     = "cpuQuota"
	BusinessQuotaCheckFieldMemoryQuota  = "memoryQuota"
	BusinessQuotaCheckFieldNodeCount    = "nodeCount"
)

type BusinessQuotaCheck struct {
	BusinessName string `json:"businessName,omitempty" yaml:"businessName,omitempty"`
	CpuQuota     int64  `json:"cpuQuota,omitempty" yaml:"cpuQuota,omitempty"`
	MemoryQuota  int64  `json:"memoryQuota,omitempty" yaml:"memoryQuota,omitempty"`
	NodeCount    int64  `json:"nodeCount,omitempty" yaml:"nodeCount,omitempty"`
}
