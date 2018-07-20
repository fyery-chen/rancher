package client

const (
	HuaweiCloudAccountSpecType             = "huaweiCloudAccountSpec"
	HuaweiCloudAccountSpecFieldAccessKey   = "accessKey"
	HuaweiCloudAccountSpecFieldDisplayName = "displayName"
	HuaweiCloudAccountSpecFieldSecretKey   = "secretKey"
	HuaweiCloudAccountSpecFieldUserId      = "userId"
)

type HuaweiCloudAccountSpec struct {
	AccessKey   string `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	SecretKey   string `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	UserId      string `json:"userId,omitempty" yaml:"userId,omitempty"`
}
