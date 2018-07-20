package client

const (
	BusinessQuotaCheckOutputType         = "businessQuotaCheckOutput"
	BusinessQuotaCheckOutputFieldMessage = "message"
)

type BusinessQuotaCheckOutput struct {
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}
