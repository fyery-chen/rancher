package client

const (
	HighwaySubnetType              = "highwaySubnet"
	HighwaySubnetFieldAdminStateUP = "admin_state_up"
	HighwaySubnetFieldName         = "name"
	HighwaySubnetFieldStatus       = "status"
)

type HighwaySubnet struct {
	AdminStateUP bool   `json:"admin_state_up,omitempty" yaml:"admin_state_up,omitempty"`
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	Status       string `json:"status,omitempty" yaml:"status,omitempty"`
}
