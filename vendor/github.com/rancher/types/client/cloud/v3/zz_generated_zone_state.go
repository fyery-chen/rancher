package client

const (
	ZoneStateType           = "zoneState"
	ZoneStateFieldAvailable = "available"
)

type ZoneState struct {
	Available bool `json:"available,omitempty" yaml:"available,omitempty"`
}
