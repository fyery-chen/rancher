package client

const (
	QueryIstioGraphType           = "queryIstioGraph"
	QueryIstioGraphFieldGraphName = "graphID"
	QueryIstioGraphFieldSeries    = "series"
)

type QueryIstioGraph struct {
	GraphName string   `json:"graphID,omitempty" yaml:"graphID,omitempty"`
	Series    []string `json:"series,omitempty" yaml:"series,omitempty"`
}
