package client

const (
	HuaweiCloudApiInformationInputType           = "huaweiCloudApiInformationInput"
	HuaweiCloudApiInformationInputFieldProjectId = "projectId"
	HuaweiCloudApiInformationInputFieldZone      = "zone"
)

type HuaweiCloudApiInformationInput struct {
	ProjectId string `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Zone      string `json:"zone,omitempty" yaml:"zone,omitempty"`
}
