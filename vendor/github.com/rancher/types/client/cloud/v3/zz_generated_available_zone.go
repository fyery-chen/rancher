package client

const (
	AvailableZoneType           = "availableZone"
	AvailableZoneFieldZoneName  = "zoneName"
	AvailableZoneFieldZoneState = "zoneState"
)

type AvailableZone struct {
	ZoneName  string `json:"zoneName,omitempty" yaml:"zoneName,omitempty"`
	ZoneState bool   `json:"zoneState,omitempty" yaml:"zoneState,omitempty"`
}
