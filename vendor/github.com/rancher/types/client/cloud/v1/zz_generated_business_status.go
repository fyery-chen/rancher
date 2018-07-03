package client

const (
	BusinessStatusType               = "businessStatus"
	BusinessStatusFieldBusinessState = "businessState"
)

type BusinessStatus struct {
	BusinessState string `json:"businessState,omitempty" yaml:"businessState,omitempty"`
}
