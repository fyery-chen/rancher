package client

const (
	QueryIstioGraphOutputType      = "queryIstioGraphOutput"
	QueryIstioGraphOutputFieldData = "data"
	QueryIstioGraphOutputFieldType = "type"
)

type QueryIstioGraphOutput struct {
	Data []QueryIstioGraph `json:"data,omitempty" yaml:"data,omitempty"`
	Type string            `json:"type,omitempty" yaml:"type,omitempty"`
}
