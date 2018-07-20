package client

const (
	BusinessQuotaCheckType              = "businessQuotaCheck"
	BusinessQuotaCheckFieldBusinessName = "businessName"
	BusinessQuotaCheckFieldNodeCount    = "nodeCount"
)

type BusinessQuotaCheck struct {
	BusinessName string `json:"businessName,omitempty" yaml:"businessName,omitempty"`
	NodeCount    int64  `json:"nodeCount,omitempty" yaml:"nodeCount,omitempty"`
}
