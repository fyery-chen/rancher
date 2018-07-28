package client

const (
	HuaweiCloudContainerEngineCongfigType                       = "huaweiCloudContainerEngineCongfig"
	HuaweiCloudContainerEngineCongfigFieldAccessKey             = "accessKey"
	HuaweiCloudContainerEngineCongfigFieldAuthenticatingProxyCa = "authenticatingProxyCa"
	HuaweiCloudContainerEngineCongfigFieldAvailableZone         = "availableZone"
	HuaweiCloudContainerEngineCongfigFieldBMSIsAutoRenew        = "bmsIsAutoRenew"
	HuaweiCloudContainerEngineCongfigFieldBMSPeriodNum          = "bmsPeriodNum"
	HuaweiCloudContainerEngineCongfigFieldBMSPeriodType         = "bmsPeriodType"
	HuaweiCloudContainerEngineCongfigFieldBillingMode           = "billingMode"
	HuaweiCloudContainerEngineCongfigFieldClusterBillingMode    = "clusterBillingMode"
	HuaweiCloudContainerEngineCongfigFieldClusterEipId          = "clusterEipId"
	HuaweiCloudContainerEngineCongfigFieldClusterFlavor         = "clusterFlavor"
	HuaweiCloudContainerEngineCongfigFieldClusterType           = "clusterType"
	HuaweiCloudContainerEngineCongfigFieldContainerNetworkCidr  = "containerNetworkCidr"
	HuaweiCloudContainerEngineCongfigFieldContainerNetworkMode  = "containerNetworkMode"
	HuaweiCloudContainerEngineCongfigFieldDataVolumeSize        = "dataVolumeSize"
	HuaweiCloudContainerEngineCongfigFieldDataVolumeType        = "dataVolumeType"
	HuaweiCloudContainerEngineCongfigFieldDescription           = "description"
	HuaweiCloudContainerEngineCongfigFieldEipBandwidthSize      = "eipBandwidthSize"
	HuaweiCloudContainerEngineCongfigFieldEipChargeMode         = "eipChargeMode"
	HuaweiCloudContainerEngineCongfigFieldEipCount              = "eipCount"
	HuaweiCloudContainerEngineCongfigFieldEipIds                = "eipIds"
	HuaweiCloudContainerEngineCongfigFieldEipShareType          = "eipShareType"
	HuaweiCloudContainerEngineCongfigFieldEipType               = "eipType"
	HuaweiCloudContainerEngineCongfigFieldExternalServerEnabled = "externalServerEnabled"
	HuaweiCloudContainerEngineCongfigFieldHighwaySubnet         = "highwaySubnet"
	HuaweiCloudContainerEngineCongfigFieldLabels                = "labels"
	HuaweiCloudContainerEngineCongfigFieldMasterVersion         = "masterVersion"
	HuaweiCloudContainerEngineCongfigFieldNodeCount             = "nodeCount"
	HuaweiCloudContainerEngineCongfigFieldNodeFlavor            = "nodeFlavor"
	HuaweiCloudContainerEngineCongfigFieldNodeLabels            = "nodeLabels"
	HuaweiCloudContainerEngineCongfigFieldNodeOperationSystem   = "nodeOperationSystem"
	HuaweiCloudContainerEngineCongfigFieldProjectID             = "projectId"
	HuaweiCloudContainerEngineCongfigFieldRootVolumeSize        = "rootVolumeSize"
	HuaweiCloudContainerEngineCongfigFieldRootVolumeType        = "rootVolumeType"
	HuaweiCloudContainerEngineCongfigFieldSecretKey             = "secretKey"
	HuaweiCloudContainerEngineCongfigFieldSshKey                = "sshKey"
	HuaweiCloudContainerEngineCongfigFieldSubnetId              = "subnetId"
	HuaweiCloudContainerEngineCongfigFieldVpcId                 = "vpcId"
	HuaweiCloudContainerEngineCongfigFieldZone                  = "zone"
)

