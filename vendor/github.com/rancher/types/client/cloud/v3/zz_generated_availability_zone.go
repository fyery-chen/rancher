package client

const (
	AvailabilityZoneType           = "availabilityZone"
	AvailabilityZoneFieldZoneName  = "zoneName"
	AvailabilityZoneFieldZoneState = "zoneState"
)

type AvailabilityZone struct {
	ZoneName  string     `json:"zoneName,omitempty" yaml:"zoneName,omitempty"`
	ZoneState *ZoneState `json:"zoneState,omitempty" yaml:"zoneState,omitempty"`
}
