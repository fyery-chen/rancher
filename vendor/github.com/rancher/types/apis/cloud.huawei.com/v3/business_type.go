package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Business struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BusinessQuotaSpec `json:"spec"`
	// Most recent observed status of the alert. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status BusinessStatus `json:"status"`
}

type BusinessStatus struct {
	BusinessState string `json:"businessState,omitempty" norman:"options=transferring|transferError|active,default=active"`
}

type BusinessQuotaSpec struct {
	DisplayName string `json:"displayName,omitempty" norman:"required"`
	Description string `json:"description,omitempty"`
	NodeCount   int    `json:"nodeCount,omitempty"`
}

type BusinessQuotaCheck struct {
	BusinessName string `json:"businessName,omitempty" norman:"type=string,required"`
	NodeCount    int    `json:"nodeCount,omitempty" norman:"required"`
}

type BusinessQuotaCheckOutput struct {
	Message string `json:"message,omitempty"`
}

type SubnetInfo struct {
	SubnetName string `json:"name,omitempty"`
	SubnetId   string `json:"id,omitempty"`
	VpcId      string `json:"vpc_id,omitempty"`
}

type HighwaySubnet struct {
	Name         string `json:"name,omitempty"`
	Status       string `json:"status,omitempty"`
	AdminStateUP bool   `json:"admin_state_up,omitempty"`
	Id           string `json:"id,omitempty"`
}

type VpcInfo struct {
	VpcName string `json:"name,omitempty"`
	VpcId   string `json:"id,omitempty"`
}

type NodeFlavor struct {
	Name  string `json:"name,omitempty"`
	Vcpus string `json:"vcpus,omitempty"`
	Ram   int    `json:"ram,omitempty"`
}

type ZoneState struct {
	Available bool `json:"available,omitempty"`
}
type AvailabilityZone struct {
	ZoneState ZoneState `json:"zoneState,omitempty"`
	ZoneName  string    `json:"zoneName,omitempty"`
}

type HuaweiCloudApiInformationInput struct {
	ProjectId string `json:"projectId,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

type SshKeyInfo struct {
	Name string `json:"name,omitempty"`
}

type SshKey struct {
	Keypair SshKeyInfo `json:"keypair,omitempty"`
}

type HuaweiCloudApiInformationOutput struct {
	VpcInfo       []VpcInfo          `json:"vpcs,omitempty"`
	SubnetInfo    []SubnetInfo       `json:"subnets,omitempty"`
	SshKeyName    []SshKey           `json:"keypairs,omitempty"`
	NodeFlavor    []NodeFlavor       `json:"flavors,omitempty"`
	AvailableZone []AvailabilityZone `json:"availabilityZoneInfo,omitempty"`
	HighwaySubnet []HighwaySubnet    `json:"networks,omitempty"`
	Message       string             `json:"message,omitempty"`
}