type HuaweiCloudContainerEngineCongfig struct {
	AccessKey             string            `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	AuthenticatingProxyCa string            `json:"authenticatingProxyCa,omitempty" yaml:"authenticatingProxyCa,omitempty"`
	AvailableZone         string            `json:"availableZone,omitempty" yaml:"availableZone,omitempty"`
	BMSIsAutoRenew        string            `json:"bmsIsAutoRenew,omitempty" yaml:"bmsIsAutoRenew,omitempty"`
	BMSPeriodNum          int64             `json:"bmsPeriodNum,omitempty" yaml:"bmsPeriodNum,omitempty"`
	BMSPeriodType         string            `json:"bmsPeriodType,omitempty" yaml:"bmsPeriodType,omitempty"`
	BillingMode           int64             `json:"billingMode,omitempty" yaml:"billingMode,omitempty"`
	ClusterBillingMode    int64             `json:"clusterBillingMode,omitempty" yaml:"clusterBillingMode,omitempty"`
	ClusterEipId          string            `json:"clusterEipId,omitempty" yaml:"clusterEipId,omitempty"`
	ClusterFlavor         string            `json:"clusterFlavor,omitempty" yaml:"clusterFlavor,omitempty"`
	ClusterType           string            `json:"clusterType,omitempty" yaml:"clusterType,omitempty"`
	ContainerNetworkCidr  string            `json:"containerNetworkCidr,omitempty" yaml:"containerNetworkCidr,omitempty"`
	ContainerNetworkMode  string            `json:"containerNetworkMode,omitempty" yaml:"containerNetworkMode,omitempty"`
	DataVolumeSize        int64             `json:"dataVolumeSize,omitempty" yaml:"dataVolumeSize,omitempty"`
	DataVolumeType        string            `json:"dataVolumeType,omitempty" yaml:"dataVolumeType,omitempty"`
	Description           string            `json:"description,omitempty" yaml:"description,omitempty"`
	EipBandwidthSize      int64             `json:"eipBandwidthSize,omitempty" yaml:"eipBandwidthSize,omitempty"`
	EipChargeMode         string            `json:"eipChargeMode,omitempty" yaml:"eipChargeMode,omitempty"`
	EipCount              int64             `json:"eipCount,omitempty" yaml:"eipCount,omitempty"`
	EipIds                []string          `json:"eipIds,omitempty" yaml:"eipIds,omitempty"`
	EipShareType          string            `json:"eipShareType,omitempty" yaml:"eipShareType,omitempty"`
	EipType               string            `json:"eipType,omitempty" yaml:"eipType,omitempty"`
	ExternalServerEnabled bool              `json:"externalServerEnabled,omitempty" yaml:"externalServerEnabled,omitempty"`
	HighwaySubnet         string            `json:"highwaySubnet,omitempty" yaml:"highwaySubnet,omitempty"`
	Labels                map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	MasterVersion         string            `json:"masterVersion,omitempty" yaml:"masterVersion,omitempty"`
	NodeCount             int64             `json:"nodeCount,omitempty" yaml:"nodeCount,omitempty"`
	NodeFlavor            string            `json:"nodeFlavor,omitempty" yaml:"nodeFlavor,omitempty"`
	NodeLabels            map[string]string `json:"nodeLabels,omitempty" yaml:"nodeLabels,omitempty"`
	NodeOperationSystem   string            `json:"nodeOperationSystem,omitempty" yaml:"nodeOperationSystem,omitempty"`
	ProjectID             string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RootVolumeSize        int64             `json:"rootVolumeSize,omitempty" yaml:"rootVolumeSize,omitempty"`
	RootVolumeType        string            `json:"rootVolumeType,omitempty" yaml:"rootVolumeType,omitempty"`
	SecretKey             string            `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	SshKey                string            `json:"sshKey,omitempty" yaml:"sshKey,omitempty"`
	SubnetId              string            `json:"subnetId,omitempty" yaml:"subnetId,omitempty"`
	VpcId                 string            `json:"vpcId,omitempty" yaml:"vpcId,omitempty"`
	Zone                  string            `json:"zone,omitempty" yaml:"zone,omitempty"`
}
