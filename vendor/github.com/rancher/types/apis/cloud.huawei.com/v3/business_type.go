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
	SubnetName string `json:"subnetName,omitempty"`
	SubnetId   string `json:"subnetId,omitempty"`
}

type VpcInfo struct {
	VpcName    string       `json:"vpcName,omitempty"`
	VpcId      string       `json:"vpcId,omitempty"`
	SubnetInfo []SubnetInfo `json:"subnetInfo,omitempty"`
}

type NodeFlavor struct {
	FlavorName string `json:"flavorName,omitempty"`
	Cpu        int    `json:"cpu,omitempty"`
	Memory     string `json:"memory,omitempty"`
}

type HuaweiCloudApiInformationInput struct {
	ProjectId string `json:"projectId,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

type HuaweiCloudApiInformationOutput struct {
	VpcInfo       []VpcInfo    `json:"vpcInfo,omitempty"`
	SshKeyName    string       `json:"sshKeyName,omitempty"`
	NodeFlavor    []NodeFlavor `json:"nodeFlavor,omitempty"`
	AvailableZone []string     `json:"availableZone,omitempty"`
	Message       string       `json:"message,omitempty"`
}
